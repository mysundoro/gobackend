package auth

import (
	"fmt"
	"os"
	"time"

	"gobackend/database"
	"gobackend/models"
	"gobackend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		DeviceID string `json:"device_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal parse request"})
	}

	// Ambil IP & User Agent
	ip := c.IP()
	userAgent := c.Get("User-Agent")

	// Ambil dari settings
	maxFailedUser := utils.GetSettingInt("max_failed_login_user", 3)
	lockDuration := utils.GetSettingInt("lock_duration_minutes", 30)
	maxFailedIP := utils.GetSettingInt("max_failed_login_ip", 5)
	blockIPMinutes := utils.GetSettingInt("ip_block_duration_minutes", 60)

	// Cek apakah IP diblokir
	var blocked models.BlockedIP
	now := time.Now().Unix()
	if err := database.DB.Where("ip = ?", ip).First(&blocked).Error; err == nil {
		if blocked.ExpiresAt > now {
			return c.Status(429).JSON(fiber.Map{
				"error": "IP ini diblokir karena terlalu banyak percobaan login gagal. Coba lagi nanti.",
			})
		}
	}

	// Cari user
	var user models.User
	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		updateIPAttempt(ip, &blocked, maxFailedIP, blockIPMinutes)
		return c.Status(401).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Cek apakah akun terkunci
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return c.Status(403).JSON(fiber.Map{
			"error": fmt.Sprintf("Akun Anda dikunci sampai %s", user.LockedUntil.Format("15:04")),
		})
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		user.FailedLoginCount++

		if user.FailedLoginCount >= maxFailedUser {
			lockTime := time.Now().Add(time.Duration(lockDuration) * time.Minute)
			user.LockedUntil = &lockTime

			// Buat token unlock
			unlockToken := uuid.New().String()
			database.DB.Create(&models.AccountUnlockToken{
				UserID:    user.ID,
				Token:     unlockToken,
				ExpiredAt: time.Now().Add(time.Duration(lockDuration) * time.Minute).Unix(),
			})

			// Kirim email
			go utils.SendAccountLockedEmail(user.Email, lockTime)
			go utils.SendAccountUnlockEmail(user.Email, unlockToken)

			database.DB.Save(&user)
			updateIPAttempt(ip, &blocked, maxFailedIP, blockIPMinutes)

			return c.Status(403).JSON(fiber.Map{
				"error": "Akun dikunci karena terlalu banyak percobaan login gagal. Silakan cek email Anda.",
			})
		}

		database.DB.Save(&user)
		updateIPAttempt(ip, &blocked, maxFailedIP, blockIPMinutes)

		return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
	}

	// Reset IP block jika berhasil login
	if blocked.ID != 0 {
		database.DB.Delete(&blocked)
	}

	// Cek session aktif
	var existingSession models.ActiveSession
	database.DB.Where("user_id = ? AND is_active = true", user.ID).First(&existingSession)

	if existingSession.ID != 0 && existingSession.DeviceID != body.DeviceID {
		return c.Status(409).JSON(fiber.Map{
			"error":  "User sudah login di device lain",
			"action": "ask_override",
		})
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	// Nonaktifkan session lama jika ada
	if existingSession.ID != 0 {
		database.DB.Model(&existingSession).Update("is_active", false)
	}

	// Simpan session baru
	newSession := models.ActiveSession{
		UserID:    user.ID,
		Token:     signed,
		DeviceID:  body.DeviceID,
		IP:        ip,
		UserAgent: userAgent,
		IsActive:  true,
	}
	database.DB.Create(&newSession)

	// Reset gagal login
	user.FailedLoginCount = 0
	user.LockedUntil = nil
	database.DB.Save(&user)

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   signed,
	})
}

func ForgotPassword(c *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Email tidak ditemukan"})
	}

	// Generate token & simpan
	token := uuid.New().String()
	exp := time.Now().Add(15 * time.Minute).Unix()

	reset := models.PasswordReset{
		Email:     user.Email,
		Token:     token,
		ExpiredAt: exp,
	}
	database.DB.Create(&reset)

	// Kirim email reset password
	if err := utils.SendResetPasswordEmail(user.Email, token); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal kirim email"})
	}

	return c.JSON(fiber.Map{"message": "Silakan cek email untuk reset password"})
}

func ResetPassword(c *fiber.Ctx) error {
	token := c.Query("token")
	var body struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}
	if body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password wajib diisi"})
	}

	var reset models.PasswordReset
	if err := database.DB.Where("token = ?", token).First(&reset).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	if time.Now().Unix() > reset.ExpiredAt {
		return c.Status(400).JSON(fiber.Map{"error": "Token sudah kadaluarsa"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", reset.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	hashed := utils.HashPassword(body.Password)
	user.Password = hashed
	database.DB.Save(&user)

	// Hapus token setelah digunakan
	database.DB.Delete(&reset)

	return c.JSON(fiber.Map{"message": "Password berhasil direset"})
}

func UnlockAccount(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Token tidak ditemukan"})
	}

	var unlock models.AccountUnlockToken
	if err := database.DB.Where("token = ?", token).First(&unlock).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	if time.Now().Unix() > unlock.ExpiredAt {
		return c.Status(400).JSON(fiber.Map{"error": "Token sudah kadaluarsa"})
	}

	var user models.User
	if err := database.DB.First(&user, unlock.UserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	user.FailedLoginCount = 0
	user.LockedUntil = nil
	database.DB.Save(&user)

	// hapus token
	database.DB.Delete(&unlock)

	return c.JSON(fiber.Map{"message": "Akun berhasil dibuka, silakan login kembali."})
}

func updateIPAttempt(ip string, existing *models.BlockedIP, maxAttempts int, blockMinutes int) {
	now := time.Now().Unix()
	if existing.ID == 0 {
		block := models.BlockedIP{
			IP:        ip,
			Attempts:  1,
			BlockedAt: now,
		}
		database.DB.Create(&block)
	} else {
		existing.Attempts++
		if existing.Attempts >= maxAttempts {
			existing.ExpiresAt = time.Now().Add(time.Duration(blockMinutes) * time.Minute).Unix()
		}
		database.DB.Save(existing)
	}
}

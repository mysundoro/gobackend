package utils

import (
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

// Kirim email reset password
func SendResetPasswordEmail(to string, token string) error {
	frontend := GetSettingString("frontend_url", "https://yourfrontend.com")
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontend, token)

	m := gomail.NewMessage()
	m.SetHeader("From", GetSettingString("mail_from", "noreply@domain.com"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset Password Akun Anda")
	m.SetBody("text/html", `
		<p>Halo,</p>
		<p>Kami menerima permintaan untuk mereset password Anda.</p>
		<p><a href="`+resetLink+`">Klik di sini untuk reset password</a></p>
		<p>Link ini hanya berlaku selama 15 menit.</p>
	`)

	d := gomail.NewDialer(
		GetSettingString("mail_host", "smtp.gmail.com"),
		GetSettingInt("mail_port", 587),
		GetSettingString("mail_user", ""),
		GetSettingString("mail_pass", ""),
	)

	return d.DialAndSend(m)
}

// Kirim email saat akun dikunci
func SendAccountLockedEmail(to string, until time.Time) error {
	m := gomail.NewMessage()
	m.SetHeader("From", GetSettingString("mail_from", "noreply@domain.com"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Akun Anda Dikunci Sementara")
	m.SetBody("text/html", fmt.Sprintf(`
		<p>Halo,</p>
		<p>Akun Anda telah dikunci sementara karena terlalu banyak percobaan login gagal.</p>
		<p>Akun akan terbuka kembali pada <strong>%s</strong></p>
		<p>Jika ini bukan Anda, segera ubah password Anda begitu akun tersedia.</p>
	`, until.Format("02 Jan 2006 15:04")))

	d := gomail.NewDialer(
		GetSettingString("mail_host", "smtp.gmail.com"),
		GetSettingInt("mail_port", 587),
		GetSettingString("mail_user", ""),
		GetSettingString("mail_pass", ""),
	)

	return d.DialAndSend(m)
}

// Kirim email untuk unlock akun
func SendAccountUnlockEmail(email string, token string) error {
	frontend := GetSettingString("frontend_url", "https://yourfrontend.com")
	unlockURL := fmt.Sprintf("%s/unlock-account?token=%s", frontend, token)

	m := gomail.NewMessage()
	m.SetHeader("From", GetSettingString("mail_from", "noreply@domain.com"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Akun Anda Terkunci - Buka Sekarang")
	m.SetBody("text/html", fmt.Sprintf(`
		<p>Halo,</p>
		<p>Akun Anda telah dikunci karena terlalu banyak percobaan login gagal.</p>
		<p>Klik link berikut untuk membuka akun Anda:</p>
		<p><a href="%s">%s</a></p>
		<p>Link ini hanya berlaku selama 30 menit.</p>
	`, unlockURL, unlockURL))

	d := gomail.NewDialer(
		GetSettingString("mail_host", "smtp.gmail.com"),
		GetSettingInt("mail_port", 587),
		GetSettingString("mail_user", ""),
		GetSettingString("mail_pass", ""),
	)

	return d.DialAndSend(m)
}

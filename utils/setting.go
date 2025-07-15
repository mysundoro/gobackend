package utils

import (
	"gobackend/database"
	"gobackend/models"
	"log"
	"strconv"
)

// GetSettingString mengambil nilai string dari DB berdasarkan key
func GetSettingString(key string, defaultVal string) string {
	var s models.Setting
	if err := database.DB.Where("`key` = ?", key).First(&s).Error; err != nil {
		return defaultVal
	}
	return s.Value
}

// GetSettingInt mengambil nilai int dari DB berdasarkan key
func GetSettingInt(key string, defaultVal int) int {
	var s models.Setting
	if err := database.DB.Where("`key` = ?", key).First(&s).Error; err != nil {
		return defaultVal
	}

	val, err := strconv.Atoi(s.Value)
	if err != nil {
		log.Printf("⚠️ Setting '%s' bukan angka: %v", key, err)
		return defaultVal
	}

	return val
}

// GetSettingBool mengambil nilai boolean dari DB berdasarkan key
func GetSettingBool(key string, defaultVal bool) bool {
	var s models.Setting
	if err := database.DB.Where("`key` = ?", key).First(&s).Error; err != nil {
		return defaultVal
	}

	val, err := strconv.ParseBool(s.Value)
	if err != nil {
		log.Printf("⚠️ Setting '%s' bukan boolean: %v", key, err)
		return defaultVal
	}
	return val
}

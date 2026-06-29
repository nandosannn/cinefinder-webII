package config

import "os"

func GetJWTKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "cinefinder_dev_secret_123456"
	}
	return []byte(secret)
}

package config

import (
	"fmt"
	"time"

	gobetterauth "github.com/GoBetterAuth/go-better-auth"
	gobetterauthdomain "github.com/GoBetterAuth/go-better-auth/pkg/domain"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func NewAuth(config *viper.Viper, db *gorm.DB, logger *logrus.Logger) *gobetterauth.Auth {
	configAuth := gobetterauthdomain.NewConfig()
	configAuth.AppName = config.GetString("app.name")
	configAuth.EmailPassword = gobetterauthdomain.EmailPasswordConfig{
		Enabled:                  true,
		RequireEmailVerification: true,
	}
	configAuth.WithTrustedOrigins(gobetterauthdomain.TrustedOriginsConfig{
		Origins: []string{
			"http://localhost:3000",
		},
	})
	configAuth.EmailVerification = gobetterauthdomain.EmailVerificationConfig{
		SendOnSignUp: true,
		SendOnSignIn: true,
		AutoSignIn:   false,
		ExpiresIn:    1 * time.Hour,
		SendVerificationEmail: func(user *gobetterauthdomain.User, url string, token string) error {
			// Implement email sending logic
			fmt.Printf("Send verification email to %s with url: %s", user.Email, url)
			return nil
		},
	}
	configAuth.DatabaseHooks.Users.BeforeCreate = func(user *gobetterauthdomain.User) error {
		logger.Infof("Creating user: %s", user.Email)
		return nil
	}

	auth := gobetterauth.New(configAuth, db)
	auth.RunMigrations()

	logger.Info("Auth Running Successfully")

	return auth
}

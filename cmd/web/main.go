package main

import (
	"fmt"
	"hr-sas/internal/config"
	"hr-sas/internal/pkg"
)

func main() {
	viperConfig := config.NewViper()
	app := config.NewFiber(viperConfig)
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	s3 := pkg.NewS3Client(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		S3Client: s3,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

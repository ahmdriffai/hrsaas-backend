package http

import (
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type NotificationController struct {
	Log *logrus.Logger
}

func NewNotificationController(log *logrus.Logger) *NotificationController {
	return &NotificationController{
		Log: log,
	}
}

func (c *NotificationController) RegisterRoutes(ctx *fiber.Ctx) error {
	app, err := firebase.NewApp(ctx.UserContext(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Messaging(ctx.UserContext())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "YOUR_REGISTRATION_TOKEN"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx.UserContext(), message)
	if err != nil {
		log.Fatalln(err)
	}

	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)

	return ctx.JSON(fiber.Map{
		"message": "Notification sent successfully",
	})
}

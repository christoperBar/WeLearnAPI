package main

import (
	"github.com/christoperBar/WeLearnAPI/controllers/studentcontroller"
	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/gofiber/fiber/v2"
)

func main() {

	models.ConnectDatabase()

	app := fiber.New()

	app.Post("/api/students/register", studentcontroller.Register)
	app.Get("/api/students/:id", studentcontroller.UserProfile)

	app.Listen(":8000")
}

package main

import (
	"github.com/christoperBar/WeLearnAPI/controllers/instructorcontroller"
	"github.com/christoperBar/WeLearnAPI/controllers/studentcontroller"
	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/gofiber/fiber/v2"
)

func main() {

	models.ConnectDatabase()

	app := fiber.New()

	//students
	app.Post("/api/students/register", studentcontroller.Register)
	app.Get("/api/students/:id", studentcontroller.UserProfile)

	//sayembara
	app.Get("/api/students/:id/sayembaras", studentcontroller.SayembaraList)
	app.Get("/api/students/:id/sayembaras/:sayembaraid", studentcontroller.SayembaraDetail)
	// app.Put("/api/students/:id/sayembaras/:sayembaraid", studentcontroller.EditStatussayembara)
	app.Post("/api/students/:id/sayembaras/", studentcontroller.CreateSayembara)

	//instrucors
	app.Get("/api/instructors", instructorcontroller.InstructorList)
	app.Get("/api/instructors/:id", instructorcontroller.InstructorDetail)

	app.Listen(":8000")
}

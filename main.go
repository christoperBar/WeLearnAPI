package main

import (
	"github.com/christoperBar/WeLearnAPI/controllers/instructorcontroller"
	"github.com/christoperBar/WeLearnAPI/controllers/studentcontroller"
	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

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

	//lessons
	app.Get("/api/instructors/:id/lessons", instructorcontroller.LessonList)
	app.Get("/api/instructors/:id/lessons/:lessonid", instructorcontroller.LessonDetail)

	app.Listen(":8000")
}

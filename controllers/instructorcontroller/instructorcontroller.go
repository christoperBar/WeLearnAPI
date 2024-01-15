package instructorcontroller

import (
	// "fmt"
	"net/http"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InstructorExpertiseDTO struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type InstructorListDTO struct {
	Id         uuid.UUID                `json:"id"`
	AuthId     string                   `json:"authid"`
	DOB        string                   `json:"dob"`
	Address    string                   `json:"address"`
	Phone      string                   `json:"phone"`
	Image_url  string                   `json:"image_url"`
	Expertises []InstructorExpertiseDTO `json:"expertises"`
}

func InstructorList(c *fiber.Ctx) error {
	var instrucors []models.Instructor

	if err := models.DB.Preload("Expertises").Find(&instrucors).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var instructorsDTO []InstructorListDTO
	for _, instructor := range instrucors {

		var expertisesDTO []InstructorExpertiseDTO
		for _, expertise := range instructor.Expertises {
			expertisesDTO = append(expertisesDTO, InstructorExpertiseDTO{
				Id:   expertise.Id,
				Name: expertise.Name,
			})
		}

		instructorDTO := InstructorListDTO{
			Id:         instructor.Id,
			AuthId:     instructor.AuthId,
			DOB:        instructor.DOB,
			Address:    instructor.Address,
			Phone:      instructor.Phone,
			Image_url:  instructor.Image_url,
			Expertises: expertisesDTO,
		}

		instructorsDTO = append(instructorsDTO, instructorDTO)
	}

	return c.Status(fiber.StatusOK).JSON(instructorsDTO)
}

func InstructorDetail(c *fiber.Ctx) error {

	id := c.Params("id")

	var instructor models.Instructor

	if err := models.DB.Preload("Expertises").Where("id = ?", id).First(&instructor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Sayembara not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var expertisesDTO []InstructorExpertiseDTO
	for _, expertise := range instructor.Expertises {
		expertisesDTO = append(expertisesDTO, InstructorExpertiseDTO{
			Id:   expertise.Id,
			Name: expertise.Name,
		})
	}

	instructorDetails := fiber.Map{
		"id":         instructor.Id,
		"authid":     instructor.AuthId,
		"dob":        instructor.DOB,
		"address":    instructor.Address,
		"phone":      instructor.Phone,
		"image_url":  instructor.Image_url,
		"expertises": expertisesDTO,
	}
	return c.Status(fiber.StatusOK).JSON(instructorDetails)
}

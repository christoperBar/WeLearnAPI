package studentcontroller

import (
	"net/http"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// sayembara
type SayembaraExpertiseDTO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SayembaraCategoryDTO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SayembaraStudentDTO struct {
	Id        int64  `json:"id"`
	AuthID    string `json:"authid"`
	DOB       string `json:"dob"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Image_url string `json:"image_url"`
}

type SayembaraBudgetDTO struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

type SayembaraListDTO struct {
	ID          int64                   `json:"id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Budget      SayembaraBudgetDTO      `json:"budget"`
	ImageURL    string                  `json:"image_url"`
	Status      string                  `json:"status"`
	Student     SayembaraStudentDTO     `json:"student"`
	Category    SayembaraCategoryDTO    `json:"category"`
	Expertises  []SayembaraExpertiseDTO `json:"expertises"`
}

func SayembaraList(c *fiber.Ctx) error {
	var sayembaras []models.Sayembara

	if err := models.DB.Preload("Student").Preload("Category").Preload("Expertises").Find(&sayembaras).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var sayembarasDTO []SayembaraListDTO
	for _, sayembara := range sayembaras {
		// Mengonversi Budget
		var expertisesDTO []SayembaraExpertiseDTO
		for _, expertise := range sayembara.Expertises {
			expertisesDTO = append(expertisesDTO, SayembaraExpertiseDTO{
				Id:   expertise.Id,
				Name: expertise.Name,
			})
		}

		categoryDTO := SayembaraCategoryDTO{
			Id:   sayembara.Category.Id,
			Name: sayembara.Category.Name,
		}

		studentDTO := SayembaraStudentDTO{
			Id:        sayembara.Student.Id,
			AuthID:    sayembara.Student.AuthId,
			DOB:       sayembara.Student.DOB,
			Address:   sayembara.Student.Address,
			Phone:     sayembara.Student.Phone,
			Image_url: sayembara.Student.Image_url,
		}

		budgetDTO := SayembaraBudgetDTO{
			Min: float32(sayembara.Budget_min),
			Max: float32(sayembara.Budget_max),
		}

		sayembaraDTO := SayembaraListDTO{
			ID:          sayembara.Id,
			Title:       sayembara.Title,
			Description: sayembara.Description,
			Budget:      budgetDTO,
			ImageURL:    sayembara.Image_url,
			Status:      sayembara.Status,
			Student:     studentDTO,
			Category:    categoryDTO,
			Expertises:  expertisesDTO,
		}

		sayembarasDTO = append(sayembarasDTO, sayembaraDTO)
	}

	return c.Status(fiber.StatusOK).JSON(sayembarasDTO)
}

func SayembaraDetail(c *fiber.Ctx) error {

	id := c.Params("sayembaraid")

	var sayembara models.Sayembara

	if err := models.DB.Preload("Student").Preload("Category").Preload("Expertises").Where("id = ?", id).First(&sayembara).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Sayembara not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var expertisesDTO []SayembaraExpertiseDTO
	for _, expertise := range sayembara.Expertises {
		expertisesDTO = append(expertisesDTO, SayembaraExpertiseDTO{
			Id:   expertise.Id,
			Name: expertise.Name,
		})
	}

	categoryDTO := SayembaraCategoryDTO{
		Id:   sayembara.Category.Id,
		Name: sayembara.Category.Name,
	}

	studentDTO := SayembaraStudentDTO{
		Id:        sayembara.Student.Id,
		AuthID:    sayembara.Student.AuthId,
		DOB:       sayembara.Student.DOB,
		Address:   sayembara.Student.Address,
		Phone:     sayembara.Student.Phone,
		Image_url: sayembara.Student.Image_url,
	}

	budgetDTO := SayembaraBudgetDTO{
		Min: float32(sayembara.Budget_min),
		Max: float32(sayembara.Budget_max),
	}

	sayembaraDetails := fiber.Map{
		"id":          sayembara.Id,
		"title":       sayembara.Title,
		"description": sayembara.Description,
		"budget":      budgetDTO,
		"image_url":   sayembara.Image_url,
		"status":      sayembara.Status,
		"student":     studentDTO,
		"category":    categoryDTO,
		"expertises":  expertisesDTO,
	}

	return c.Status(fiber.StatusOK).JSON(sayembaraDetails)
}

type CreateSayembaraRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	CategoryID  int64   `json:"category_id" validate:"required"`
	Expertises  []int64 `json:"expertises_id" validate:"required"`
	Budget      struct {
		Min float32 `json:"min" validate:"required"`
		Max float32 `json:"max" validate:"required"`
	} `json:"budget" validate:"required"`
	Image_url string `json:"image_url" validate:"required"`
}

func CreateSayembara(c *fiber.Ctx) error {

	id := c.Params("id")

	var student models.Student
	if err := models.DB.First(&student, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Data not Found",
			})
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Data not Found",
		})
	}

	var request CreateSayembaraRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  err.Error(),
		})
	}

	sayembara := models.Sayembara{
		Title:       request.Title,
		Description: request.Description,
		Budget_min:  request.Budget.Min,
		Budget_max:  request.Budget.Max,
		Category_ID: int64(request.CategoryID),
		Student_ID:  student.Id,
		Image_url:   request.Image_url,
		Status:      "On Going",
	}

	if err := models.DB.Create(&sayembara).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Sayembara",
		})
	}

	for _, expertiseID := range request.Expertises {
		expertise := models.Expertise{}
		if err := models.DB.First(&expertise, expertiseID).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to find Expertise",
			})
		}
		models.DB.Model(&sayembara).Association("Expertises").Append(&expertise)
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

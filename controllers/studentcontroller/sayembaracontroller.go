package studentcontroller

import (
	"net/http"
	"time"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// sayembara
type SayembaraExpertiseDTO struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SayembaraCategoryDTO struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SayembaraStudentDTO struct {
	Id        uuid.UUID `json:"id"`
	AuthID    string    `json:"authid"`
	DOB       string    `json:"dob"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Image_url string    `json:"image_url"`
}

type SayembaraBudgetDTO struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

type SayembaraListDTO struct {
	ID          uuid.UUID               `json:"id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Budget      SayembaraBudgetDTO      `json:"budget"`
	ImageURL    string                  `json:"image_url"`
	Status      string                  `json:"status"`
	Student     SayembaraStudentDTO     `json:"student"`
	Category    SayembaraCategoryDTO    `json:"category"`
	Expertises  []SayembaraExpertiseDTO `json:"expertises"`
	CreatedAt   time.Time               `json:"createdat"`
	ClosedAt    time.Time               `json:"closedat"`
}

func SayembaraList(c *fiber.Ctx) error {

	id := c.Params("id")
	searchExpertise := c.Query("expertise")
	searchCategory := c.Query("category")

	var sayembaras []models.Sayembara

	query := models.DB.Preload("Student").Preload("Category").Preload("Expertises").Where("student_id = ?", id)

	if searchExpertise != "" {
		query = query.Joins("JOIN sayembara_expertises ON sayembara_expertises.sayembara_id = sayembaras.id").
			Joins("JOIN expertises ON expertises.id = sayembara_expertises.expertise_id").
			Where("expertises.name LIKE ?", "%"+searchExpertise+"%")
	}

	if searchCategory != "" {
		query = query.Joins("JOIN categories ON categories.id = sayembaras.category_id").
			Where("categories.name LIKE ?", "%"+searchCategory+"%")
	}

	if err := query.Find(&sayembaras).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var sayembarasDTO []SayembaraListDTO
	for _, sayembara := range sayembaras {

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
			CreatedAt:   sayembara.CreatedAt,
			ClosedAt:    sayembara.ClosedAt,
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
		"createdat:":  sayembara.CreatedAt,
		"closedat":    sayembara.ClosedAt,
	}

	return c.Status(fiber.StatusOK).JSON(sayembaraDetails)
}

type CreateSayembaraRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	CategoryID  string   `json:"category_id" validate:"required"`
	Expertises  []string `json:"expertises_id" validate:"required"`
	Budget      struct {
		Min float32 `json:"min" validate:"required"`
		Max float32 `json:"max" validate:"required"`
	} `json:"budget" validate:"required"`
	Image_url string `json:"image_url"`
	ClosedAt  string `json:"closedat" validate:"required"`
}

func CreateSayembara(c *fiber.Ctx) error {

	id := c.Params("id")

	var student models.Student
	if err := models.DB.Where("id = ?", id).First(&student).Error; err != nil {
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

	categoryUUID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid UUID format for CategoryID",
		})
	}

	closedAtTime, err := time.Parse("2006-01-02 15:04:05", request.ClosedAt)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid format for Closed_At",
		})
	}

	sayembara := models.Sayembara{
		Title:       request.Title,
		Description: request.Description,
		Budget_min:  request.Budget.Min,
		Budget_max:  request.Budget.Max,
		Category_ID: categoryUUID,
		Student:     student,
		Image_url:   request.Image_url,
		ClosedAt:    closedAtTime,
		Status:      "On Going",
	}

	if err := models.DB.Create(&sayembara).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Sayembara",
		})
	}

	for _, expertiseID := range request.Expertises {
		newsayembara := models.Sayembara{}
		expertise := models.Expertise{}
		if err := models.DB.Where("id = ?", expertiseID).First(&expertise).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to find Expertise",
			})
		}
		if err := models.DB.Where("title = ?", request.Title).First(&newsayembara).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to find Sayembara",
			})
		}

		models.DB.Model(&newsayembara).Association("Expertises").Append(&expertise)
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

package instructorcontroller

import (
	// "fmt"

	"net/http"
	"strings"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//Lesson

type CreateLessonRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Price       float32  `json:"price" validate:"required"`
	Tags        []string `json:"tags" validate:"required"`
	Method      []string `json:"method" validate:"required"`
	Image_url   string   `json:"image_url"`
	Category_ID string   `json:"category_id" validate:"required"`
}

func CreateLesson(c *fiber.Ctx) error {

	id := c.Params("id")

	var instructor models.Instructor
	if err := models.DB.Where("id = ?", id).First(&instructor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Data not Found",
			})
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Data not Found",
		})
	}

	var request CreateLessonRequest
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

	categoryUUID, err := uuid.Parse(request.Category_ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid UUID format for CategoryID",
		})
	}

	tags := strings.Join(request.Tags, ", ")
	method := strings.Join(request.Method, ", ")

	lesson := models.Lesson{
		Title:         request.Title,
		Description:   request.Description,
		Price:         request.Price,
		Tags:          tags,
		Method:        method,
		Image_url:     request.Image_url,
		Category_ID:   categoryUUID,
		Instructor_ID: instructor.Id,
	}

	if err := models.DB.Create(&lesson).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Sayembara",
		})
	}
	return c.Status(fiber.StatusOK).Send(nil)
}

type LessonInstructorDTO struct {
	Id         uuid.UUID                `json:"id"`
	DOB        string                   `json:"dob"`
	Address    string                   `json:"address"`
	Phone      string                   `json:"phone"`
	Image_url  string                   `json:"image_url"`
	Expertises []InstructorExpertiseDTO `json:"expertises"`
}

type LessonCategoryDTO struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type LessonListDTO struct {
	Id          uuid.UUID           `json:"id"`
	Instructor  LessonInstructorDTO `json:"instructor"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Price       float32             `json:"price"`
	Category    LessonCategoryDTO   `json:"category"`
	Tags        []string            `json:"tags"`
	Method      []string            `json:"method"`
	Image_url   string              `json:"image_url"`
}

func LessonList(c *fiber.Ctx) error {
	id := c.Params("id")
	var lessons []models.Lesson

	if err := models.DB.Preload("Category").Where("instructor_id = ?", id).Find(&lessons).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	var instructor models.Instructor
	models.DB.Preload("Expertises").Where("id = ?", id).First(&instructor)

	var lessonsDTO []LessonListDTO
	for _, lesson := range lessons {

		var expertisesDTO []InstructorExpertiseDTO
		for _, expertise := range instructor.Expertises {
			expertisesDTO = append(expertisesDTO, InstructorExpertiseDTO{
				Id:   expertise.Id,
				Name: expertise.Name,
			})
		}

		categoryDTO := LessonCategoryDTO{
			Id:   lesson.Category.Id,
			Name: lesson.Category.Name,
		}

		instructorDTO := LessonInstructorDTO{
			Id:         instructor.Id,
			DOB:        instructor.DOB,
			Address:    instructor.Address,
			Phone:      instructor.Phone,
			Image_url:  instructor.Image_url,
			Expertises: expertisesDTO,
		}

		tags := strings.Split(lesson.Tags, ",")
		methods := strings.Split(lesson.Method, ",")

		lessonDTO := LessonListDTO{
			Id:          lesson.Id,
			Instructor:  instructorDTO,
			Title:       lesson.Title,
			Description: lesson.Description,
			Price:       lesson.Price,
			Category:    categoryDTO,
			Tags:        tags,
			Method:      methods,
			Image_url:   lesson.Image_url,
		}

		lessonsDTO = append(lessonsDTO, lessonDTO)
	}

	return c.Status(fiber.StatusOK).JSON(lessonsDTO)
}

func LessonDetail(c *fiber.Ctx) error {
	id := c.Params("lessonid")
	instructorid := c.Params("id")

	var lesson models.Lesson

	if err := models.DB.Preload("Category").Preload("Instructor").Where("id = ?", id).First(&lesson).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Lesson not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var instructor models.Instructor
	models.DB.Preload("Expertises").Where("id = ?", instructorid).First(&instructor)

	var expertisesDTO []InstructorExpertiseDTO
	for _, expertise := range instructor.Expertises {
		expertisesDTO = append(expertisesDTO, InstructorExpertiseDTO{
			Id:   expertise.Id,
			Name: expertise.Name,
		})
	}

	categoryDTO := LessonCategoryDTO{
		Id:   lesson.Category.Id,
		Name: lesson.Category.Name,
	}

	instructorDTO := LessonInstructorDTO{
		Id:         instructor.Id,
		DOB:        instructor.DOB,
		Address:    instructor.Address,
		Phone:      instructor.Phone,
		Image_url:  instructor.Image_url,
		Expertises: expertisesDTO,
	}

	tags := strings.Split(lesson.Tags, ",")
	methods := strings.Split(lesson.Method, ",")

	lessonDTO := LessonListDTO{
		Id:          lesson.Id,
		Instructor:  instructorDTO,
		Title:       lesson.Title,
		Description: lesson.Description,
		Price:       lesson.Price,
		Category:    categoryDTO,
		Tags:        tags,
		Method:      methods,
		Image_url:   lesson.Image_url,
	}

	return c.Status(fiber.StatusOK).JSON(lessonDTO)
}

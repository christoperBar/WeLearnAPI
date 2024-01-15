package instructorcontroller

import (
	// "fmt"

	"net/http"
	"strings"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

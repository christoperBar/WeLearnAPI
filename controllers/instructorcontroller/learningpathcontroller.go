package instructorcontroller

import (
	"net/http"
	"strings"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//learningpath

type CreateLearningPathRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Lessons     []string `json:"lessons" validate:"required"`
}

func CreateLearningPath(c *fiber.Ctx) error {

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

	var request CreateLearningPathRequest
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

	learningpath := models.Learning_path{
		Title:         request.Title,
		Description:   request.Description,
		Instructor_ID: instructor.Id,
	}

	if err := models.DB.Create(&learningpath).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Sayembara",
		})
	}

	for _, lessonID := range request.Lessons {
		newlearningpath := models.Learning_path{}
		lesson := models.Lesson{}
		if err := models.DB.Where("id = ?", lessonID).First(&lesson).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to find Expertise",
			})
		}
		if err := models.DB.Where("title = ?", request.Title).First(&newlearningpath).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to find Sayembara",
			})
		}

		models.DB.Model(&newlearningpath).Association("Lessons").Append(&lesson)
	}
	return c.Status(fiber.StatusOK).Send(nil)
}

type Learning_pathListDTO struct {
	Id          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Lessons     []LessonListDTO `json:"lesson"`
}

func LearningPathList(c *fiber.Ctx) error {

	id := c.Params("id")
	var learningpaths []models.Learning_path

	if err := models.DB.Preload("Lessons").Where("instructor_id = ?", id).Find(&learningpaths).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var instructor models.Instructor
	models.DB.Preload("Expertises").Where("id = ?", id).First(&instructor)

	var learning_pathsDTO []Learning_pathListDTO
	for _, learningpath := range learningpaths {

		var lessonsDTO []LessonListDTO
		for _, lesson := range learningpath.Lessons {
			var expertisesDTO []InstructorExpertiseDTO
			for _, expertise := range instructor.Expertises {
				expertisesDTO = append(expertisesDTO, InstructorExpertiseDTO{
					Id:   expertise.Id,
					Name: expertise.Name,
				})
			}
			var category models.Category
			models.DB.Where("id = ?", lesson.Category_ID).First(&category)

			categoryDTO := LessonCategoryDTO{
				Id:   category.Id,
				Name: category.Name,
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
		learning_pathDTO := Learning_pathListDTO{
			Id:          learningpath.Id,
			Title:       learningpath.Title,
			Description: learningpath.Description,
			Lessons:     lessonsDTO,
		}
		learning_pathsDTO = append(learning_pathsDTO, learning_pathDTO)
	}

	return c.Status(fiber.StatusOK).JSON(learning_pathsDTO)
}

func LearningPathDetail(c *fiber.Ctx) error {

	id := c.Params("learningpathid")
	instructorid := c.Params("id")

	var learning_path models.Learning_path

	if err := models.DB.Preload("Lessons").Where("id = ?", id).Find(&learning_path).Error; err != nil {
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

	var lessonsDTO []LessonListDTO
	for _, lesson := range learning_path.Lessons {
		var expertisesDTO []InstructorExpertiseDTO
		for _, expertise := range instructor.Expertises {
			expertisesDTO = append(expertisesDTO, InstructorExpertiseDTO{
				Id:   expertise.Id,
				Name: expertise.Name,
			})
		}
		var category models.Category
		models.DB.Where("id = ?", lesson.Category_ID).First(&category)

		categoryDTO := LessonCategoryDTO{
			Id:   category.Id,
			Name: category.Name,
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
	learning_pathDTO := Learning_pathListDTO{
		Id:          learning_path.Id,
		Title:       learning_path.Title,
		Description: learning_path.Description,
		Lessons:     lessonsDTO,
	}
	return c.Status(fiber.StatusOK).JSON(learning_pathDTO)
}

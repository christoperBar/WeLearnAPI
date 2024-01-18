package studentcontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	fbAdmin "github.com/christoperBar/WeLearnAPI/config/auth"
)

var validate = validator.New()

// user
func Register(c *fiber.Ctx) error {

	var requestData struct {
		Username    string   `json:"username"`
		Phone       string   `json:"phone"`
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		DOB         string   `json:"DOB"`
		Category_ID string   `json:"category_id"`
		Method      []string `json:"method"`
		Address     struct {
			StreetNumber             int    `json:"street_number"`
			Route                    string `json:"route"`
			PostalCode               int    `json:"postal_code"`
			Locality                 string `json:"locality"`
			AdministrativeAreaLevel1 string `json:"administrative_area_level_1"`
		} `json:"address"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var imageUrl string = "https://api.multiavatar.com/" + requestData.Username + ".svg"

	student := models.Student{
		DOB:       requestData.DOB,
		Phone:     requestData.Phone,
		Image_url: imageUrl,
	}
	client, InitAutherr := fbAdmin.InitAtuh().Auth(c.Context())
	if InitAutherr != nil {
		panic("init fb failed" + InitAutherr.Error())
	}
	newUser := fbAdmin.CreateUser(
		c.Context(),
		client,
		student,
		requestData.Email,
		requestData.Password,
		requestData.Username,
		false,
	)

	var address = fmt.Sprintf("%d %s, %d, %s, %s", requestData.Address.StreetNumber, requestData.Address.Route, requestData.Address.PostalCode, requestData.Address.Locality, requestData.Address.AdministrativeAreaLevel1)

	categoryUUID, err := uuid.Parse(requestData.Category_ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid UUID format for CategoryID",
		})
	}
	method := strings.Join(requestData.Method, ", ")

	createdStudent := models.Student{
		Id:          newUser.UID,
		DOB:         requestData.DOB,
		Address:     address,
		Phone:       requestData.Phone,
		Image_url:   imageUrl,
		Category_ID: categoryUUID,
		Method:      method,
	}

	if err := validate.Struct(createdStudent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}

	if err := models.DB.Create(&createdStudent).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

func UserProfile(c *fiber.Ctx) error {

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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":        student.Id,
		"dob":       student.DOB,
		"address":   student.Address,
		"phone":     student.Phone,
		"image_url": student.Image_url,
	})

}

type ExpertiseDTO struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
type CategoryDTO struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func Filter(c *fiber.Ctx) error {

	onlyCategory := c.Query("only")

	if onlyCategory == "category" {
		var categories []models.Category
		if err := models.DB.Find(&categories).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		var categoriesDTO []CategoryDTO
		for _, category := range categories {
			categoriesDTO = append(categoriesDTO, CategoryDTO{
				Id:   category.Id,
				Name: category.Name,
			})
		}
		return c.Status(fiber.StatusOK).JSON(categoriesDTO)
	}
	var categories []models.Category
	if err := models.DB.Find(&categories).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	var categoriesDTO []CategoryDTO
	for _, category := range categories {
		categoriesDTO = append(categoriesDTO, CategoryDTO{
			Id:   category.Id,
			Name: category.Name,
		})
	}

	var tags []string
	if err := models.DB.Raw("SELECT DISTINCT tags FROM lessons").Pluck("tags", &tags).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Proses untuk memisahkan tag
	allTags := make([]string, 0)
	for _, tagString := range tags {
		tagList := strings.Split(tagString, ",")
		allTags = append(allTags, tagList...)
	}

	// Menghapus duplikat tag
	uniqueTags := make(map[string]bool)
	for _, tag := range allTags {
		uniqueTags[tag] = true
	}

	// Mengubah map ke slice
	finalTags := make([]string, 0, len(uniqueTags))
	for tag := range uniqueTags {
		finalTags = append(finalTags, tag)
	}

	var expertises []models.Expertise
	if err := models.DB.Find(&expertises).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	var expertisesDTO []ExpertiseDTO
	for _, expertise := range expertises {
		expertisesDTO = append(expertisesDTO, ExpertiseDTO{
			Id:   expertise.Id,
			Name: expertise.Name,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tags":       finalTags,
		"categories": categoriesDTO,
		"expertises": expertisesDTO,
	})
}

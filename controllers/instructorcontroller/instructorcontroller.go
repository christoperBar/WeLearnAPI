package instructorcontroller

import (
	// "fmt"
	"fmt"
	"net/http"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// instructor

var validate = validator.New()

func Register(c *fiber.Ctx) error {

	var requestData struct {
		AuthID  string `json:"authId"`
		DOB     string `json:"DOB"`
		Address struct {
			StreetNumber             int    `json:"street_number"`
			Route                    string `json:"route"`
			PostalCode               int    `json:"postal_code"`
			Locality                 string `json:"locality"`
			AdministrativeAreaLevel1 string `json:"administrative_area_level_1"`
		} `json:"address"`
		Phone         string `json:"phone"`
		ImageURL      string `json:"image_url"`
		Email         string `json:"email"`
		ID_type       string `json:"id_type"`
		IDcard_number string `json:"idcard_number"`
		IDcard_url    string `json:"idcard_url"`
		Selfie_url    string `json:"selfie_url"`
		Category_ID   string `json:"category_id"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var address = fmt.Sprintf("%d %s, %d, %s, %s", requestData.Address.StreetNumber, requestData.Address.Route, requestData.Address.PostalCode, requestData.Address.Locality, requestData.Address.AdministrativeAreaLevel1)
	categoryUUID, err := uuid.Parse(requestData.Category_ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid UUID format for CategoryID",
		})
	}
	instructor := models.Instructor{
		AuthId:        requestData.AuthID,
		DOB:           requestData.DOB,
		Address:       address,
		Phone:         requestData.Phone,
		Image_url:     requestData.ImageURL,
		Email:         requestData.Email,
		ID_type:       requestData.ID_type,
		IDcard_number: requestData.IDcard_number,
		IDcard_url:    requestData.IDcard_url,
		Selfie_url:    requestData.Selfie_url,
		Category_ID:   categoryUUID,
	}

	if err := validate.Struct(instructor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}

	if err := models.DB.Create(&instructor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

type InstructorExpertiseDTO struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type InstructorListDTO struct {
	Id           uuid.UUID                `json:"id"`
	AuthId       string                   `json:"authid"`
	DOB          string                   `json:"dob"`
	Address      string                   `json:"address"`
	Phone        string                   `json:"phone"`
	Image_url    string                   `json:"image_url"`
	Expertises   []InstructorExpertiseDTO `json:"expertises"`
	Rate         float32                  `json:"rate"`
	JumlahReview int                      `json:"jumlah_review"`
}

type InstructorRatingDTO struct {
	Rate float32 `json:"rate"`
}

func InstructorList(c *fiber.Ctx) error {

	searchAddress := c.Query("address")

	var instructors []models.Instructor

	if searchAddress != "" {
		if err := models.DB.Preload("Expertises").Where("address LIKE ?", "%"+searchAddress+"%").Find(&instructors).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
	} else {
		if err := models.DB.Preload("Expertises").Find(&instructors).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
	}

	var instructorsDTO []InstructorListDTO
	for _, instructor := range instructors {

		var result struct {
			InstructorID string  `gorm:"column:instructor_id"`
			JumlahReview int     `gorm:"column:jumlah_review"`
			TotalRate    float32 `gorm:"column:total_rate"`
			RataRataRate float32 `gorm:"column:rata_rata_rate"`
		}
		if err := models.DB.
			Table("ratings").
			Select("instructor_id, COUNT(*) AS jumlah_review, SUM(rate) AS total_rate, ROUND(SUM(rate) / COUNT(*), 2) AS rata_rata_rate").
			Where("instructor_id = ?", instructor.Id).
			Group("instructor_id").
			Scan(&result).Error; err != nil {
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

		instructorDTO := InstructorListDTO{
			Id:           instructor.Id,
			AuthId:       instructor.AuthId,
			DOB:          instructor.DOB,
			Address:      instructor.Address,
			Phone:        instructor.Phone,
			Image_url:    instructor.Image_url,
			Expertises:   expertisesDTO,
			Rate:         result.RataRataRate,
			JumlahReview: result.JumlahReview,
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

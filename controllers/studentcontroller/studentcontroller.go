package studentcontroller

import (
	"fmt"
	"net/http"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

// user
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
		Phone    string `json:"phone"`
		ImageURL string `json:"image_url"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var address = fmt.Sprintf("%d %s, %d, %s, %s", requestData.Address.StreetNumber, requestData.Address.Route, requestData.Address.PostalCode, requestData.Address.Locality, requestData.Address.AdministrativeAreaLevel1)

	student := models.Student{
		AuthId:    requestData.AuthID,
		DOB:       requestData.DOB,
		Address:   address,
		Phone:     requestData.Phone,
		Image_url: requestData.ImageURL,
	}

	if err := validate.Struct(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}

	if err := models.DB.Create(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

func UserProfile(c *fiber.Ctx) error {

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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":        student.Id,
		"authid":    student.AuthId,
		"dob":       student.DOB,
		"address":   student.Address,
		"phone":     student.Phone,
		"image_url": student.Image_url,
	})

}

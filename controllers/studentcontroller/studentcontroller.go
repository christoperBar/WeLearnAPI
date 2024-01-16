package studentcontroller

import (
	"fmt"
	"net/http"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	fbAdmin "github.com/christoperBar/WeLearnAPI/config/auth"
)

var validate = validator.New()

// user
func Register(c *fiber.Ctx) error {

	var requestData struct {
		Username string `json:"username"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Password string `json:"password"`
		DOB      string `json:"DOB"`
		Address  struct {
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

	createdStudent := models.Student{
		AuthId:    newUser.UID,
		DOB:       requestData.DOB,
		Address:   address,
		Phone:     requestData.Phone,
		Image_url: imageUrl,
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
		"authid":    student.AuthId,
		"dob":       student.DOB,
		"address":   student.Address,
		"phone":     student.Phone,
		"image_url": student.Image_url,
	})

}

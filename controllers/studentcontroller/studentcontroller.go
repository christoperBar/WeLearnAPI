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

func SayembaraList(c *fiber.Ctx) error {

	id := c.Params("id")

	var sayembaras []models.Sayembara
	var student models.Student

	if err := models.DB.First(&student, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Student not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	if err := models.DB.Preload("Student").Preload("Category").Preload("Expertises").Where("student_id = ?", student.Id).Find(&sayembaras).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching sayembara details",
		})
	}

	return c.Status(fiber.StatusOK).JSON(sayembaras)

}

func SayembaraDetail(c *fiber.Ctx) error {

	id := c.Params("sayembaraid")

	var sayembara models.Sayembara

	if err := models.DB.First(&sayembara, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Sayembara not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	if err := models.DB.Preload("Student").Preload("Category").Preload("Expertises").Where("id = ?", sayembara.Id).Find(&sayembara).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching sayembara details",
		})
	}

	return c.Status(fiber.StatusOK).JSON(sayembara)

}

// todo

// func CreateSayembara(c *fiber.Ctx) error {

// 	id := c.Params("id")
// 	studentID, err := strconv.ParseInt(id, 10, 64)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid student ID",
// 		})
// 	}

// 	var requestDatas struct {
// 		Title       string   `json:"title"`
// 		Description string   `json:"description"`
// 		Phone       string   `json:"phone"`
// 		ImageURL    string   `json:"image_url"`
// 		CategoryID  string   `json:"category_id"`
// 		Expertises  []string `json:"expertises_id"`
// 		Budget      struct {
// 			Min float64 `json:"min"`
// 			Max float64 `json:"max"`
// 		} `json:"budget"`
// 	}

// 	categoryID, err := strconv.ParseInt(requestDatas.CategoryID, 10, 64)

// 	if err := c.BodyParser(&requestDatas); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": err.Error(),
// 		})
// 	}

// 	var expertises []models.Expertise
//     for _, expertiseID := range requestDatas.ExpertisesIDs {
//         expertise, err := GetExpertiseByID(expertiseID)
//         if err != nil {
//             // Handle the error
//         }
//         expertises = append(expertises, expertise)
//     }

// 	sayembara := models.Sayembara{
// 		Title:       requestDatas.Title,
// 		Description: requestDatas.Description,
// 		Budget_min:  float32(requestDatas.Budget.Min),
// 		Budget_max:  float32(requestDatas.Budget.Max),
// 		Image_url:   requestDatas.ImageURL,
// 		Student_ID:  studentID,
// 		Category_ID: categoryID,
// 		Expertises:  requestDatas.Expertises,
// 	}

// 	if err := validate.Struct(sayembara); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid Request",
// 		})
// 	}

// 	if err := models.DB.Create(&sayembara).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(sayembara)
// }

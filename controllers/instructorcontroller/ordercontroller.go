package instructorcontroller

import (
	"net/http"
	"strings"
	"time"

	"github.com/christoperBar/WeLearnAPI/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetailDTO struct {
	Id             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Method         []string  `json:"method"`
	Note           string    `json:"note"`
	Total          float32   `json:"total"`
	Schedule       time.Time `json:"schedule"`
	Expired        time.Time `json:"expired"`
	IsPaid         bool      `json:"ispaid"`
	Transaction_ID string    `json:"transaction_id"`
	Instructor_ID  uuid.UUID `json:"instructor_id"`
	Student_ID     uuid.UUID `json:"student_id"`
}

func OrderDetail(c *fiber.Ctx) error {

	orderid := c.Params("orderid")

	var order models.Order
	if err := models.DB.Where("id = ?", orderid).First(&order).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Data not Found",
		})
	}

	methods := strings.Split(order.Method, ",")

	orderDTO := OrderDetailDTO{
		Id:             order.Id,
		Title:          order.Title,
		Description:    order.Description,
		Method:         methods,
		Note:           order.Note,
		Total:          order.Total,
		Schedule:       order.Schedule,
		Expired:        order.Expired,
		IsPaid:         order.IsPaid,
		Transaction_ID: order.Transaction_ID,
		Instructor_ID:  order.Instructor_ID,
		Student_ID:     order.Student_ID,
	}

	return c.Status(fiber.StatusOK).JSON(orderDTO)
}

type CreateOrderRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Method      []string `json:"method" validate:"required"`
	Note        string   `json:"note"`
	Total       float32  `json:"total" validate:"required"`
	Schedule    string   `json:"schedule" validate:"required"`
	Expired     string   `json:"expired" validate:"required"`
	Student_ID  string   `json:"student_id" validate:"required"`
}

func CreateOrder(c *fiber.Ctx) error {

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

	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	var student models.Student
	if err := models.DB.Where("id = ?", request.Student_ID).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Data not Found",
			})
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Data not Found",
		})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  err.Error(),
		})
	}

	studentUUID, err := uuid.Parse(request.Student_ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid UUID format for CategoryID",
		})
	}

	scheduleTime, err := time.Parse("2006-01-02 15:04:05", request.Schedule)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid format for schedule",
		})
	}

	expiredTime, err := time.Parse("2006-01-02 15:04:05", request.Expired)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid format for expired",
		})
	}
	method := strings.Join(request.Method, ", ")

	order := models.Order{
		Title:         request.Title,
		Description:   request.Description,
		Method:        method,
		Note:          request.Note,
		Total:         request.Total,
		Schedule:      scheduleTime,
		Expired:       expiredTime,
		Instructor_ID: instructor.Id,
		Student_ID:    studentUUID,
	}
	if err := models.DB.Create(&order).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create order",
		})
	}
	return c.Status(fiber.StatusOK).Send(nil)
}

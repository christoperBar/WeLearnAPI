package models

import "github.com/google/uuid"

type Student struct {
	Id          string    `gorm:"type:varchar(130)" json:"id"`
	DOB         string    `gorm:"type:date" json:"dob" validate:"required"`
	Address     string    `gorm:"type:varchar(200)" json:"address" validate:"required"`
	Phone       string    `gorm:"type:varchar(15)" json:"phone" validate:"required"`
	Image_url   string    `gorm:"type:varchar(200)" json:"image_url"`
	Method      string    `gorm:"type:SET('Offline','Online')" json:"method" validate:"required"`
	Category_ID uuid.UUID `json:"category_id" validate:"required"`

	Ratings    []Rating    `gorm:"foreignKey:Student_ID" json:"rating"`
	Sayembaras []Sayembara `gorm:"foreignKey:Student_ID" json:"sayembaras"`
	Category   Category    `gorm:"foreignKey:Category_ID" json:"category"`
}

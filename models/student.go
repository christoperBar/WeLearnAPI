package models

import "github.com/google/uuid"

type Student struct {
	Id          uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	AuthId      string    `gorm:"type:varchar(45)" json:"authid" validate:"required"`
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

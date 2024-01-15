package models

import "github.com/google/uuid"

type Student struct {
	Id        uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	AuthId    string    `gorm:"type:varchar(45)" json:"authid" validate:"required"`
	DOB       string    `gorm:"type:date" json:"dob" validate:"required"`
	Address   string    `gorm:"type:varchar(200)" json:"address" validate:"required"`
	Phone     string    `gorm:"type:varchar(13)" json:"phone" validate:"required"`
	Image_url string    `gorm:"type:varchar(200)" json:"image_url"`

	Sayembaras []Sayembara `gorm:"foreignKey:Student_ID" json:"sayembaras"`
}

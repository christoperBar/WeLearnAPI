package models

import "github.com/google/uuid"

type Instructor struct {
	Id            uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	AuthId        string    `gorm:"type:varchar(45)" json:"authid" validate:"required"`
	DOB           string    `gorm:"type:date" json:"dob" validate:"required"`
	Address       string    `gorm:"type:varchar(200)" json:"address" validate:"required"`
	Email         string    `gorm:"type:varchar(200)" json:"email" validate:"required"`
	Phone         string    `gorm:"type:varchar(13)" json:"phone" validate:"required"`
	Image_url     string    `gorm:"type:varchar(200)" json:"image_url"`
	ID_type       string    `gorm:"type:ENUM('ID Card', 'Passport', 'Student Card') DEFAULT 'ID Card'" json:"id_type" validate:"required"`
	IDcard_number string    `gorm:"type:varchar(100)" json:"idcard_number" validate:"required"`
	IDcard_url    string    `gorm:"type:varchar(200)" json:"idcard_url" validate:"required"`
	Selfie_url    string    `gorm:"type:varchar(200)" json:"selfie_url" validate:"required"`
	Category_ID   uuid.UUID `json:"category_id" validate:"required"`

	Category       Category        `gorm:"foreignKey:Category_ID" json:"category"`
	Ratings        []Rating        `gorm:"foreignKey:Instructor_ID" json:"rating"`
	Lessons        []Lesson        `gorm:"foreignKey:Instructor_ID" json:"lesson"`
	Learning_paths []Learning_path `gorm:"foreignKey:Instructor_ID" json:"learning_path"`
	Expertises     []Expertise     `gorm:"many2many:instrucor_expertises;" json:"expertises"`
}

package models

import (
	"github.com/google/uuid"
)

type Lesson struct {
	Id            uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	Title         string    `gorm:"type:varchar(50)" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	Price         float32   `gorm:"type:decimal(10,2)" json:"price"`
	Tags          string    `gorm:"type:text" json:"tags"`
	Method        string    `gorm:"type:SET('Offline','Online')" json:"method"`
	Image_url     string    `gorm:"type:varchar(200)" json:"image_url"`
	Instructor_ID uuid.UUID `json:"instructor_id"`
	Category_ID   uuid.UUID `json:"category_id"`

	Learning_path []Learning_path `gorm:"many2many:lesson_learningpath" json:"learning_path"`

	Instructor Instructor `gorm:"foreignKey:Instructor_ID" json:"instructor"`
	Category   Category   `gorm:"foreignKey:Category_ID" json:"category"`
}

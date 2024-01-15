package models

import "github.com/google/uuid"

type Learning_path struct {
	Id          uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	Title       string    `gorm:"type:varchar(50)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`

	Lessons    []Lesson    `gorm:"many2many:lesson_learningpath" json:"lesson"`
	Expertises []Expertise `gorm:"many2many:sayembara_expertises;" json:"expertises"`
}

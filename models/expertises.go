package models

import "github.com/google/uuid"

type Expertise struct {
	Id   uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	Name string    `gorm:"type:varchar(50)" json:"name"`

	Instructors []Instructor `gorm:"many2many:instrucor_expertises;" json:"instructors"`
	Sayembaras  []Sayembara  `gorm:"many2many:sayembara_expertises;" json:"sayembaras"`
}

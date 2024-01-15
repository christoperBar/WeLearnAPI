package models

import "github.com/google/uuid"

type Category struct {
	Id   uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	Name string    `gorm:"type:varchar(50)" json:"name"`

	Sayembaras []Sayembara `gorm:"foreignKey:Category_ID" json:"sayembaras"`
	Lessons    []Lesson    `gorm:"foreignKey:Category_ID" json:"lesson"`
}

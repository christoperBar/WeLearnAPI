package models

import "github.com/google/uuid"

type Rating struct {
	Id            uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	Rate          float32   `gorm:"type:decimal(10,2)" json:"rate"`
	Review        string    `gorm:"type:text" json:"review"`
	Instructor_ID uuid.UUID `json:"instructor_id"`
	Student_ID    uuid.UUID `json:"student_id"`

	Instructor Instructor `gorm:"foreignKey:Instructor_ID" json:"instructor"`
	Student    Student    `gorm:"foreignKey:Student_ID" json:"student"`
}

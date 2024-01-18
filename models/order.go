package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id             uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	Title          string    `gorm:"type:varchar(50)" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	Method         string    `gorm:"type:SET('Offline','Online')" json:"method"`
	Note           string    `gorm:"type:text" json:"note"`
	Total          float32   `gorm:"type:decimal(10,2)" json:"total"`
	Schedule       time.Time `gorm:"type:datetime" json:"schedule"`
	Expired        time.Time `gorm:"type:datetime" json:"expired"`
	IsPaid         bool      `gorm:"default:false" json:"ispaid"`
	Transaction_ID string    `gorm:"type:varchar(36)" json:"transaction_id"`
	Instructor_ID  uuid.UUID `json:"instructor_id"`
	Student_ID     uuid.UUID `json:"student_id"`

	Instructor Instructor `gorm:"foreignKey:Instructor_ID" json:"instructor"`
	Student    Student    `gorm:"foreignKey:Student_ID" json:"student"`
}

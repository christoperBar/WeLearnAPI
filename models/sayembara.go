package models

import (
	"time"

	"github.com/google/uuid"
)

type Sayembara struct {
	Id          uuid.UUID `gorm:"type:varchar(36);default:(UUID());primary_key;" json:"id"`
	Title       string    `gorm:"type:varchar(50)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Budget_min  float32   `gorm:"type:decimal(10,2)" json:"budget_min"`
	Budget_max  float32   `gorm:"type:decimal(10,2)" json:"budget_max"`
	Image_url   string    `gorm:"type:varchar(200)" json:"image_url"`
	Status      string    `gorm:"type:ENUM('On Going', 'Closed') DEFAULT 'On Going'" json:"status"`
	CreatedAt   time.Time `gorm:"type:timestamp" json:"createdat"`
	ClosedAt    time.Time `gorm:"type:datetime" json:"closedat"`

	Student_ID  uuid.UUID `json:"student_id"`
	Category_ID uuid.UUID `json:"category_id"`

	Expertises []Expertise `gorm:"many2many:sayembara_expertises;" json:"expertises"`

	Student  Student  `gorm:"foreignKey:Student_ID" json:"student"`
	Category Category `gorm:"foreignKey:Category_ID" json:"category"`
}

package models

type Sayembara struct {
	Id          int64   `gorm:"primaryKey" json:"id"`
	Title       string  `gorm:"type:varchar(50)" json:"title"`
	Description string  `gorm:"type:text" json:"description"`
	Budget_min  float32 `gorm:"type:decimal(10,2)" json:"budget_min"`
	Budget_max  float32 `gorm:"type:decimal(10,2)" json:"budget_max"`
	Image_url   string  `gorm:"type:varchar(200)" json:"image_url"`
	Student_ID  int64   `json:"student_id"`
	Category_ID int64   `json:"category_id"`

	Expertises []Expertise `gorm:"many2many:sayembara_expertises;" json:"expertises"`

	Student  Student  `gorm:"foreignKey:Student_ID" json:"student"`
	Category Category `gorm:"foreignKey:Category_ID" json:"category"`
}

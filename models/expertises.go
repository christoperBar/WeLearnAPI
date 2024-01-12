package models

type Expertise struct {
	Id   int64  `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(50)" json:"name"`

	Instructors []Instructor `gorm:"many2many:instrucor_expertises;" json:"instructors"`
	Sayembaras  []Sayembara  `gorm:"many2many:sayembara_expertises;" json:"sayembaras"`
}

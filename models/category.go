package models

type Category struct {
	Id   int64  `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(50)" json:"name"`

	Sayembaras []Sayembara `gorm:"foreignKey:Category_ID" json:"sayembaras"`
	Lessons    []Lesson    `gorm:"foreignKey:Category_ID" json:"lesson"`
}

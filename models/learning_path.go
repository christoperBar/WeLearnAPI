package models

type Learning_path struct {
	Id          int64  `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(50)" json:"title"`
	Description string `gorm:"type:text" json:"description"`

	Lessons []Lesson `gorm:"foreignKey:Learning_path_ID" json:"lesson"`
}

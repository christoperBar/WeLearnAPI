package models

type Lesson struct {
	Id               int64    `gorm:"primaryKey" json:"id"`
	Title            string   `gorm:"type:varchar(50)" json:"title"`
	Description      string   `gorm:"type:text" json:"description"`
	Price            float32  `gorm:"type:decimal(10,2)" json:"price"`
	Tags             []string `gorm:"type:text" json:"tags"`
	Method           []string `gorm:"type:text" json:"method"`
	Image_url        string   `gorm:"type:varchar(200)" json:"image_url"`
	Instructor_ID    int64    `json:"instructor_id"`
	Category_ID      int64    `json:"category_id"`
	Learning_path_ID int64    `json:"leaarning_path_id"`

	Instructor    Instructor    `gorm:"foreignKey:Instructor_ID" json:"-"`
	Category      Category      `gorm:"foreignKey:Category_ID" json:"-"`
	Learning_path Learning_path `gorm:"foreignKey:Learning_path_ID" json:"-"`
}

package models

type Instructor struct {
	Id            int64  `gorm:"primaryKey" json:"id"`
	AuthId        string `gorm:"type:varchar(45)" json:"authid"`
	DOB           string `gorm:"type:date" json:"dob"`
	Address       string `gorm:"type:varchar(200)" json:"address"`
	Email         string `gorm:"type:varchar(200)" json:"email"`
	Phone         string `gorm:"type:varchar(13)" json:"phone"`
	Image_url     string `gorm:"type:varchar(200)" json:"image_url"`
	ID_type       string `gorm:"type:ENUM('ID Card', 'Passport', 'Student Card') DEFAULT 'ID Card'" json:"id_type"`
	IDcard_number string `gorm:"type:varchar(100)" json:"idcard_number"`
	IDcard_url    string `gorm:"type:varchar(200)" json:"idcard_url"`
	Selfie_url    string `gorm:"type:varchar(200)" json:"selfie_url"`

	Lessons    []Lesson    `gorm:"foreignKey:Instructor_ID" json:"lesson"`
	Expertises []Expertise `gorm:"many2many:instrucor_expertises;" json:"expertises"`
}

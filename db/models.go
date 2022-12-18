package db

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	*gorm.Model
	Title         string
	Url           string
	Author        string `gorm:"type:varchar(255)"`
	Text          string
	Date          time.Time
	MassMediaName string
	ImgUrl        string
}

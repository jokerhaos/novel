package models

import "time"

type Book struct {
	ID         uint      `gorm:"primary_key"`
	Name       string    `gorm:"column:name"`
	Author     string    `gorm:"column:author"`
	Describe   string    `gorm:"column:describe"`
	Href       string    `gorm:"column:href"`
	Status     int       `gorm:"column:status"`
	Lock       int       `gorm:"column:lock"`
	Type       string    `gorm:"column:type"`
	Down       string    `gorm:"column:down"`
	Chapter    string    `gorm:"column:chapter"`
	NewChapter string    `gorm:"column:new_chapter"`
	Image      string    `gorm:"column:image"`
	F          string    `gorm:"column:f"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
	Loeva      string
}

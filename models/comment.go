package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserId  uint
	User    *User
	PhotoId uint `json:"photo_id"`
	Photo   *Photo
	Message string `gorm:"not null" json:"message" form:"message" valid:"required~Message must not be empty"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)

	if errCreate != nil {
		err = errCreate
		return
	}
	err = nil
	return
}

func (c *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errUpdate := govalidator.ValidateStruct(c)

	if errUpdate != nil {
		err = errUpdate
		return
	}

	err = nil
	return
}

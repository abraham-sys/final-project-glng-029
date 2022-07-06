package models

import (
	"final-project/helpers"
	"strconv"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null; uniqueIndex" json:"username" form:"username" valid:"required~Username must not be empty"`
	Email    string `gorm:"not null; uniqueIndex" json:"email" form:"email" valid:"required~Email must not be empty, email~Invalid email format"`
	Password string `gorm:"not null" json:"password" form:"password" valid:"required~Password must not be empty, minstringlength(6)~Password must have a minimum length of 6 characters"`
	Age      int    `gorm:"not null" json:"age" form:"age" valid:"required~Age must not be empty, minimum~Minimum requirement for age is at least must be 8"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	govalidator.TagMap["minimum"] = govalidator.Validator(func(num string) bool {
		intNum, _ := strconv.Atoi(num)
		return intNum > 8
	})
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPassword(u.Password)

	err = nil
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	govalidator.TagMap["minimum"] = govalidator.Validator(func(num string) bool {
		intNum, _ := strconv.Atoi(num)
		return intNum > 8
	})

	_, errUpdate := govalidator.ValidateStruct(u)

	if errUpdate != nil {
		err = errUpdate
		return
	}

	err = nil
	return
}

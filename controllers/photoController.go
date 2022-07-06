package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PhotoData struct {
	ID        uint
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserId    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      UserData  `json:"User"`
}

func CreatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Photo := models.Photo{}

	if contentType == AppJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserId = uint(userData["id"].(float64))

	err := db.Debug().Create(&Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"created_at": Photo.CreatedAt,
	})
}

func UpdatePhoto(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	Photo := models.Photo{}
	photoId, _ := strconv.Atoi(c.Param("photoId"))

	db.First(&Photo, photoId)
	if contentType == AppJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	err := db.Model(&Photo).Where("id = ?", photoId).Updates(models.Photo{Title: Photo.Title, Caption: Photo.Caption, PhotoUrl: Photo.PhotoUrl}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"updated_at": Photo.UpdatedAt,
	})
}

func DeletePhoto(c *gin.Context) {
	db := database.GetDB()

	photoId, _ := strconv.Atoi(c.Param("photoId"))

	Photo := models.Photo{}

	err := db.Delete(&Photo, photoId).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}

func GetAllPhotos(c *gin.Context) {
	db := database.GetDB()

	Photo := PhotoData{}

	PhotoDatas := []PhotoData{}

	rows, err := db.Debug().Table("photos").Select("photos.id, photos.title, photos.caption, photos.photo_url, photos.user_id, photos.created_at, photos.updated_at, users.email, users.username").Joins("left join users on users.id = photos.user_id").Where("photos.deleted_at is null").Rows()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for rows.Next() {
		err := rows.Scan(
			&Photo.ID,
			&Photo.Title,
			&Photo.Caption,
			&Photo.PhotoUrl,
			&Photo.UserId,
			&Photo.CreatedAt,
			&Photo.UpdatedAt,
			&Photo.User.Email,
			&Photo.User.Username,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		PhotoDatas = append(PhotoDatas, Photo)
	}

	c.JSON(http.StatusOK, PhotoDatas)
}

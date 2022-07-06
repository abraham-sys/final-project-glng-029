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

type UserDataSocialMedia struct {
	ID              string `json:"id"`
	Username        string `json:"username"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type SocialMediaData struct {
	ID             uint                `json:"id"`
	Name           string              `json:"name"`
	SocialMediaUrl string              `json:"social_media_url"`
	UserId         uint                `json:"user_id"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	User           UserDataSocialMedia `json:"User"`
}

func CreateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	SocialMedia := models.SocialMedia{}

	if contentType == AppJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserId = uint(userData["id"].(float64))

	err := db.Debug().Create(&SocialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":               SocialMedia.ID,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"created_at":       SocialMedia.CreatedAt,
	})
}

func UpdateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}
	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))

	db.First(&SocialMedia, socialMediaId)
	if contentType == AppJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	err := db.Model(&SocialMedia).Where("id = ?", socialMediaId).Updates(models.SocialMedia{Name: SocialMedia.Name, SocialMediaUrl: SocialMedia.SocialMediaUrl}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               SocialMedia.ID,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"updated_at":       SocialMedia.UpdatedAt,
	})
}

func DeleteSocialMedia(c *gin.Context) {
	db := database.GetDB()

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))

	SocialMedia := models.SocialMedia{}

	err := db.Delete(&SocialMedia, socialMediaId).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}

func GetAllSocialMedias(c *gin.Context) {
	db := database.GetDB()

	SocialMedia := SocialMediaData{}

	SocialMediaDatas := []SocialMediaData{}

	rows, err := db.Debug().Table("social_media").Select("social_media.id, social_media.name, social_media.social_media_url, social_media.user_id, social_media.created_at, social_media.updated_at, users.id, users.username").Joins("left join users on users.id = social_media.user_id").Where("social_media.deleted_at is null").Rows()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for rows.Next() {
		err := rows.Scan(
			&SocialMedia.ID,
			&SocialMedia.Name,
			&SocialMedia.SocialMediaUrl,
			&SocialMedia.UserId,
			&SocialMedia.CreatedAt,
			&SocialMedia.UpdatedAt,
			&SocialMedia.User.ID,
			&SocialMedia.User.Username,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		SocialMediaDatas = append(SocialMediaDatas, SocialMedia)
	}

	c.JSON(http.StatusOK, gin.H{
		"social_medias": SocialMediaDatas,
	})
}

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

type UserDataComment struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PhotoDataComment struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   uint   `json:"user_id"`
}

type CommentData struct {
	ID        uint             `json:"id"`
	Message   string           `json:"message"`
	PhotoId   string           `json:"photo_id"`
	UserId    uint             `json:"user_id"`
	UpdatedAt time.Time        `json:"updated_at"`
	CreatedAt time.Time        `json:"created_at"`
	User      UserDataComment  `json:"User"`
	Photo     PhotoDataComment `json:"Photo"`
}

func CreateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := models.Comment{}

	if contentType == AppJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserId = uint(userData["id"].(float64))

	err := db.Debug().Create(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoId,
		"user_id":    Comment.UserId,
		"created_at": Comment.CreatedAt,
	})
}

func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	Comment := models.Comment{}
	commentId, _ := strconv.Atoi(c.Param("commentId"))

	db.First(&Comment, commentId)
	if contentType == AppJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	err := db.Model(&Comment).Where("id = ?", commentId).Updates(models.Comment{Message: Comment.Message}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoId,
		"user_id":    Comment.UserId,
		"updated_at": Comment.UpdatedAt,
	})
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()

	commentId, _ := strconv.Atoi(c.Param("commentId"))

	Comment := models.Comment{}

	err := db.Delete(&Comment, commentId).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}

func GetAllComments(c *gin.Context) {
	db := database.GetDB()

	Comment := CommentData{}

	CommentDatas := []CommentData{}

	rows, err := db.Debug().Table("comments").Select("comments.id, comments.message, comments.photo_id, comments.user_id, comments.created_at, comments.updated_at, users.id, users.age, users.username, photos.id, photos.title, photos.caption, photos.photo_url, photos.user_id").Joins("left join users on users.id = comments.user_id").Joins("left join photos on photos.id = comments.photo_id").Where("comments.deleted_at is null").Rows()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for rows.Next() {
		err := rows.Scan(
			&Comment.ID,
			&Comment.Message,
			&Comment.PhotoId,
			&Comment.UserId,
			&Comment.CreatedAt,
			&Comment.UpdatedAt,
			&Comment.User.ID,
			&Comment.User.Email,
			&Comment.User.Username,
			&Comment.Photo.ID,
			&Comment.Photo.Title,
			&Comment.Photo.Caption,
			&Comment.Photo.PhotoUrl,
			&Comment.Photo.UserId,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		CommentDatas = append(CommentDatas, Comment)
	}

	c.JSON(http.StatusOK, CommentDatas)
}

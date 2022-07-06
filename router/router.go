package router

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
		userRouter.PUT("/:userId", middlewares.Authenticate(), controllers.UpdateUser)
		userRouter.DELETE("/:userId", middlewares.Authenticate(), controllers.DeleteUser)
	}

	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authenticate())
		photoRouter.POST("/", controllers.CreatePhoto)
		photoRouter.GET("/", controllers.GetAllPhotos)
		photoRouter.PUT("/:photoId", middlewares.AuthorizePhoto(), controllers.UpdatePhoto)
		photoRouter.DELETE("/:photoId", middlewares.AuthorizePhoto(), controllers.DeletePhoto)
	}

	commentRouter := r.Group("/comments")
	{
		commentRouter.Use(middlewares.Authenticate())
		commentRouter.POST("/", controllers.CreateComment)
		commentRouter.GET("/", controllers.GetAllComments)
		commentRouter.PUT("/:commentId", middlewares.AuthorizeComment(), controllers.UpdateComment)
		commentRouter.DELETE("/:commentId", middlewares.AuthorizeComment(), controllers.DeleteComment)
	}

	socialMediaRouter := r.Group("/socialmedias")
	{
		socialMediaRouter.Use(middlewares.Authenticate())
		socialMediaRouter.POST("/", controllers.CreateSocialMedia)
		socialMediaRouter.GET("/", controllers.GetAllSocialMedias)
		socialMediaRouter.PUT("/:socialMediaId", middlewares.AuthorizeSocialMedia(), controllers.UpdateSocialMedia)
		socialMediaRouter.DELETE("/:socialMediaId", middlewares.AuthorizeSocialMedia(), controllers.DeleteSocialMedia)
	}

	return r
}

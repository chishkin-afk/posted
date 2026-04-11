package handlers

import (
	"github.com/chishkin-afk/posted/http-gateway/pkg/errs"
	"github.com/gin-gonic/gin"
)

const (
	ModeDev   = "dev"
	ModeProd  = "prod"
	ModeLocal = "local"
)

func New(env string, authService authService, postsService postsService) (*gin.Engine, error) {
	var router *gin.Engine
	switch env {
	case ModeDev, ModeLocal:
		router = gin.Default()
	case ModeProd:
		router = gin.New()
		router.Use(gin.Recovery())
	default:
		return nil, errs.ErrInvalidEnvironment
	}

	handlers := handlers{
		authService:  authService,
		postsService: postsService,
	}

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Auth block
			v1.POST("/register", handlers.Register())
			v1.POST("/login", handlers.Login())
			v1.PATCH("/user", handlers.UpdateUser())
			v1.DELETE("/user", handlers.DeleteUser())
			v1.GET("/user/:id", handlers.GetUserByID())
			v1.GET("/user", handlers.GetUserSelf())

			// Posts block
			v1.POST("/post", handlers.CreatePost())
			v1.PATCH("/post", handlers.UpdatePost())
			v1.DELETE("/post/:id", handlers.DeletePost())
			v1.GET("/post/:id", handlers.GetPostByID())
			v1.GET("/posts", handlers.GetSelfPosts()) // ?page=&limit=
		}
	}

	return router, nil
}

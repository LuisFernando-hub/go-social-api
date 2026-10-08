package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LuisFernando-hub/go-social-api/internal/auth"
	"github.com/LuisFernando-hub/go-social-api/internal/config"
	"github.com/LuisFernando-hub/go-social-api/internal/database"
	authMiddleware "github.com/LuisFernando-hub/go-social-api/internal/middleware/auth"
	"github.com/LuisFernando-hub/go-social-api/internal/models"
	"github.com/LuisFernando-hub/go-social-api/internal/post"
	postinteraction "github.com/LuisFernando-hub/go-social-api/internal/post_interaction"
	"github.com/LuisFernando-hub/go-social-api/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	var cfg *config.Config
	var err error

	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//Auto-Migrations
	if err := db.AutoMigrate(&models.Post{}, &models.User{}, &models.PostInteraction{}); err != nil {
		log.Fatal(err)
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	postRepository := post.NewRepository(db)
	postService := post.NewService(postRepository, userService)
	postHandler := post.NewHandler(postService)

	postInteractionRepository := postinteraction.NewRepository(db)
	postInteractionService := postinteraction.NewService(postInteractionRepository, postService)
	postInteractionHandler := postinteraction.NewHandler(postInteractionService)


	authService := auth.NewService(userRepository)
	authHandler := auth.NewHandler(authService, cfg)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		fmt.Println("Hello World!")
		log.Printf("Server running on :%s", port)

		c.JSON(200, gin.H{
			"message": "Server running",
		})
	})

	//Auth
	router.POST("/api/v1/login", authHandler.Login)


	usersGroup := router.Group("/api/v1/register")
	{
		usersGroup.POST("", userHandler.Create)
	}
	
	// -- Authorization Auth -- //
	postsGroup := router.Group("/api/v1/posts")
	postsGroup.Use(authMiddleware.Middleware(cfg))
	{
		postsGroup.POST("", postHandler.Create)
		postsGroup.GET("", postHandler.List)

		postsGroup.POST("/interation", postInteractionHandler.CreateInteractionPost)
	}

	// ------------------------ //

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}

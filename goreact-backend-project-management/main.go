package main

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/controllers"
	"arfian-nurdavicky/project-management/database/seed"
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/repositories"
	"arfian-nurdavicky/project-management/routes"
	"arfian-nurdavicky/project-management/services"
)

func main() {
	// Load .env and connect to DB
	config.LoadEnv()
	config.ConnectDB()
err := config.DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}
	// Run seeder
	seed.SeedAdmin()

	// Inisialisasi Fiber
	app := fiber.New()

	// Ambil origins dari ENV dan split koma
	origins := os.Getenv("CORS_ORIGIN")
	if origins == "" {
		origins = "*" // fallback untuk testing
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(strings.Split(origins, ","), ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// ================== Init Layer ==================

	// User
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Board
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	// List
	listPosRepo := repositories.NewListPositionRepository()
	listRepo := repositories.NewListRepository()
	listService := services.NewListService(listRepo, boardRepo, listPosRepo)
	listController := controllers.NewListController(listService)

	// Card
	cardRepo := repositories.NewCardRepository()
	labelRepo := repositories.NewLabelRepository()
	cardService := services.NewCardService(cardRepo, listRepo, labelRepo, userRepo)

	// === Attachment (baru) ===
	attachmentRepo := repositories.NewAttachmentRepository(config.DB)
	attachmentService := services.NewAttachmentService(attachmentRepo, cardRepo, userRepo)

	// ambil folder upload dari ENV, fallback ke ./uploads/attachments
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads/attachments"
	}

	cardController := controllers.NewCardController(cardService, attachmentService, uploadDir)

	// di main.go
	app.Static("/files", "./uploads/attachments")

	// Comment
	commentRepo := repositories.NewCommentRepository()
	commentService := services.NewCommentService(commentRepo, cardRepo, userRepo)
	commentController := controllers.NewCommentController(commentService)

	// ================== Routing ==================
	routes.Setup(app,
		userController,
		boardController,
		listController,
		cardController,
		commentController)

	// Port dari ENV
	port := config.AppConfig.AppPort
	log.Printf("🚀 Server is running on http://localhost:%s\n", port)

	log.Fatal(app.Listen(":" + port))
}

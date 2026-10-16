package routes

import (
	"arfian-nurdavicky/project-management/controllers"
	"log"
	"os"

	"arfian-nurdavicky/project-management/middleware"
	"arfian-nurdavicky/project-management/utils"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App,
	uc *controllers.UserController,
	bc *controllers.BoardController,
	lc *controllers.ListController,
	cc *controllers.CardController,
	cmc *controllers.CommentController) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Public routes
	app.Post("/v1/auth/register", uc.Register)
	app.Post("/v1/auth/login", uc.Login)
	app.Get("/v1/swagger/*", swagger.HandlerDefault) // Swagger route

	// JWT Protected routes
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, err.Error())
		},
	}))

	// User routes
	userGroup := api.Group("/users")
	userGroup.Get("/", middleware.RoleMiddleware("admin"), uc.GetUsers)
	userGroup.Get("/page", middleware.RoleMiddleware("admin"), uc.GetUsersPaginate) //done
	userGroup.Get("/:id", uc.GetUser)                                               // done
	userGroup.Post("/", middleware.RoleMiddleware("admin"), uc.CreateUser)
	userGroup.Put("/:id", middleware.RoleMiddlewareWithException("admin"), uc.UpdateUser) //done
	userGroup.Delete("/:id", middleware.RoleMiddleware("admin"), uc.DeleteUser)           //done

	// Board routes
	boardGroup := api.Group("/boards")
	boardGroup.Post("/", bc.CreateBoard) //done
	boardGroup.Get("/all", middleware.RoleMiddleware("admin"), bc.GetBoards)
	boardGroup.Get("/", middleware.RoleMiddleware("admin"), bc.GetBoardsPaginate)
	//boardGroup.Get("/my", bc.GetUserBoards)
	boardGroup.Get("/my", bc.GetMyBoardsPaginate) //GetMyBoardsPaginate // progress
	boardGroup.Get("/:id", bc.GetBoard)
	boardGroup.Put("/:id", bc.UpdateBoard) //done
	boardGroup.Delete("/:id", bc.DeleteBoard)
	//boardGroup.Post("/:id/members", bc.AddBoardMember)
	boardGroup.Post("/:id/members", bc.AddBoardMembers) //Done
	//boardGroup.Delete("/:id/members", bc.RemoveBoardMember)
	boardGroup.Delete("/:id/members", bc.RemoveBoardMembers) //done
	boardGroup.Get("/:id/members", bc.GetBoardMembers)
	boardGroup.Get("/:board_id/lists", lc.GetBoardLists)
	boardGroup.Put("/:board_id/positions", lc.UpdateListPositions)

	// List routes
	listGroup := api.Group("/lists")
	listGroup.Post("/", lc.CreateList)
	listGroup.Put("/:id", lc.UpdateList)
	listGroup.Delete("/:id", lc.DeleteList)

	// Nested under boards
	boardsGroup := api.Group("/boards")
	boardsGroup.Get("/:board_id/lists", lc.GetBoardLists) //err
	boardsGroup.Put("/:board_id/positions", lc.UpdateListPositions)
	//{1f6d23b6-36ef-4abf-b97b-25ea8ea01043,29df6de0-2f89-4478-8ac4-51848c24c58f,a03a1ab0-bf17-47c3-aaa4-39c3654aa7ae}/

	// Card routes
	cardGroup := api.Group("/cards")
	cardGroup.Post("/", cc.CreateCard)      //
	cardGroup.Put("/:id", cc.UpdateCard)    //
	cardGroup.Get("/:id", cc.GetCardDetail) //
	cardGroup.Delete("/:id", cc.DeleteCard) //
	cardGroup.Post("/:id/labels", cc.AddCardLabel)
	cardGroup.Delete("/:id/labels", cc.RemoveCardLabel)

	cardGroup.Post("/:id/assignees", cc.AddCardAssignee)
	cardGroup.Delete("/:id/assignees", cc.RemoveCardAssignee)
	//card attachment
	cardGroup.Post("/:id/attachments", cc.UploadAttachment)                       // upload
	cardGroup.Get("/:id/attachments", cc.GetAttachments)                          // list
	cardGroup.Delete("/:card_id/attachments/:attachment_id", cc.DeleteAttachment) // delete

	// Nested under lists
	listsGroup := api.Group("/lists")
	listsGroup.Get("/:list_id/cards", cc.GetListCards)
	listsGroup.Put("/:list_id/positions", cc.UpdateCardPositions)

	// Comment routes
	commentGroup := api.Group("/comments")
	commentGroup.Post("/", cmc.CreateComment)
	commentGroup.Put("/:id", cmc.UpdateComment)
	commentGroup.Delete("/:id", cmc.DeleteComment)

	// Nested under cards
	cardsGroup := api.Group("/cards")
	cardsGroup.Get("/:card_id/comments", cmc.GetCardComments)
}

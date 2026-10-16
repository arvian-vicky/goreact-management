package controllers

import (
	"arfian-nurdavicky/project-management/models"
	"math"
	"strconv"

	"arfian-nurdavicky/project-management/services"
	"arfian-nurdavicky/project-management/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

// CreateBoard godoc
// @Summary Create a new board
// @Description Create a new board with the provided details.
// @Tags boards
// @Accept json
// @Produce json
// @Param board body models.Board true "Board data to create"
// @Success 200 {object} ApiResponse{data=models.Board} "Successfully created board"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /boards [post]
func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	board := new(models.Board)
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat Board", err.Error())
	}

	return utils.Success(ctx, "Board berhasil dibuat", board)
}

// GetBoards godoc
// @Summary Get all boards
// @Description Retrieve a list of all boards.
// @Tags boards
// @Produce json
// @Success 200 {object} ApiResponse{data=[]models.Board} "Successfully retrieved boards"
// @Failure 500 {object} ApiErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /boards [get]
func (c *BoardController) GetBoards(ctx *fiber.Ctx) error {
	boards, err := c.service.GetAll()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Mengambil Data", err.Error())
	}

	return utils.Success(ctx, "Data boards berhasil diambil", boards)
}

// GetBoardsPaginate godoc
// @Summary Get list of boards with pagination
// @Description Ambil daftar board dengan filter, sort, dan pagination
// @Tags boards
// @Accept json
// @Produce json
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah data per halaman (default: 10)"
// @Param filter query string false "Cari berdasarkan title"
// @Param sort query string false "Urutkan data, contoh: created_at desc"
// @Success 200 {object} utils.PaginatedResponse
// @Failure 404 {object} utils.PaginatedResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /boards [get]
func (c *BoardController) GetBoardsPaginate(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	boards, total, err := c.service.GetAllPaginate(filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Mengambil Data", err.Error()) //err 500
	}

	//var boardResp []models. BoardResponse
	//_ = copier.Copy(&boardResp, &boards)

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}

	if total == 0 {
		return utils.PaginatedNotFound(ctx, "Data board tidak ditemukan", boards, meta)
	}

	return utils.PaginatedSuccess(ctx, "Data board berhasil diambil", boards, meta)
}

// GetUserBoards godoc
// @Summary Get boards for current user
// @Description Retrieve boards where current user is a member
// @Tags boards
// @Produce json
// @Success 200 {object} ApiResponse{data=[]models.Board} "Successfully retrieved user boards"
// @Failure 500 {object} ApiErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /boards/my [get]
func (c *BoardController) GetUserBoards(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["pub_id"].(string)

	boards, err := c.service.GetAllByUser(userID)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Mengambil Data", err.Error())
	}

	return utils.Success(ctx, "Data boards pengguna berhasil diambil", boards)
}

func (c *BoardController) GetMyBoardsPaginate(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["pub_id"].(string)

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	boards, total, err := c.service.GetAllByUserPaginate(userID, filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Mengambil Data", err.Error()) //err 500
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}

	if total == 0 {
		return utils.PaginatedNotFound(ctx, "Data board tidak ditemukan", boards, meta)
	}

	return utils.PaginatedSuccess(ctx, "Data board berhasil diambil", boards, meta)
}

// GetBoard godoc
// @Summary Get board by ID
// @Description Get a board by its public ID
// @Tags boards
// @Produce json
// @Param id path string true "Board public ID"
// @Success 200 {object} ApiResponse{data=models.Board} "Successfully retrieved board"
// @Failure 404 {object} ApiErrorResponse "Board not found"
// @Security BearerAuth
// @Router /boards/{id} [get]
func (c *BoardController) GetBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	board, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board tidak ditemukan")
	}

	return utils.Success(ctx, "Data board berhasil diambil", board)
}

// UpdateBoard godoc
// @Summary Update board
// @Description Update an existing board
// @Tags boards
// @Accept json
// @Produce json
// @Param id path string true "Board public ID"
// @Param board body models.Board true "Updated board data"
// @Success 200 {object} ApiResponse{data=models.Board} "Successfully updated board"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /boards/{id} [put]
func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	existingBoard, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board tidak ditemukan")
	}

	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID

	if err := c.service.Update(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Board", err.Error())
	}

	return utils.Success(ctx, "Board berhasil diperbarui", board)
}

// DeleteBoard godoc
// @Summary Delete board
// @Description Delete a board by its public ID
// @Tags boards
// @Produce json
// @Param id path string true "Board public ID"
// @Success 200 {object} ApiResponse{data=string} "Board successfully deleted"
// @Failure 404 {object} ApiErrorResponse "Board not found"
// @Security BearerAuth
// @Router /boards/{id} [delete]
func (c *BoardController) DeleteBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	board, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board tidak ditemukan")
	}

	if err := c.service.Delete(uint(board.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Board", err.Error())
	}

	return utils.Success(ctx, "Board berhasil dihapus", publicID)
}

// AddBoardMember godoc
// @Summary Add member to board
// @Description Add a user to a board as a member
// @Tags boards
// @Accept json
// @Produce json
// @Param id path string true "Board public ID"
// @Param member body object{user_id=string} true "User to add"
// @Success 200 {object} ApiResponse "Member successfully added"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /boards/{id}/members [post]
func (c *BoardController) AddBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var body struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.AddMember(publicID, body.UserID); err != nil {
		return utils.BadRequest(ctx, "Gagal Menambahkan Member", err.Error())
	}

	return utils.Success(ctx, "Member berhasil ditambahkan", nil)
}

// AddBoardMembers godoc
// @Summary Add multiple members to board
// @Description Add multiple users to a board as members
// @Tags boards
// @Accept json
// @Produce json
// @Param id path string true "Board public ID"
// @Param members body []string true "Array of user public IDs to add"
// @Success 200 {object} ApiResponse "Members successfully added"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /boards/{id}/members [post]
func (c *BoardController) AddBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if len(userIDs) == 0 {
		return utils.BadRequest(ctx, "Data tidak valid", "Minimal satu user harus dipilih")
	}

	if err := c.service.AddMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Menambahkan Members", err.Error())
	}

	return utils.Success(ctx, "Members berhasil ditambahkan", nil)
}

// RemoveBoardMember godoc
// @Summary Remove member from board
// @Description Remove a user from a board
// @Tags boards
// @Accept json
// @Produce json
// @Param id path string true "Board public ID"
// @Param member body object{user_id=string} true "User to remove"
// @Success 200 {object} ApiResponse "Member successfully removed"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /boards/{id}/members [delete]
func (c *BoardController) RemoveBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var body struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.RemoveMember(publicID, body.UserID); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Member", err.Error())
	}

	return utils.Success(ctx, "Member berhasil dihapus", nil)
}

// RemoveBoardMembers godoc
// @Summary Remove multiple members from board
// @Description Remove multiple users from a board
// @Tags boards
// @Accept json
// @Produce json
// @Param id path string true "Board public ID"
// @Param members body []string true "Array of user public IDs to remove" example(["550e8400-e29b-41d4-a716-446655440000", "f47ac10b-58cc-0372-8567-0e02b2c3d479"])
// @Success 200 {object} utils.Response{data=interface{}} "Members successfully removed"
// @Failure 400 {object} utils.ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /boards/{id}/members [delete]
func (c *BoardController) RemoveBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if len(userIDs) == 0 {
		return utils.BadRequest(ctx, "Data tidak valid", "Minimal satu user harus dipilih")
	}

	if err := c.service.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Members", err.Error())
	}

	return utils.Success(ctx, "Members berhasil dihapus", nil)
}

func (c *BoardController) GetBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	members, err := c.service.GetMembers(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Members tidak ditemukan")
	}

	var userResp []models.UserResponse
	_ = copier.Copy(&userResp, &members)

	return utils.Success(ctx, "Data member board berhasil diambil", userResp)
}

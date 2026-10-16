package controllers

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4" // gunakan versi yang sesuai
	"github.com/google/uuid"

	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/services"
	"arfian-nurdavicky/project-management/utils"
)

type CardController struct {
	service           services.CardService
	attachmentService services.AttachmentService
	uploadDir         string
}

func NewCardController(s services.CardService, a services.AttachmentService, uploadDir string) *CardController {
	return &CardController{service: s, attachmentService: a, uploadDir: uploadDir}
}

// CreateCard godoc
// @Summary Create a new card
// @Description Create a new card in a list
// @Tags cards
// @Accept json
// @Produce json
// @Param card body models.Card true "Card data to create"
// @Success 200 {object} ApiResponse{data=models.Card} "Successfully created card"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /cards [post]
func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}

	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Position:    req.Position,
	}

	if err := c.service.Create(card, req.ListPublicID); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat Card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dibuat", card)
}

// GetListCards godoc
// @Summary Get cards in a list
// @Description Retrieve all cards in a list
// @Tags cards
// @Produce json
// @Param list_id path string true "List public ID"
// @Success 200 {object} ApiResponse{data=[]models.Card} "Successfully retrieved cards"
// @Failure 404 {object} ApiErrorResponse "List not found"
// @Security BearerAuth
// @Router /lists/{list_id}/cards [get]
func (c *CardController) GetListCards(ctx *fiber.Ctx) error {
	listID := ctx.Params("list_id")
	if _, err := uuid.Parse(listID); err != nil {
		return utils.BadRequest(ctx, "ID List tidak valid", err.Error())
	}

	cards, err := c.service.GetByListID(listID)
	if err != nil {
		return utils.NotFound(ctx, "Card tidak ditemukan")
	}

	return utils.Success(ctx, "Data card berhasil diambil", cards)
}

// UpdateCard godoc
// @Summary Update card
// @Description Update an existing card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "Card public ID"
// @Param card body models.Card true "Updated card data"
// @Success 200 {object} ApiResponse{data=models.Card} "Successfully updated card"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /cards/{id} [put]
func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	type UpdateCardRequest struct {
		ListPublicID string     `json:"list_id"`
		Title        string     `json:"title"`
		Description  string     `json:"description"`
		DueDate      *time.Time `json:"due_date"`
		Position     int        `json:"position"`
	}

	var req UpdateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Position:    req.Position,
		PublicID:    uuid.MustParse(publicID),
	}

	if err := c.service.Update(card, req.ListPublicID); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil diperbarui", card)
}

// DeleteCard godoc
// @Summary Delete card
// @Description Delete a card by its public ID
// @Tags cards
// @Produce json
// @Param id path string true "Card public ID"
// @Success 200 {object} ApiResponse{data=string} "Card successfully deleted"
// @Failure 404 {object} ApiErrorResponse "Card not found"
// @Security BearerAuth
// @Router /cards/{id} [delete]
func (c *CardController) DeleteCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	card, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Card tidak ditemukan")
	}

	if err := c.service.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dihapus", publicID)
}

// AddCardLabel godoc
// @Summary Add label to card
// @Description Add a label to a card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "Card public ID"
// @Param label body object{label_id=string} true "Label to add"
// @Success 200 {object} ApiResponse "Label added successfully"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /cards/{id}/labels [post]
func (c *CardController) AddCardLabel(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")

	var body struct {
		LabelID string `json:"label_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.AddLabel(cardID, body.LabelID); err != nil {
		return utils.BadRequest(ctx, "Gagal Menambahkan Label", err.Error())
	}

	return utils.Success(ctx, "Label berhasil ditambahkan", nil)
}

// RemoveCardLabel godoc
// @Summary Remove label from card
// @Description Remove a label from a card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "Card public ID"
// @Param label body object{label_id=string} true "Label to remove"
// @Success 200 {object} ApiResponse "Label removed successfully"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /cards/{id}/labels [delete]
func (c *CardController) RemoveCardLabel(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")

	var body struct {
		LabelID string `json:"label_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.RemoveLabel(cardID, body.LabelID); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Label", err.Error())
	}

	return utils.Success(ctx, "Label berhasil dihapus", nil)
}

// AddCardAssignee godoc
// @Summary Add assignee to card
// @Description Add a user to a card as assignee
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "Card public ID"
// @Param assignee body object{user_id=string} true "User to add"
// @Success 200 {object} ApiResponse "Assignee added successfully"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /cards/{id}/assignees [post]
func (c *CardController) AddCardAssignee(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")

	var body struct {
		UserID []string `json:"user_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.AddAssignees(cardID, body.UserID); err != nil {
		return utils.BadRequest(ctx, "Gagal Menambahkan Assignee", err.Error())
	}

	return utils.Success(ctx, "Assignee berhasil ditambahkan", nil)
}

func (cc *CardController) GetCardDetail(c *fiber.Ctx) error {
	cardPublicID := c.Params("id")

	// Ambil detail card dari service
	card, err := cc.service.GetCardDetail(cardPublicID)
	if err != nil {
		return utils.InternalServerError(c, err.Error(), "Error")
	}
	if card == nil {
		return utils.NotFound(c, "Card tidak ditemukan")
	}

	// Response sukses
	return utils.Success(c, "Detail card berhasil diambil", card)
}

// RemoveCardAssignee godoc
// @Summary Remove assignee from card
// @Description Remove a user from a card assignees
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "Card public ID"
// @Param assignee body object{user_id=string} true "User to remove"
// @Success 200 {object} ApiResponse "Assignee removed successfully"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /cards/{id}/assignees [delete]
func (c *CardController) RemoveCardAssignee(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")

	var body struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.RemoveAssignee(cardID, body.UserID); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Assignee", err.Error())
	}

	return utils.Success(ctx, "Assignee berhasil dihapus", nil)
}

// UpdateCardPositions godoc
// @Summary Update card positions in a list
// @Description Update the order of cards in a list
// @Tags cards
// @Accept json
// @Produce json
// @Param list_id path string true "List public ID"
// @Param positions body []string true "Array of card public IDs in order"
// @Success 200 {object} ApiResponse "Positions updated successfully"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /lists/{list_id}/positions [put]
func (c *CardController) UpdateCardPositions(ctx *fiber.Ctx) error {
	listID := ctx.Params("list_id")
	if _, err := uuid.Parse(listID); err != nil {
		return utils.BadRequest(ctx, "ID List tidak valid", err.Error())
	}

	// Buat struct sesuai request body
	var req struct {
		Positions []string `json:"positions"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.UpdatePositions(listID, req.Positions); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Posisi", err.Error())
	}

	return utils.Success(ctx, "Posisi card berhasil diperbarui", nil)
}

//card attachment

// POST /api/v1/cards/:id/attachments
func (cc *CardController) UploadAttachment(c *fiber.Ctx) error {
	cardPublicID := c.Params("id")
	if cardPublicID == "" {
		return utils.BadRequest(c, "card id is required", "Error")
	}

	// ambil user public_id dari JWT claims
	u := c.Locals("user")
	if u == nil {
		return utils.Unauthorized(c, "unauthorized")
	}
	token, ok := u.(*jwt.Token)
	if !ok {
		return utils.Unauthorized(c, "unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.Unauthorized(c, "unauthorized")
	}

	userPublicID, _ := claims["pub_id"].(string)
	if userPublicID == "" {
		// fallback: kadang claim bernama "sub" atau "user_id" -> sesuaikan jika perlu
		if sub, ok := claims["sub"].(string); ok && sub != "" {
			userPublicID = sub
		} else {
			return utils.Unauthorized(c, "user id claim not found")
		}
	}

	// ambil file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return utils.BadRequest(c, "file is required", "Error")
	}

	// validasi ekstensi / size (contoh simple)
	const maxSize = 20 << 20 // 20 MB
	if fileHeader.Size > maxSize {
		return utils.BadRequest(c, "file too large (max 20MB)", "Error")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".pdf": true, ".docx": true, ".xlsx": true}
	if !allowed[ext] {
		return utils.BadRequest(c, "file type not allowed", "Error")
	}

	// buat folder jika belum ada
	if err := os.MkdirAll(cc.uploadDir, os.ModePerm); err != nil {
		return utils.InternalServerError(c, "cannot create upload dir", "Error")
	}

	// simpan dengan nama unik
	dstFilename := uuid.New().String() + ext
	dstPath := filepath.Join(cc.uploadDir, dstFilename)

	src, err := fileHeader.Open()
	if err != nil {
		return utils.InternalServerError(c, "cannot open uploaded file", "Error")
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return utils.InternalServerError(c, "cannot create destination file", "Error")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return utils.InternalServerError(c, "cannot save file", "Error")
	}

	// simpan record via service
	attachment, err := cc.attachmentService.Create(cardPublicID, userPublicID, dstPath)
	if err != nil {
		// kalau gagal, hapus file yang sudah disimpan
		_ = os.Remove(dstPath)
		return utils.InternalServerError(c, err.Error(), "Error")
	}

	return c.Status(fiber.StatusCreated).JSON(attachment)
}

// GET /api/v1/cards/:id/attachments
func (cc *CardController) GetAttachments(c *fiber.Ctx) error {
	cardPublicID := c.Params("id")
	attachments, err := cc.attachmentService.GetByCardID(cardPublicID)
	if err != nil {
		return utils.InternalServerError(c, err.Error(), "Error")
	}
	return c.JSON(attachments)
}

// DELETE /api/v1/cards/:card_id/attachments/:attachment_id
func (cc *CardController) DeleteAttachment(c *fiber.Ctx) error {
	pubIDStr := c.Params("attachment_id")
	pubID, err := uuid.Parse(pubIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid attachment public_id", "Error")
	}

	att, err := cc.attachmentService.GetByPublicID(pubID)
	if err != nil {
		return utils.InternalServerError(c, err.Error(), "Error")
	}
	if att == nil {
		return utils.NotFound(c, "attachment not found")
	}

	// 🔐 Cek permission: hanya owner yang boleh hapus
	u := c.Locals("user")
	token, _ := u.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return utils.Forbidden(c, "invalid user claims")
	}
	userID := int64(userIDFloat)

	if att.UserID != userID {
		return utils.Forbidden(c, "you don't have permission to delete this attachment")
	}

	// Hapus dari database
	if err := cc.attachmentService.DeleteByPublicID(pubID); err != nil {
		return utils.InternalServerError(c, err.Error(), "Error")
	}

	// Hapus file fisik
	_ = os.Remove(att.File)

	return utils.Success(c, "Berhasil menghapus attachment", nil)
}

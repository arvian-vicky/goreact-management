package controllers

import (
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/services"
	"arfian-nurdavicky/project-management/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CommentController struct {
	service services.CommentService
}

func NewCommentController(s services.CommentService) *CommentController {
	return &CommentController{service: s}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment on a card
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Comment data to create"
// @Success 200 {object} ApiResponse{data=models.Comment} "Successfully created comment"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /comments [post]
func (c *CommentController) CreateComment(ctx *fiber.Ctx) error {
	comment := new(models.Comment)
	if err := ctx.BodyParser(comment); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	// Set user ID from JWT
	userID := ctx.Locals("userID").(string)
	comment.UserPubID = uuid.MustParse(userID)

	if err := c.service.Create(comment); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat Komentar", err.Error())
	}

	return utils.Success(ctx, "Komentar berhasil dibuat", comment)
}

// GetCardComments godoc
// @Summary Get comments for a card
// @Description Retrieve all comments for a card
// @Tags comments
// @Produce json
// @Param card_id path string true "Card public ID"
// @Success 200 {object} ApiResponse{data=[]models.Comment} "Successfully retrieved comments"
// @Failure 404 {object} ApiErrorResponse "Card not found"
// @Security BearerAuth
// @Router /cards/{card_id}/comments [get]
func (c *CommentController) GetCardComments(ctx *fiber.Ctx) error {
	cardID := ctx.Params("card_id")
	if _, err := uuid.Parse(cardID); err != nil {
		return utils.BadRequest(ctx, "ID Card tidak valid", err.Error())
	}

	comments, err := c.service.GetByCardID(cardID)
	if err != nil {
		return utils.NotFound(ctx, "Komentar tidak ditemukan")
	}

	return utils.Success(ctx, "Data komentar berhasil diambil", comments)
}

// UpdateComment godoc
// @Summary Update comment
// @Description Update an existing comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment public ID"
// @Param comment body models.Comment true "Updated comment data"
// @Success 200 {object} ApiResponse{data=models.Comment} "Successfully updated comment"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /comments/{id} [put]
func (c *CommentController) UpdateComment(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	comment := new(models.Comment)

	if err := ctx.BodyParser(comment); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	//existingComment, err := c.service.GetByID(publicID)
	existingComment, err := c.service.GetByPublicID(publicID)

	if err != nil {
		return utils.NotFound(ctx, "Komentar tidak ditemukan")
	}

	// Check if user is the owner of the comment
	userID := ctx.Locals("userID").(string)
	if userID != existingComment.UserPubID.String() {
		return utils.Unauthorized(ctx, "Anda tidak memiliki akses untuk mengubah komentar ini")
	}

	comment.InternalID = existingComment.InternalID
	comment.PublicID = existingComment.PublicID

	if err := c.service.Update(comment); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Komentar", err.Error())
	}

	return utils.Success(ctx, "Komentar berhasil diperbarui", comment)
}

// DeleteComment godoc
// @Summary Delete comment
// @Description Delete a comment by its public ID
// @Tags comments
// @Produce json
// @Param id path string true "Comment public ID"
// @Success 200 {object} ApiResponse{data=string} "Comment successfully deleted"
// @Failure 404 {object} ApiErrorResponse "Comment not found"
// @Security BearerAuth
// @Router /comments/{id} [delete]
func (c *CommentController) DeleteComment(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	comment, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Komentar tidak ditemukan")
	}

	// Check if user is the owner of the comment
	userID := ctx.Locals("userID").(string)
	if userID != comment.UserPubID.String() {
		return utils.Unauthorized(ctx, "Anda tidak memiliki akses untuk menghapus komentar ini")
	}

	if err := c.service.Delete(uint(comment.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Komentar", err.Error())
	}

	return utils.Success(ctx, "Komentar berhasil dihapus", publicID)
}

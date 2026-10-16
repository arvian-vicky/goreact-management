package controllers

import (
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/services"
	"arfian-nurdavicky/project-management/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

// CreateList godoc
// @Summary Create a new list
// @Description Create a new list in a board
// @Tags lists
// @Accept json
// @Produce json
// @Param list body models.List true "List data to create"
// @Success 200 {object} ApiResponse{data=models.List} "Successfully created list"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /lists [post]
func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat List", err.Error())
	}

	return utils.Success(ctx, "List berhasil dibuat", list)
}

// GetBoardLists godoc
// @Summary Get lists in a board
// @Description Retrieve all lists in a board
// @Tags lists
// @Produce json
// @Param board_id path string true "Board public ID"
// @Success 200 {object} ApiResponse{data=[]models.List} "Successfully retrieved lists"
// @Failure 404 {object} ApiErrorResponse "Board not found"
// @Security BearerAuth
// @Router /boards/{board_id}/lists [get]
func (c *ListController) GetBoardLists(ctx *fiber.Ctx) error {
	boardID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardID); err != nil {
		return utils.BadRequest(ctx, "ID Board tidak valid", err.Error())
	}

	lists, err := c.service.GetByBoardID(boardID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan")
	}
	fmt.Println(lists)
	return utils.Success(ctx, "Data list berhasil diambil", lists)
}

// UpdateList godoc
// @Summary Update list
// @Description Update an existing list
// @Tags lists
// @Accept json
// @Produce json
// @Param id path string true "List public ID"
// @Param list body models.List true "Updated list data"
// @Success 200 {object} ApiResponse{data=models.List} "Successfully updated list"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /lists/{id} [put]
func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	existingList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan")
	}

	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.service.Update(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Update List", err.Error())
	}

	updatedList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan")
	}

	return utils.Success(ctx, "List berhasil diperbarui", updatedList)
}

// DeleteList godoc
// @Summary Delete list
// @Description Delete a list by its public ID
// @Tags lists
// @Produce json
// @Param id path string true "List public ID"
// @Success 200 {object} ApiResponse{data=string} "List successfully deleted"
// @Failure 404 {object} ApiErrorResponse "List not found"
// @Security BearerAuth
// @Router /lists/{id} [delete]
func (c *ListController) DeleteList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	list, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan")
	}

	if err := c.service.Delete(uint(list.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus List", err.Error())
	}

	return utils.Success(ctx, "List berhasil dihapus", publicID)
}

// UpdateListPositions godoc
// @Summary Update list positions in a board
// @Description Update the order of lists in a board
// @Tags lists
// @Accept json
// @Produce json
// @Param board_id path string true "Board public ID"
// @Param positions body []string true "Array of list public IDs in order"
// @Success 200 {object} ApiResponse "Positions updated successfully"
// @Failure 400 {object} ApiErrorResponse "Bad request"
// @Security BearerAuth
// @Router /boards/{board_id}/positions [put]
func (c *ListController) UpdateListPositions(ctx *fiber.Ctx) error {
	boardID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardID); err != nil {
		return utils.BadRequest(ctx, "ID Board tidak valid", err.Error())
	}

	// Ubah untuk menerima array of UUID bukan array of string
	var positionUUIDs []uuid.UUID
	if err := ctx.BodyParser(&positionUUIDs); err != nil {
		// Jika gagal, coba parse sebagai array of string
		var positionStrings []string
		if err := ctx.BodyParser(&positionStrings); err != nil {
			return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
		}

		// Konversi string ke UUID
		for _, s := range positionStrings {
			u, err := uuid.Parse(s)
			if err != nil {
				return utils.BadRequest(ctx, "Format UUID tidak valid", err.Error())
			}
			positionUUIDs = append(positionUUIDs, u)
		}
	}

	if err := c.service.UpdatePositions(boardID, positionUUIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Posisi", err.Error())
	}

	return utils.Success(ctx, "Posisi list berhasil diperbarui", nil)
}

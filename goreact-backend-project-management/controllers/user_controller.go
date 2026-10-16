package controllers

import (
	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/services"
	"arfian-nurdavicky/project-management/utils"

	"github.com/gofiber/fiber/v2"
)

// ApiResponse godoc
// @Description Standard success response structure
type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ApiErrorResponse godoc
// @Description Standard error response structuret
type ApiErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user with the provided details.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data to register"
// @Success 200 {object} utils.Response{data=models.User} "Successfully registered user with success message"
// @Failure 400 {object} ApiErrorResponse "Bad request, e.g., invalid input or user already exists"
// @Router /v1/auth/register [post]
func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}
	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal ", err.Error())
	}

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	return utils.Success(ctx, "Register successful", userResp)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password to obtain access and refresh tokens.
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body object{email=string,password=string} true "User credentials (email and password)"
// @Success 200 {object} ApiResponse{data=object{access_token=string,refresh_token=string,user=models.UserResponse}} "Login successful, returns access_token, refresh_token, and user data"
// @Failure 400 {object} ApiErrorResponse "Bad request, e.g., invalid input format"
// @Failure 401 {object} ApiErrorResponse "Unauthorized, e.g., invalid credentials"
// @Router /v1/auth/login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Invalid request", err.Error())
	}

	user, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "Gagal Masuk")
	}

	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)

	return utils.Success(ctx, "Login successful", fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          userResp,
	})
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all registered users.
// @Tags users
// @Produce json
// @Success 200 {object} ApiResponse{data=[]models.UserResponse} "Successfully retrieved list of users with success message"
// @Failure 500 {object} ApiErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users [get]
// @Security BearerAuth
func (c *UserController) GetUsers(ctx *fiber.Ctx) error {
	users, err := c.service.GetAll()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Mengambil Data", err.Error()) //err 500
	}

	var userResp []models.UserResponse
	err = copier.Copy(&userResp, &users)
	if err != nil {
		return utils.InternalServerError(ctx, "Internal Server Error", err.Error())
	}

	return utils.Success(ctx, "Data pengguna berhasil diambil", userResp)
}

// GetUsersPaginate godoc
// @Summary Get all users with pagination
// @Description Retrieve a paginated list of users, with optional filtering and sorting.
// @Tags users
// @Produce json
// @Param page query int false "Page number for pagination (default 1)" default(1)
// @Param limit query int false "Number of items per page (default 10)" default(10)
// @Param filter query string false "Filter users by name or email (case-insensitive)"
// @Param sort query string false "Sort order, e.g., '-name for DESC by name', '-email for DESC by email'. Use comma for multiple fields."
// @Success 200 {object} utils.PaginatedResponse "Successfully retrieved paginated user data"
// @Failure 404 {object} utils.PaginatedResponse "No users found matching the filter criteria with paginated metadata"
// @Failure 500 {object} ApiErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users/page [get]
func (c *UserController) GetUsersPaginate(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	users, total, err := c.service.GetAllPaginate(filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Mengambil Data", err.Error()) //err 500
	}

	var userResp []models.UserResponse
	_ = copier.Copy(&userResp, &users)

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}

	//cek apakah ada data sesuai filter
	if total == 0 {
		return utils.PaginatedNotFound(ctx, "Data pengguna tidak ditemukan", userResp, meta)
	}
	return utils.PaginatedSuccess(ctx, "Data pengguna berhasil diambil", userResp, meta)
}

// GetUser godocut
// @Summary Get user by ID
// @Description Get a single user record by their unique ID.
// @Tags users
// @Produce json
// @Param id path int true "User ID to retrieve"
// @Success 200 {object} ApiResponse{data=models.UserResponse} "Successfully retrieved user with success message"
// @Failure 404 {object} ApiErrorResponse "User not found with the given ID"
// @Failure 500 {object} ApiErrorResponse "Internal server error during data copying"
// @Security BearerAuth
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "Data tidak ditemukan")
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp, &user)
	if err != nil {
		return utils.InternalServerError(ctx, "Internal Server Error", err.Error())
	}
	return utils.Success(ctx, "Data pengguna berhasil diambil", userResp)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user record.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data to create"
// @Success 200 {object} ApiResponse{data=models.User} "Successfully created user with success message"
// @Failure 400 {object} ApiErrorResponse "Bad request, e.g., invalid input"
// @Security BearerAuth
// @Router /users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}
	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal ", err.Error())
	}
	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	return utils.Success(ctx, "Berhasil membuat user", userResp)

}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user record by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID to update"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} ApiResponse{data=models.User} "Successfully updated user with success message"
// @Failure 400 {object} ApiErrorResponse "Bad request, e.g., invalid input or user not found"
// @Failure 500 {object} ApiErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var err error = nil
	id := ctx.Params("id")
	user := new(models.User)

	if user.PublicID, err = uuid.Parse(id); err != nil {
		return utils.BadRequest(ctx, "Invalid ID format", err.Error())
	}
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}
	if err := c.service.Update(user); err != nil {
		return utils.InternalServerError(ctx, "Gagal Update Data", err.Error())
	}
	userUpdated, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "Data tidak ditemukan")
	}
	var userResp models.UserResponse
	err = copier.Copy(&userResp, &userUpdated)
	if err != nil {
		return utils.InternalServerError(ctx, "Internal Server Error", err.Error())
	}
	return utils.Success(ctx, "Berhasil update user", userResp)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user record by ID.
// @Tags users
// @Produce json
// @Param id path int true "User ID to delete"
// @Success 200 {object} ApiResponse{data=int} "User successfully deleted, returns a success message and the deleted user ID"
// @Failure 500 {object} ApiErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := c.service.Delete(uint(id)); err != nil {
		return utils.InternalServerError(ctx, "Gagal Menghapus Data", err.Error())
	}
	return utils.Success(ctx, "Berhasil Menghapus User", id)
}

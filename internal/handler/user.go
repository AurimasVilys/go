package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"scootin/internal/models"
	"scootin/internal/repository"
)

type UserHandler struct {
	userRepository repository.UserRepositoryInterface
}

func NewUserHandler(userRepository repository.UserRepositoryInterface) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	user.Identifier = uuid.NewString()

	err := h.userRepository.Insert(c.Request().Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	user, err := h.userRepository.FindByIdentifier(c.Request().Context(), c.Param("identifier"))
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, user)
}

package handler

import (
	"errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	dto "scootin/internal/dto"
	"scootin/internal/models"
	"scootin/internal/repository"
	"strings"
)

type ScooterHandler struct {
	scooterRepository repository.ScooterRepositoryInterface
	userRepository    repository.UserRepositoryInterface
	validator         *validator.Validate
}

func NewScooterHandler(scooterRepository repository.ScooterRepositoryInterface, userRepository repository.UserRepositoryInterface, validator *validator.Validate) *ScooterHandler {
	return &ScooterHandler{
		scooterRepository: scooterRepository,
		userRepository:    userRepository,
		validator:         validator,
	}
}

func (h *ScooterHandler) CreateScooter(c echo.Context) error {
	scooter := models.Scooter{}
	createModel := dto.CreateScooterDTO{}
	if err := c.Bind(&createModel); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("failed to parse request body"))
	}
	scooter.Identifier = uuid.NewString()
	scooter.LastConfirmedLongitude = createModel.Longitude
	scooter.LastConfirmedLatitude = createModel.Latitude

	err := h.validator.Struct(createModel)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	err = h.scooterRepository.Insert(c.Request().Context(), &scooter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, map[string]string{"identifier": scooter.Identifier})
}

func (h *ScooterHandler) PatchScooter(c echo.Context) error {
	updateModel := dto.UpdateScooterDTO{}
	if err := c.Bind(&updateModel); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	scooter, _ := h.scooterRepository.FindOneByIdentifier(c.Request().Context(), c.Param("scooter_identifier"))
	if scooter == nil {
		return c.JSON(http.StatusNotFound, "Scooter not found")
	}

	updateModel.Identifier = scooter.Identifier

	err := h.validator.Struct(updateModel)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	if updateModel.Longitude != nil {
		scooter.LastConfirmedLongitude = *updateModel.Longitude
	}
	if updateModel.Latitude != nil {
		scooter.LastConfirmedLatitude = *updateModel.Latitude
	}
	scooter.OccupiedUserIdentifier = updateModel.OccupiedUserIdentifier

	err = h.scooterRepository.Update(c.Request().Context(), scooter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to update scooter"))
	}

	return c.JSON(http.StatusOK, scooter)
}

func (h *ScooterHandler) OccupyScooter(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.QueryParam("user_identifier")
	scooterID := c.Param("scooter_identifier")

	user, _ := h.userRepository.FindByIdentifier(ctx, userID)
	if user == nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	scooter, _ := h.scooterRepository.FindOneByIdentifier(ctx, scooterID)
	if scooter == nil {
		return c.JSON(http.StatusNotFound, "Scooter not found")
	}
	if scooter.OccupiedUserIdentifier != nil {
		return c.JSON(http.StatusBadRequest, "Scooter already occupied")
	}

	err := h.scooterRepository.UpdateOccupiedByUser(ctx, scooterID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to occupy scooter")
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *ScooterHandler) ReleaseScooter(c echo.Context) error {
	ctx := c.Request().Context()
	scooterID := c.Param("scooter_identifier")

	scooter, err := h.scooterRepository.FindOneByIdentifier(ctx, scooterID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch scooter")
	}
	if scooter == nil {
		return c.JSON(http.StatusNotFound, "Scooter not found")
	}
	if scooter.OccupiedUserIdentifier == nil {
		return c.JSON(http.StatusBadRequest, "Scooter already released")
	}

	err = h.scooterRepository.ReleaseOccupied(ctx, scooterID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to release scooter")
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *ScooterHandler) QueryScooter(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse query parameters for coordinates
	startPoint := strings.Split(c.QueryParam("start_point"), ",")
	endPoint := strings.Split(c.QueryParam("end_point"), ",")
	includeOccupied := c.QueryParam("include_occupied") == "1"

	if len(startPoint) != 2 || len(endPoint) != 2 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid coordinates format")
	}

	startLat, startLng := startPoint[0], startPoint[1]
	endLat, endLng := endPoint[0], endPoint[1]

	// Query scooters by coordinates
	scooters, err := h.scooterRepository.FindByCoordinatesAndStatus(ctx, startLat, endLat, startLng, endLng, includeOccupied)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to query scooters")
	}

	// Return response
	if len(scooters) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"Message": "No scooters found in the specified area",
		})
	}
	return c.JSON(http.StatusOK, scooters)
}

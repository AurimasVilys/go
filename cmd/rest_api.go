package cmd

import (
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"scootin/internal/contraints"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/stephenafamo/bob"
	"scootin/internal/handler"
	"scootin/internal/repository"

	"github.com/go-playground/validator"
)

func InitializeRestfulApi(wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.KeyAuth(
		func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("REST_API_KEY"), nil
		},
	))

	db, _ := openDB()

	// Scooter
	scooterHandler := InitializeScooterHandler(db)
	e.POST("/scooters", scooterHandler.CreateScooter)
	e.PATCH("/scooters/:scooter_identifier", scooterHandler.PatchScooter)
	e.GET("/scooters", scooterHandler.QueryScooter)

	// User
	userHandler := InitializeUserHandler(db)
	e.POST("/user", userHandler.CreateUser)
	e.GET("/user/:identifier", userHandler.GetUser)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("REST_API_PORT"))))
}

func InitializeScooterHandler(db *bob.DB) *handler.ScooterHandler {
	scooterRepository := repository.NewScooterRepository(db)
	userRepository := repository.NewUserRepository(db)

	return handler.NewScooterHandler(scooterRepository, userRepository, InitializeValidator(db))
}

func InitializeUserHandler(db *bob.DB) *handler.UserHandler {
	userRepository := repository.NewUserRepository(db)

	return handler.NewUserHandler(userRepository)
}

func InitializeValidator(db *bob.DB) *validator.Validate {
	validate := validator.New()

	scooterRepository := repository.NewScooterRepository(db)
	checkOccupiedValidator := contraints.NewCheckOccupiedValidator(scooterRepository)
	err := validate.RegisterValidation("scooter_dto:check_occupied", checkOccupiedValidator.Validate, true)
	if err != nil {
		return nil
	}

	return validate
}

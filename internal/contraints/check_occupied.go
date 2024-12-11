package contraints

import (
	"context"
	"github.com/go-playground/validator"
	"scootin/internal/dto"
	"scootin/internal/repository"
)

type CheckOccupiedValidator struct {
	scooterRepository repository.ScooterRepositoryInterface
}

func NewCheckOccupiedValidator(scooterRepository repository.ScooterRepositoryInterface) *CheckOccupiedValidator {
	return &CheckOccupiedValidator{scooterRepository: scooterRepository}
}

func (v *CheckOccupiedValidator) Validate(fl validator.FieldLevel) bool {
	identifier := fl.Parent().Interface().(dto.UpdateScooterDTO)
	scooter, _ := v.scooterRepository.FindOneByIdentifier(context.Background(), identifier.Identifier)
	if scooter == nil {
		return true
	}

	if scooter.OccupiedUserIdentifier != nil && !fl.Field().IsNil() {
		return false
	}

	if scooter.OccupiedUserIdentifier == nil && fl.Field().IsValid() && fl.Field().String() == "" {
		return false

	}

	return true
}

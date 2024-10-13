package usecase

import (
	"errors"

	"SavingBooks/internal/domain"
)

func findSavingType(inputTerm int, regulation * domain.SavingRegulation) (*domain.SavingType, error) {
	for _, value := range regulation.SavingTypes {
		if value.Term == inputTerm {
			return &value, nil
		}
	}
	return nil, errors.New("Regulation term not found")
}

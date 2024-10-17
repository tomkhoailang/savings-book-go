package presenter

import (
	"errors"
	"sort"
)

type SavingRegulationInput struct {
	MinWithdrawValue float64      `json:"minWithdrawValue" validate:"min=10"`
	SavingTypes      []SavingType `json:"savingTypes" validate:"required"`
	MinWithdrawDay   int          `json:"minWithdrawDay" validate:"required"`
	IsActive         bool         `json:"isActive" validate:"required"`
}

type SavingType struct {
	Name         string  `json:"name"`
	Term         int     `json:"term" validate:"min=1"`
	InterestRate float64 `json:"interestRate"`
}

func(input *SavingRegulationInput) Validate() error {
	if len(input.SavingTypes) == 1{
		return errors.New(SavingTypeMinimumLength)
	}
	termSet := make(map[int]int, len(input.SavingTypes))
	sort.Slice(input.SavingTypes, func(i, j int) bool {
		return input.SavingTypes[i].Term < input.SavingTypes[j].Term
	})
	interestRate := -1.0
	for _, savingType := range input.SavingTypes {
		termSet[savingType.Term] += 1
		if interestRate >= savingType.InterestRate {
			return errors.New(SavingTypeRateMismatch)
		}
		interestRate = savingType.InterestRate
	}
	if _,ok := termSet[0]; !ok {
		return errors.New(SavingTypeMissingNoTerm)
	}
	for _, value := range termSet {
		if value >= 2{
			return errors.New(SavingTypeTermDuplicate)
		}
	}
	return nil
}
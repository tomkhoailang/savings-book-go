package saving_regulation

import "SavingBooks/internal/domain"

type SavingRegulationRepository interface {
	domain.GenericRepository[domain.SavingRegulation]
}

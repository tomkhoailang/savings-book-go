package saving_book

import "SavingBooks/internal/domain"

type SavingBookRepository interface {
	domain.GenericRepository[domain.SavingBook]
}

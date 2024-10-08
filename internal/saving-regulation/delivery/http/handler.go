package http

import (
	"SavingBooks/internal/domain"
	saving_regulation "SavingBooks/internal/saving-regulation"
	"SavingBooks/internal/saving-regulation/presenter"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type savingRegulationHandler struct {
	savingRegulationUC saving_regulation.SavingRegulationUseCase
}



func (s *savingRegulationHandler) CreateRegulation() gin.HandlerFunc {
	return utils.HandleCreateRequest[presenter.SavingRegulationInput, presenter.SavingRegulationOutput, domain.SavingRegulation](s.savingRegulationUC.CreateRegulation)
}

func (s *savingRegulationHandler) UpdateRegulation() gin.HandlerFunc {
	return utils.HandleUpdateRequest[presenter.SavingRegulationInput, presenter.SavingRegulationOutput, domain.SavingRegulation](s.savingRegulationUC.UpdateRegulation)
}

func (s *savingRegulationHandler) DeleteManyRegulations() gin.HandlerFunc {
	return utils.HandleDeleteManyRequest[domain.SavingRegulation](s.savingRegulationUC.DeleteManyRegulations)
}

func (s *savingRegulationHandler) GetListRegulations() gin.HandlerFunc {
	return utils.HandleGetListRequest[domain.SavingRegulation](s.savingRegulationUC.GetListRegulation)
}

func NewSavingRegulationHandler(savingRegulationUC saving_regulation.SavingRegulationUseCase) saving_regulation.Handler {
	return &savingRegulationHandler{savingRegulationUC: savingRegulationUC}
}


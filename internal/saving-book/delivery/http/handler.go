﻿package http

import (
	"net/http"

	"SavingBooks/internal/domain"
	saving_book "SavingBooks/internal/saving-book"
	"SavingBooks/internal/saving-book/presenter"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type savingBookHandler struct {
	savingBookUC saving_book.UseCase
}

func (s *savingBookHandler) DepositOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input presenter.DepositInput

		savingBookId := c.Param("id")
		if savingBookId == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"err":"id can not be empty"})
			return
		}

		err := utils.ReadRequest(c, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ticket, err := s.savingBookUC.DepositOnline(c, &input, savingBookId, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ticket)
		return

	}
}

func (s *savingBookHandler) WithDrawOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input presenter.WithDrawInput

		savingBookId := c.Param("id")
		if savingBookId == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"err":"id can not be empty"})
			return
		}

		err := utils.ReadRequest(c, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = s.savingBookUC.WithdrawOnline(c, &input, savingBookId, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"error": "Your request is being processed"})
		return

	}
}

func (s *savingBookHandler) GetListSavingBook() gin.HandlerFunc {
	return utils.HandleGetListRequestAuth[domain.SavingBook](s.savingBookUC.GetListSavingBook)
}

func (s *savingBookHandler) CreateSavingBookOnline() gin.HandlerFunc {
	return utils.HandleCreateRequest[presenter.SavingBookGuestInput, presenter.SavingBookOutput, domain.SavingBook](s.savingBookUC.CreateSavingBookOnline)
}

func (s *savingBookHandler) ConfirmPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input presenter.ConfirmPaymentInput

		errors := utils.ReadRequest(c, &input)
		if errors != nil {
			c.JSON(http.StatusBadRequest, errors)
			return
		}

		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = s.savingBookUC.ConfirmPaymentOnline(c.Request.Context(), input.PaymentId, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "ok")
		return
	}
}


func NewSavingBookHandler(savingBookUC saving_book.UseCase) saving_book.Handler {
	return &savingBookHandler{savingBookUC: savingBookUC}
}

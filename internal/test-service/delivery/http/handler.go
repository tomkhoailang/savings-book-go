package http

import (
	test_service "SavingBooks/internal/test-service"
	"github.com/gin-gonic/gin"
)

type testServiceHandler struct {
	test_service.UseCase
}

func (t *testServiceHandler) TestKafkaProducer() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := t.UseCase.TestProducer()
		if err != nil {
			c.JSON(400, gin.H{"error:": err.Error()})
			return
		}
		c.JSON(200, "ok day")
		return
	}
}

func NewTestServiceHandler(useCase test_service.UseCase) test_service.Handler {
	return &testServiceHandler{UseCase: useCase}
}



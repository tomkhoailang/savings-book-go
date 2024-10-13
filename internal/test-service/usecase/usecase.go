package usecase

import (
	"SavingBooks/internal/services/kafka"
	test_service "SavingBooks/internal/test-service"
)

type testServiceUseCase struct {
	producer *kafka.KafkaProducer
}

func (t *testServiceUseCase) TestProducer() error{
	err := t.producer.SendMessage("aa",[]byte("that su la test"))
	if err != nil {
		return err
	}
	return nil
}

func NewTestServiceUseCase(producer *kafka.KafkaProducer) test_service.UseCase {
	return &testServiceUseCase{producer: producer}
}
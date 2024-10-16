package usecase

import (
	"SavingBooks/internal/services/kafka"
	"SavingBooks/internal/services/websocket"
	test_service "SavingBooks/internal/test-service"
)

type testServiceUseCase struct {
	producer *kafka.KafkaProducer
	socket *websocket.Hub
}

func (t *testServiceUseCase) TestProducer() error{
	err := t.producer.SendMessage("aa",[]byte("that su la test"))
	if err != nil {
		return err
	}
	t.socket.SendAll("hehe", "heheheheheheh")
	return nil
}

func NewTestServiceUseCase(producer *kafka.KafkaProducer, socket *websocket.Hub) test_service.UseCase {
	return &testServiceUseCase{producer: producer, socket: socket}
}
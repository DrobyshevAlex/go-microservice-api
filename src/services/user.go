package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"main/src/ampq/requests"
	"main/src/ampq/responses"

	"time"

	"github.com/google/uuid"
	go_amqp_lib "github.com/lan143/go-amqp-lib"
)

type UserService struct {
	amqp *go_amqp_lib.Client
}

func (s *UserService) GetUserById(ctx context.Context, userId int64) (*responses.UserResponse, error) {
	timeout := 20 * time.Second
	timeoutCxt, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	userProfileRequest := requests.UserGetByIdRequest{Id: userId}
	request := go_amqp_lib.Request[requests.UserGetByIdRequest]{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(timeout),
		Payload:   userProfileRequest,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	consume, err := s.amqp.Consume("", true, false)
	if err != nil {
		return nil, err
	}
	defer func(amqp *go_amqp_lib.Client, consume *go_amqp_lib.Consume) {
		err := amqp.CloseConsume(consume)
		if err != nil {
			log.Printf("user-service.get-user: %s", err.Error())
		}
	}(s.amqp, consume)

	err = s.amqp.Publish(consume.Channel, "user.get", requestBytes, &consume.QueueName)
	if err != nil {
		return nil, err
	}

	response := go_amqp_lib.Response[responses.UserResponse]{}

	select {
	case <-timeoutCxt.Done():
		err = errors.New("timed out waiting for a auth microservice")
		return nil, err
	case delivery := <-consume.Delivery:
		_ = delivery.Ack(false)
		err := json.Unmarshal(delivery.Body, &response)

		if err != nil {
			return nil, err
		}

		break
	}

	if response.Success {
		return &response.Payload, nil
	} else {
		return nil, errors.New(response.Message)
	}
}

func NewUserService(amqp *go_amqp_lib.Client) *UserService {
	return &UserService{
		amqp: amqp,
	}
}

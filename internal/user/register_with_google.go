package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/vardius/go-api-boilerplate/pkg/domain"
)

// RegisterWithGoogle command bus contract
const RegisterWithGoogle = "register-user-with-google"

type registerWithGoogle struct {
	Email     string `json:"email"`
	AuthToken string `json:"authToken"`
}

func (c *registerWithGoogle) fromJSON(payload json.RawMessage) error {
	return json.Unmarshal(payload, c)
}

// OnRegisterWithGoogle creates command handler
func OnRegisterWithGoogle(es domain.EventStore, eb domain.EventBus) domain.CommandHandler {
	repository := newRepository(fmt.Sprintf("%T", User{}), es, eb)

	return func(ctx context.Context, payload json.RawMessage, out chan<- error) {
		c := &registerWithGoogle{}
		err := c.fromJSON(payload)
		if err != nil {
			out <- err
			return
		}

		//todo: validate if email is taken or if user already connected with google

		id, err := uuid.NewRandom()
		if err != nil {
			out <- err
			return
		}

		user := New()
		err = user.RegisterWithGoogle(id, c.Email, c.AuthToken)
		if err != nil {
			out <- err
			return
		}

		out <- nil

		repository.Save(domain.ContextWithFlag(ctx, domain.LIVE), user)
	}
}
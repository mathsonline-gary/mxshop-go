package v1

import (
	"context"
	"mxshop-go/app/user/svc/repository/v1/mock"
	"testing"
)

func TestUserIndex(t *testing.T) {
	us := NewUserService(mock.NewUserRepository())

	us.Index(context.Background(), ListMeta{})
}

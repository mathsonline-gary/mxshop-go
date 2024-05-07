package v1

import (
	"context"
	"testing"

	"github.com/zycgary/mxshop-go/app/user/svc/repository/v1/mock"
)

func TestUserIndex(t *testing.T) {
	us := NewUserService(mock.NewUserRepository())

	_, _ = us.Index(context.Background(), ListMeta{})
}

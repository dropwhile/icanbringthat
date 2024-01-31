package rpc

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/dropwhile/icbt/internal/app/service/mockservice"
)

func NewTestServer(t *testing.T) (*Server, *mockservice.MockServicer) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mock := mockservice.NewMockServicer(ctrl)
	return &Server{svc: mock}, mock
}

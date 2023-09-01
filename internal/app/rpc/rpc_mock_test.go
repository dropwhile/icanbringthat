// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/dropwhile/icanbringthat/internal/app/service/mockservice"
)

func NewTestServer(t *testing.T) (*Server, *mockservice.MockServicer) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mock := mockservice.NewMockServicer(ctrl)
	return &Server{svc: mock}, mock
}

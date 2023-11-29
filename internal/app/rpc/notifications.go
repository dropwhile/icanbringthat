package rpc

import (
	"context"
	"errors"

	pb "github.com/dropwhile/icbt/rpc"
)

func (s *Server) ListNotifications(ctx context.Context,
	request *pb.ListNotificationsRequest,
) (*pb.ListNotificationsResponse, error) {
	return nil, errors.New("nothing yet")
}

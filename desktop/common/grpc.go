package common

import (
	"context"
	"sid-desktop/proto"
)

type DeskTop struct{}

func (s *DeskTop) UploadWebInfo(ctx context.Context, req *proto.UploadWebInfoRequest) (*proto.UploadWebInfoResponse, error) {
	return &proto.UploadWebInfoResponse{}, nil
}

package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"boilerplate/internal/pkg/errors"
)

func Error(err error) error {
	if errors.IsErrBadRequest(err) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.IsErrNotFound(err) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.IsErrForbidden(err) {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	if errors.IsErrUnauthorized(err) {
		return status.Error(codes.Unauthenticated, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}

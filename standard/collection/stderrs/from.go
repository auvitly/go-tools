package stderrs

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StandardError - interface to get error from source.
type StandardError interface {
	StandardFrom(err error) *Error
}

// FromFunc - function to get error from source.
type FromFunc func(err error) *Error

var (
	fromHandlers = []FromFunc{
		fromModel,
		fromGRPC,
	}
	fromUserHandlers []FromFunc
)

var grpcMapper = map[codes.Code]*Error{
	codes.Canceled:           Canceled,
	codes.Unknown:            Unknown,
	codes.InvalidArgument:    InvalidArgument,
	codes.DeadlineExceeded:   DeadlineExceeded,
	codes.NotFound:           NotFound,
	codes.AlreadyExists:      AlreadyExists,
	codes.PermissionDenied:   PermissionDenied,
	codes.ResourceExhausted:  ResourceExhausted,
	codes.FailedPrecondition: FailedPrecondition,
	codes.Aborted:            Aborted,
	codes.OutOfRange:         OutOfRange,
	codes.Unimplemented:      Unimplemented,
	codes.Internal:           Internal,
	codes.Unavailable:        Unavailable,
	codes.DataLoss:           DataLoss,
	codes.Unauthenticated:    Unauthenticated,
}

// RegistryFrom - allows you to register a custom error parser for recovery from the standard interface.
func RegistryFrom(fns ...FromFunc) {
	for _, fn := range fns {
		if fn == nil {
			continue
		}

		fromUserHandlers = append(fromUserHandlers, fn)
	}
}

// From - error parser from the base interface.
// If the passed error is nil, a processing error will be returned.
// If the error could not be determined, an enriched Undefined error will be returned.
func From(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}

	if stderr, ok := err.(StandardError); ok {
		var std = stderr.StandardFrom(err)

		if std == nil {
			return nil, false
		}

		return std, true
	}

	for _, handler := range fromUserHandlers {
		if std := handler(err); std != nil {
			return std, true
		}
	}

	for _, handler := range fromHandlers {
		if std := handler(err); std != nil {
			return std, true
		}
	}

	return nil, false
}

func fromModel(err error) (std *Error) {
	var result *Error

	if errors.As(err, &result) {
		return result
	}

	return nil
}

func fromGRPC(err error) *Error {
	if parsed, ok := status.FromError(err); ok {
		std, found := grpcMapper[parsed.Code()]
		if !found {
			return nil
		}

		return std.
			SetMessage(parsed.Message()).
			EmbedErrors(errors.New(parsed.Message()))
	}

	return nil
}

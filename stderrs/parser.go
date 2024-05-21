package stderrs

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ParseFunc - function to get error from source.
type ParseFunc func(err error) *Error

var (
	g_parsers = []ParseFunc{
		parseModel,
		parseGRPC,
	}
	l_parsers []ParseFunc
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

// RegistryParsers - allows you to register a custom error parser for recovery from the standard interface.
func RegistryParsers(fns ...ParseFunc) {
	for _, fn := range fns {
		if fn == nil {
			continue
		}

		l_parsers = append(l_parsers, fn)
	}
}

// Parse - error parser from the base interface.
// If the passed error is nil, a processing error will be returned.
// If the error could not be determined, an enriched Undefined error will be returned.
func Parse(err error) *Error {
	if err == nil {
		return Undefined.
			SetMessage("incorrect error handling")
	}

	for _, handler := range l_parsers {
		if std := handler(err); std != nil {
			return std
		}
	}

	for _, handler := range g_parsers {
		if std := handler(err); std != nil {
			return std
		}
	}

	return Undefined.
		SetMessage("internal server error").
		EmbedErrors(err)
}

func parseModel(err error) (std *Error) {
	switch v := err.(type) {
	case Error:
		return &v
	case *Error:
		return v
	default:
		var result *Error

		if errors.As(err, result) {
			return result
		}
	}

	return nil
}

func parseGRPC(err error) *Error {
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

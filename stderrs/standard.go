package stderrs

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

const (
	msgInternal = "internal server error"
)

var (
	// Undefined - unspecified error.
	Undefined = New("").
			SetMessage(msgInternal).
			SetGRPCCode(codes.Internal).
			SetHTTPCode(http.StatusInternalServerError)

	// Panic - unhandled exception.
	Panic = New("panic").
		SetMessage(msgInternal).
		SetGRPCCode(codes.Internal).
		SetHTTPCode(http.StatusInternalServerError)

	// Canceled - indicates the operation was canceled (typically by the caller).
	Canceled = New("canceled").
			SetMessage("canceled").
			SetGRPCCode(codes.Canceled).
			SetHTTPCode(499)

	// Unknown error. An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown = New("unknown").
		SetMessage(msgInternal).
		SetGRPCCode(codes.Unknown).
		SetHTTPCode(http.StatusInternalServerError)

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument = New("invalid_argument").
			SetMessage("bad request").
			SetGRPCCode(codes.InvalidArgument).
			SetHTTPCode(http.StatusBadRequest)

	// DeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	DeadlineExceeded = New("deadline_exceeded").
				SetMessage("deadline exceeded").
				SetGRPCCode(codes.DeadlineExceeded).
				SetHTTPCode(http.StatusBadGateway)

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	NotFound = New("not_found").
			SetMessage("not found").
			SetGRPCCode(codes.NotFound).
			SetHTTPCode(http.StatusNotFound)

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists = New("already_exists").
			SetMessage("already exists").
			SetGRPCCode(codes.AlreadyExists).
			SetHTTPCode(http.StatusConflict)

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	PermissionDenied = New("permission_denied").
				SetMessage("permission denied").
				SetGRPCCode(codes.PermissionDenied).
				SetHTTPCode(http.StatusForbidden)

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted = New("resource_exhausted").
				SetMessage("resource has been exhausted").
				SetGRPCCode(codes.ResourceExhausted).
				SetHTTPCode(http.StatusTooManyRequests)

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed. E.g., if an "rmdir"
	//      fails because the directory is non-empty, FailedPrecondition
	//      should be returned since the client should not retry unless
	//      they have first fixed up the directory by deleting files from it.
	//  (d) Use FailedPrecondition if the client performs conditional
	//      REST Get/Update/Delete on a resource and the resource on the
	//      server does not match the condition. E.g., conflicting
	//      read-modify-write on the same resource.
	FailedPrecondition = New("failed_precondition").
				SetMessage("system is not in a state required for the operation's execution").
				SetGRPCCode(codes.FailedPrecondition).
				SetHTTPCode(http.StatusBadRequest)

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Aborted = New("aborted").
		SetMessage("aborted").
		SetGRPCCode(codes.Aborted).
		SetHTTPCode(http.StatusConflict)

	// OutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange. We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	OutOfRange = New("out_of_range").
			SetMessage("attempted past the valid range").
			SetGRPCCode(codes.OutOfRange).
			SetHTTPCode(http.StatusBadRequest)

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented = New("unimplemented").
			SetMessage("not implemented or not supported/enabled").
			SetGRPCCode(codes.Unimplemented).
			SetHTTPCode(http.StatusNotImplemented)

	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal = New("internal").
			SetMessage(msgInternal).
			SetGRPCCode(codes.Internal).
			SetHTTPCode(http.StatusInternalServerError)

	// Unavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Unavailable = New("unavailable").
			SetMessage("service unavailable").
			SetGRPCCode(codes.Unavailable).
			SetHTTPCode(http.StatusServiceUnavailable)

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss = New("data_loss").
			SetMessage("unrecoverable data loss or corruption").
			SetGRPCCode(codes.DataLoss).
			SetHTTPCode(http.StatusInternalServerError)

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated = New("unauthenticated").
			SetMessage("request does not have valid authentication credentials").
			SetGRPCCode(codes.Unauthenticated).
			SetHTTPCode(http.StatusUnauthorized)
)

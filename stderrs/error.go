package stderrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auvitly/go-tools/stderrs/internal/unwrap"
	"google.golang.org/grpc/codes"
	"net/http"
	"strings"
)

// Error - unified model.
type Error struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Embed   error          `json:"embed"`
	Wraps   []string       `json:"wraps"`
	Fields  map[string]any `json:"fields"`
	Codes   struct {
		GRPC codes.Code `json:"grpc"`
		HTTP int        `json:"http"`
	} `json:"codes"`
}

// New - create error.
func New(code string) *Error {
	return &Error{
		Code: code,
		Codes: struct {
			GRPC codes.Code `json:"grpc"`
			HTTP int        `json:"http"`
		}{
			GRPC: codes.Unknown,
			HTTP: http.StatusInternalServerError,
		},
	}
}

// SetCode - set general code. The code influences the error definition.
// Errors are considered equal if their codes match.
func (e *Error) SetCode(code *Error) *Error {
	if e == nil {
		return e
	}

	var result = *e

	result.Code = code.Code
	result.Codes = code.Codes

	return &result
}

// SetMessage - set general message.
func (e *Error) SetMessage(format string, args ...any) *Error {
	if e == nil {
		return e
	}

	var result = *e

	result.Message = fmt.Sprintf(format, args...)

	return &result
}

// SetHTTPCode - set HTTP status code.
func (e *Error) SetHTTPCode(code int) *Error {
	if e == nil {
		return e
	}

	var result = *e

	result.Codes.HTTP = code

	return &result
}

// SetGRPCCode - set GRPC status code.
func (e *Error) SetGRPCCode(code codes.Code) *Error {
	if e == nil {
		return e
	}

	var result = *e

	result.Codes.GRPC = code

	return &result
}

// EmbedErrors - add a nested errors.
func (e *Error) EmbedErrors(errs ...error) *Error {
	if e == nil {
		return e
	}

	var result = *e

	var list = make([]error, 0, len(errs))

	for _, err := range errs {
		if err != nil {
			list = append(list, err)
		}
	}

	switch v := e.Embed.(type) {
	case interface{ Unwrap() error }:
		var join = append([]error{v.Unwrap()}, list...)

		result.Embed = errors.Join(join...)
	case interface{ Unwrap() []error }:
		var join = append(v.Unwrap(), list...)

		result.Embed = errors.Join(join...)
	case nil:
		result.Embed = errors.Join(list...)
	default:
		var join = append([]error{e.Embed}, list...)

		result.Embed = errors.Join(join...)
	}

	return &result
}

// Wrap - add a nested.
func (e *Error) Wrap(msg ...string) *Error {
	if e == nil {
		return e
	}

	var (
		result = *e
		wraps  = make([]string, len(e.Wraps))
	)

	copy(wraps, e.Wraps)

	result.Wraps = append(wraps, msg...)

	return &result
}

// WithField - add fields to description.
func (e *Error) WithField(key string, value any) *Error {
	if e == nil {
		return e
	}

	var (
		result = *e
		data   = make(map[string]any)
	)

	for k, v := range e.Fields {
		data[k] = v
	}

	data[key] = value

	result.Fields = data

	return &result
}

// WithFieldIf - add fields to description with condition.
func (e *Error) WithFieldIf(condition bool, key string, value any) *Error {
	if e == nil {
		return e
	}

	if !condition {
		return e
	}

	var (
		result = *e
		data   = make(map[string]any)
	)

	for k, v := range e.Fields {
		data[k] = v
	}

	data[key] = value

	result.Fields = data

	return &result
}

// WithFields - add fields to description.
func (e *Error) WithFields(fields map[string]any) *Error {
	if e == nil {
		return e
	}

	var (
		result = *e
		data   = make(map[string]any)
	)

	for k, v := range e.Fields {
		data[k] = v
	}

	for k, v := range fields {
		data[k] = v
	}

	result.Fields = data

	return &result
}

// Error - implementation of the standard interface.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	var parts []string

	if len(e.Code) != 0 {
		parts = append(parts, fmt.Sprintf(`"code": "%s"`, e.Code))
	} else {
		parts = append(parts, fmt.Sprintf(`"code": "undefined"`))
	}

	if len(e.Message) != 0 {
		parts = append(parts, fmt.Sprintf(`"message": "%s"`, e.Message))
	}

	if len(e.Fields) != 0 {
		raw, err := json.Marshal(e.Fields)
		if err == nil {
			parts = append(parts, fmt.Sprintf(`"fields": %s`, raw))
		}
	}

	if e.Embed != nil {
		msg := strings.Replace(e.Embed.Error(), "\n", ",", -1)

		parts = append(parts, fmt.Sprintf(`"embed": [%s]`, msg))
	}

	var message = strings.Join(parts, ", ")

	for i := 0; i < len(e.Wraps); i++ {
		message = fmt.Sprintf("%s > %s", e.Wraps[i], message)
	}

	return fmt.Sprintf("{%s}", message)
}

// Unwrap - implementation of the standard interface.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Embed
}

// Is - implementation of the standard interface.
func (e *Error) Is(err error) bool {
	if e == nil || err == nil {
		return false
	}

	if std, ok := err.(*Error); ok {
		switch {
		case e.Code == std.Code && e.Embed != nil && std.Embed == nil:
			return true
		case e.Code == std.Code && e.Embed == nil && std.Embed == nil:
			return true
		case e.Code == std.Code && e.Embed == nil && std.Embed != nil:
			return false
		case e.Code == std.Code && e.Embed != nil && std.Embed != nil:
			for _, item := range unwrap.Do(std.Embed) {
				if !errors.Is(e.Embed, item) {
					return false
				}
			}

			return true
		case e.Code != std.Code && e.Embed != nil:
			return errors.Is(e.Embed, err)
		}
	}

	for _, item := range unwrap.Do(err) {
		if errors.Is(e, item) {
			return true
		}
	}

	if errors.Is(e.Embed, err) {
		return true
	}

	return false
}

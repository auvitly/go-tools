package stderrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auvitly/go-tools/stderrs/internal/utils"
	"google.golang.org/grpc/codes"
	"net/http"
	"regexp"
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
func (e Error) SetCode(code string) *Error {
	e.Code = code

	return &e
}

// SetMessage - set general message.
func (e Error) SetMessage(format string, args ...any) *Error {
	e.Message = fmt.Sprintf(format, args...)

	return &e
}

// SetHTTPCode - set HTTP status code.
func (e Error) SetHTTPCode(code int) *Error {
	e.Codes.HTTP = code

	return &e
}

// SetGRPCCode - set GRPC status code.
func (e Error) SetGRPCCode(code codes.Code) *Error {
	e.Codes.GRPC = code

	return &e
}

// EmbedErrors - add a nested errors.
func (e Error) EmbedErrors(errs ...error) *Error {
	switch v := e.Embed.(type) {
	case interface{ Unwrap() error }:
		var join = append([]error{v.Unwrap()}, errs...)

		e.Embed = errors.Join(join...)
	case interface{ Unwrap() []error }:
		var join = append(v.Unwrap(), errs...)

		e.Embed = errors.Join(join...)
	case nil:
		e.Embed = errors.Join(errs...)
	default:
		var join = append([]error{e.Embed}, errs...)

		e.Embed = errors.Join(join...)
	}

	return &e
}

// Wrap - add a nested.
func (e Error) Wrap(msg ...string) *Error {
	var wraps = make([]string, len(e.Wraps))

	copy(wraps, e.Wraps)

	e.Wraps = append(wraps, msg...)

	return &e
}

// WithField - add fields to description.
func (e Error) WithField(key string, value any) *Error {
	var data = make(map[string]any)

	for k, v := range e.Fields {
		data[k] = v
	}

	data[key] = value

	e.Fields = data

	return &e
}

// WithFieldIf - add fields to description with condition.
func (e Error) WithFieldIf(condition bool, key string, value any) *Error {
	if !condition {
		return &e
	}

	var data = make(map[string]any)

	for k, v := range e.Fields {
		data[k] = v
	}

	data[key] = value

	e.Fields = data

	return &e
}

// WithFields - add fields to description.
func (e Error) WithFields(fields map[string]any) *Error {
	var data = make(map[string]any)

	for k, v := range e.Fields {
		data[k] = v
	}

	for k, v := range fields {
		data[k] = v
	}

	e.Fields = data

	return &e
}

// Error - implementation of the standard interface.
func (e Error) Error() string {
	var parts []string

	if len(e.Code) != 0 {
		parts = append(parts, fmt.Sprintf(`"code": "%s"`, e.Code))
	} else {
		parts = append(parts, fmt.Sprintf(`"code": "undefined"`))
	}

	if e.Embed != nil {
		parts = append(parts, fmt.Sprintf(`"embed": "%s"`, e.Embed.Error()))
	}

	if len(e.Message) != 0 {
		parts = append(parts, fmt.Sprintf(`"message": "%s"`, e.Message))
	}

	if len(e.Fields) != 0 {
		raw, err := json.Marshal(e.Fields)
		if err == nil {
			parts = append(parts, fmt.Sprintf(`"fields": "%s"`, raw))
		}
	}

	var message = strings.Join(parts, "; ")

	for i := 0; i < len(e.Wraps); i++ {
		message = fmt.Sprintf("%s > %s", e.Wraps[i], message)
	}

	return message
}

// Unwrap - implementation of the standard interface.
func (e Error) Unwrap() error {
	return e.Embed
}

// Is - implementation of the standard interface.
func (e Error) Is(err error) bool {
	if err == nil {
		return false
	}

	if std, ok := Parse(err); ok {
		switch {
		case e.Code == std.Code && e.Embed != nil && std.Embed == nil:
			return true
		case e.Code == std.Code && e.Embed == nil && std.Embed == nil:
			return true
		case e.Code == std.Code && e.Embed == nil && std.Embed != nil:
			return false
		case e.Code == std.Code && e.Embed != nil && std.Embed != nil:
			for _, item := range utils.Unwrap(std.Embed) {
				if !errors.Is(e.Embed, item) {
					return false
				}
			}

			return true
		case e.Code != std.Code && e.Embed != nil:
			return errors.Is(e.Embed, err)
		}
	}

	if errors.Is(e.Embed, err) {
		return true
	}

	for _, item := range utils.Unwrap(err) {
		if errors.Is(e, item) {
			return true
		}
	}

	return false
}

// Contains - check contains error. Contains method is an extended version Is,
// where the matching of the error text is alo checked. It is used as a replacement for Is
// when analyzing an error recovered from JSON.
func (e Error) Contains(err error) bool {
	switch v := err.(type) {
	case Error:
		return e.Contains(&v)
	case *Error:
		switch {
		case e.Code == v.Code && e.Embed != nil && v.Embed != nil:
			switch emb := v.Embed.(type) {
			case interface{ Unwrap() error }:
				if !strings.Contains(emb.Unwrap().Error(), e.Embed.Error()) {
					return false
				}

				return true
			case interface{ Unwrap() []error }:
				for _, sub := range emb.Unwrap() {
					if !e.Contains(sub) {
						return false
					}
				}

				return true
			default:
				return strings.Contains(e.Embed.Error(), v.Embed.Error())
			}
		case e.Code == v.Code && v.Embed == nil:
			return true
		}
	}

	if e.Embed == nil {
		return false
	}

	return strings.Contains(e.Error(), err.Error())
}

// Match - checking for the presence of a regular expression.
func (e Error) Match(expr string) bool {
	re, err := regexp.Compile(expr)
	if err != nil {
		return false
	}

	return re.MatchString(e.Embed.Error())
}

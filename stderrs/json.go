package stderrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auvitly/go-tools/stderrs/internal/models"
	"strings"
)

// MarshalJSON - implementation of the standard interface.
func (e *Error) MarshalJSON() (_ []byte, err error) {
	var raw models.Error

	raw.Code = e.Code

	if len(raw.Code) == 0 {
		raw.Code = "undefined"
	}

	raw.Message = e.Message
	raw.Fields = e.Fields
	raw.Codes = e.Codes
	raw.Embed, err = marshalJSONEmbedError(e.Embed)
	if err != nil {
		return nil, err
	}

	return json.Marshal(raw)
}

func marshalJSONEmbedError(in error) (*models.EmbedError, error) {
	if in == nil {
		return nil, nil
	}

	var res = new(models.EmbedError)

	switch v := in.(type) {
	case json.Marshaler:
		res.Value = v
	case interface{ Unwrap() error }:
		var unwrapped = v.Unwrap()

		res.Value = strings.ReplaceAll(in.Error(), unwrapped.Error(), "%w")

		item, err := marshalJSONEmbedError(unwrapped)
		if err != nil {
			return nil, err
		}

		if item != nil {
			res.Embed = append(res.Embed, item)
		}
	case interface{ Unwrap() []error }:
		for _, sub := range v.Unwrap() {
			item, err := marshalJSONEmbedError(sub)
			if err != nil {
				return nil, err
			}

			if item != nil {
				res.Embed = append(res.Embed, item)
			}
		}
	default:
		res.Value = in.Error()
	}

	return res, nil
}

func (e *Error) UnmarshalJSON(data []byte) (err error) {
	var raw models.Error

	err = json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	*e = Error{
		Code:    raw.Code,
		Message: raw.Message,
		Fields:  raw.Fields,
		Codes:   raw.Codes,
	}

	if raw.Embed != nil {
		e.Embed, err = unmarshalJSONEmbedError(raw.Embed)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalJSONEmbedError(raw *models.EmbedError) (error, error) {
	var errs []error

	for _, item := range raw.Embed {
		switch value := item.Value.(type) {
		case map[string]any:
			data, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}

			var std *Error

			err = json.Unmarshal(data, &std)
			if err != nil {
				return nil, err
			}

			if len(std.Code) == 0 {
				break
			}

			if std.Code == "undefined" {
				std.Code = ""
			}

			errs = append(errs, std)
		case string:
			if strings.Contains(value, "%w") && len(item.Embed) == 1 {
				uErr, err := unmarshalJSONEmbedError(item)
				if err != nil {
					return nil, err
				}

				std, ok := uErr.(*Error)
				if ok {
					return fmt.Errorf(value, std), nil
				}

				errs = append(errs, fmt.Errorf(value, uErr))

				break
			}

			errs = append(errs, errors.New(value))
		case nil:
			uErr, err := unmarshalJSONEmbedError(item)
			if err != nil {
				return nil, err
			}

			errs = append(errs, uErr)
		}
	}

	if len(errs) == 1 {
		return errs[0], nil
	}

	return errors.Join(errs...), nil
}

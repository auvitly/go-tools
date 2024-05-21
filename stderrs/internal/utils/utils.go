package utils

func Unwrap(err error) []error {
	switch v := err.(type) {
	case interface{ Unwrap() error }:
		return []error{v.Unwrap()}
	case interface{ Unwrap() []error }:
		return v.Unwrap()
	default:
		return nil
	}
}

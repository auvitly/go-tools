package memory

import "errors"

type ProtectMode uint32

const (
	ProtectModeRead ProtectMode = 1 << iota
	ProtectModeReadWrite
)

var ErrUnsupportedProtectMode = errors.New("unsupported protect mode")

package kit

import (
	"github.com/auvitly/go-tools/lab/random"
	"math/rand"
	"time"
)

var (
	// Now - timestamp within the current locale.
	// The method for generating timestamps is assumed relative to this object.
	// If necessary, you can use the Add method to define time boundaries.
	Now    = time.Now().Local()
	IPv4   = random.IPv4()
	IPv6   = random.IPv6()
	Int    = random.Integer[int]()
	Int8   = random.Integer[int8]()
	Int16  = random.Integer[int16]()
	Int32  = random.Integer[int32]()
	Int64  = random.Integer[int64]()
	Uint   = random.Integer[uint]()
	Uint8  = random.Integer[uint8]()
	Uint16 = random.Integer[uint16]()
	Uint32 = random.Integer[uint32]()
	Uint64 = random.Integer[uint64]()
	Byte   = random.Integer[byte]()
	Bytes  = random.Bytes(rand.Intn(64))
	String = random.String(rand.Intn(64))
)

package function_test

import (
	"github.com/auvitly/go-tools/nuclear/function"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/arch/x86/x86asm"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

func TestReplace(t *testing.T) {
	var ts = time.Now()

	var (
		oldTimeFunc func() time.Time
		newTimeFunc = func() time.Time {
			return ts
		}
	)

	for i := 0; i < 100; i++ {
		patch := function.Replace(time.Now, newTimeFunc, &oldTimeFunc)
		require.NotNil(t, patch)

		assert.Equal(t, time.Now(), ts, "equal")

		patch.Unpatch()
		time.Sleep(time.Second)
		assert.NotEqual(t, time.Now(), ts, "not equal")
	}
}

func TestDecode(t *testing.T) {
	var (
		value int
		to    = uintptr(unsafe.Pointer(&value))
	)

	inst, err := x86asm.Decode([]byte{
		0xFF, 0x25, 0, 0, 0, 0, // ff 25 00 00 00 00   jmp    QWORD PTR
		byte(to), // movabs rdx, to
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56),
	}, strconv.IntSize)
	require.NoError(t, err)

	t.Log(inst)

	var raw = []byte{
		0xE9,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24), // JMP to
	}

	inst, err = x86asm.Decode(raw, 32)
	require.NoError(t, err)

	t.Log(inst)
}

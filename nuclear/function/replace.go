package function

import (
	"errors"
	"fmt"
	"github.com/auvitly/go-tools/nuclear/memory"
	"golang.org/x/arch/x86/x86asm"
	"strconv"
	"sync"
	"unsafe"
)

type sections struct {
	Header          uintptr
	Body            uintptr
	Footer          uintptr
	EOF             uintptr
	JBEInstruction  x86asm.Inst
	JBEAddress      uintptr
	CallAddress     uintptr
	CallInstruction x86asm.Inst
	JMPAddress      uintptr
}

var mu sync.Mutex

func findInstruction(ptr uintptr, size int, desired x86asm.Op, allowed ...x86asm.Op) (
	instr *x86asm.Inst, addr uintptr, _ error) {
	var (
		code = memory.As(ptr, size)
		pos  = 0
	)

loop:
	for {
		if pos > (size - 10) {
			ptr += uintptr(pos)
			code, pos = memory.As(ptr, size), 0
		}

		instruction, err := x86asm.Decode(code[pos:], strconv.IntSize)
		if err != nil {
			return nil, 0, fmt.Errorf("x86asm.Decode: %w", err)
		}

		if instruction.Op == desired {
			return &instruction, ptr + uintptr(pos), nil
		}

		for _, allow := range allowed {
			if instruction.Op == allow {
				pos += instruction.Len

				continue loop
			}
		}

		return nil, 0, fmt.Errorf("unexpected instruction found: %v desired: %v", instruction, desired)
	}
}

func inspectFunction(ptr uintptr) (res sections, err error) {
	inst, addr, err := findInstruction(ptr, 16, x86asm.JBE, x86asm.LEA, x86asm.CMP, x86asm.NOP, x86asm.INT)
	if err != nil {
		return res, fmt.Errorf("findInstruction[0]: %w", err)
	}

	res.JBEInstruction = *inst
	res.JBEAddress = addr

	rel, ok := inst.Args[0].(x86asm.Rel)
	if !ok {
		return res, errors.New("inspectFunction")
	}

	res.Body = addr + uintptr(inst.Len)
	res.Footer = res.Body + uintptr(rel)

	inst, addr, err = findInstruction(res.Footer, 128, x86asm.CALL, x86asm.MOV, x86asm.CMP, x86asm.NOP, x86asm.INT)
	if err != nil {
		return res, fmt.Errorf("findInstruction[1]: %w", err)
	}

	res.CallInstruction, res.CallAddress = *inst, addr

	inst, addr, err = findInstruction(addr+uintptr(inst.Len), 128, x86asm.JMP, x86asm.MOV, x86asm.CMP, x86asm.NOP, x86asm.INT)
	if err != nil {
		return res, fmt.Errorf("findInstruction[2]: %w", err)
	}

	rel, ok = inst.Args[0].(x86asm.Rel)
	if !ok && addr+uintptr(rel)+uintptr(inst.Len) != ptr {
		return res, errors.New("inspectFunction")
	}

	res.Header = ptr
	res.EOF = addr + uintptr(inst.Len) + 1
	res.JMPAddress = addr

	return res, nil
}

func Replace[T any](tg, rp T, oldTo ...*T) *Patch[T] {
	mu.Lock()
	defer mu.Unlock()

	sec, err := inspectFunction(**(**uintptr)(unsafe.Pointer(&tg)))
	if err != nil {
		panic(fmt.Sprintf("inspectFunction: %v", err))
	}

	var (
		newFunc = *(*uintptr)(unsafe.Pointer(&rp))
		oldFunc = *(*uintptr)(unsafe.Pointer(&tg))
	)

	var (
		beforeJBE  = memory.As(sec.Header, int(sec.JBEAddress-sec.Header))
		beforeCall = memory.As(sec.Footer, int(sec.CallAddress-sec.Footer))
		afterCall  = memory.As(
			sec.CallAddress+uintptr(sec.CallInstruction.Len),
			int(sec.JMPAddress-sec.CallAddress-uintptr(sec.CallInstruction.Len)),
		)
	)

	rel, ok := sec.CallInstruction.Args[0].(x86asm.Rel)
	if !ok {
		panic("call of split_stack is not relative")
	}

	var newHeader = make([]byte, sec.JBEAddress-sec.Header, sec.Body-sec.Header)

	for i := range newHeader {
		newHeader[i] = 0x90
	}

	switch sec.JBEInstruction.Len {
	case 2:
		newHeader = append(newHeader, 0xEB) // JMP SHORT
	case 6:
		newHeader = append(newHeader, 0x90, 0xE9) // JMP NEAR
	default:
		panic(fmt.Sprintf("unsupported JBE instruction %v", sec.JBEInstruction))
	}

	var newFooter = append(moveDX(newFunc), 0xFF, 0x22) // JMP to new impls

	var (
		oldHeader = memory.Clone(sec.Header, len(newHeader))
		oldFooter = memory.Clone(sec.Footer, len(newFooter))
	)

	var (
		jmpBytes    = append(moveDX(oldFunc), jumpFar(sec.Body)...)
		callAddr    = sec.CallAddress + uintptr(sec.CallInstruction.Len) + uintptr(rel)
		callBytes   = append(moveDX(callAddr), 0xFF, 0xD2) // call more stack_ctx
		size        = len(beforeJBE) + len(beforeCall) + len(callBytes) + len(afterCall) + len(jmpBytes) + 4
		originTramp = make([]byte, 0, size)
	)

	originTramp = append(originTramp, beforeJBE...)              //
	originTramp = append(originTramp, 0x76, byte(len(jmpBytes))) // JBE rel8
	originTramp = append(originTramp, jmpBytes...)               // JMP FAR TO ORIGIN BODY
	originTramp = append(originTramp, beforeCall...)             //
	originTramp = append(originTramp, callBytes...)              //
	originTramp = append(originTramp, 0xEB, byte(-size))         // JMP rel8

	memory.SetProtect(uintptr(unsafe.Pointer(&originTramp[0])), size, memory.ProtectModeReadWrite)

	var (
		entrypoint = &originTramp[0]
		fn         = &entrypoint
		proxyFn    = *(*T)(unsafe.Pointer(&fn))
	)

	for _, item := range oldTo {
		*item = proxyFn
	}

	doPatch(sec.Header, sec.Footer, sec.EOF, newHeader, newFooter)

	var res = NewPatch(proxyFn, rp).
		WithUnpatch(func() {
			mu.Lock()
			defer mu.Unlock()

			doPatch(sec.Header, sec.Footer, sec.EOF, oldHeader, oldFooter)
		})

	return res
}

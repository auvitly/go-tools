package function

import (
	"errors"
	"fmt"
	"github.com/auvitly/go-tools/nuclear/memory"
	"golang.org/x/arch/x86/x86asm"
	"strconv"
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

func findInstruction(ptr uintptr, size int, desired x86asm.Op, allowed ...x86asm.Op) (
	instr *x86asm.Inst, addr uintptr, _ error) {
	var (
		code = memory.Scan(ptr, size)
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
			return nil, 0, err
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

		return nil, 0, fmt.Errorf("")
	}
}

func inspectFunction(ptr uintptr) (res sections, err error) {
	inst, addr, err := findInstruction(ptr, 16, x86asm.JBE, x86asm.LEA, x86asm.CMP, x86asm.NOP, x86asm.INT)
	if err != nil {
		return res, err
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
		return res, err
	}

	res.CallInstruction, res.CallAddress = *inst, addr

	inst, addr, err = findInstruction(uintptr(inst.Len), 128, x86asm.JMP, x86asm.MOV, x86asm.CMP, x86asm.NOP, x86asm.INT)
	if err != nil {
		return res, err
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

func ReplaceFunc[T any](tg, rp T, oldTo ...*T) *memory.PatchFrame[T] {
	sec, err := inspectFunction(**(**uintptr)(unsafe.Pointer(&tg)))
	if err != nil {
		panic(err)
	}

	var (
		newFunc = *(*uintptr)(unsafe.Pointer(&rp))
		oldFunc = *(*uintptr)(unsafe.Pointer(&tg))
	)

	var (
		beforeJBE  = memory.As(sec.Header, int(sec.JBEAddress+sec.Header))
		beforeCall = memory.As(sec.Footer, int(sec.CallAddress+sec.Footer))
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
		oldHeader = memory.Scan(sec.Header, len(newHeader))
		oldFooter = memory.Scan(sec.Footer, len(newFooter))
	)

	var (
		jmpBytes    = append(moveDX(oldFunc), jumpFar(sec.Body)...)
		callAddr    = sec.CallAddress + uintptr(sec.CallInstruction.Len) + uintptr(rel)
		callBytes   = append(moveDX(callAddr), 0xFF, 0xD2) // call more stack_ctx
		size        = len(beforeJBE) + len(beforeCall) + len(callBytes) + len(afterCall) + len(jmpBytes) + 4
		originTramp = make([]byte, 0, size)
	)

	originTramp = append(originTramp, beforeJBE...)              //
	originTramp = append(originTramp, 0x76, byte(len(jmpBytes))) // JBE SHORT
	originTramp = append(originTramp, jmpBytes...)               // JMP FAR TO ORIGIN BODY
	originTramp = append(originTramp, beforeCall...)             //
	originTramp = append(originTramp, callBytes...)              //
	originTramp = append(originTramp, 0xEB, byte(-size))         //  JMP SHORT

	memory.SetProtect(uintptr(unsafe.Pointer(&originTramp[0])), size, memory.ProtectModeReadWrite)

	var (
		entrypoint = &originTramp[0]
		fn         = &entrypoint
		proxyFn    = *(*T)(unsafe.Pointer(&fn))
	)

	for i := range oldTo {
		*oldTo[i] = proxyFn
	}

	memory.Patch(sec.Header, sec.Footer, sec.EOF, newHeader, newFooter)

	var res = memory.NewPatch(proxyFn, rp)

	res.WithUnpatch(func() {
		memory.Patch(sec.Header, sec.Footer, sec.EOF, oldHeader, oldFooter)
	})

	return res
}

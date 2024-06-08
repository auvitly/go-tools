//go:build darwin && 386
// +build darwin,386

#include "go_asm.h"
#include "textflag.h"

// func taskSelfTrap() (ret uint32)
TEXT ·taskSelfTrap(SB), $8-0
        PUSHQ AX
        MOVL $(0x1000000+28), AX // task_self_trap
        SYSCALL
        MOVL AX, ret+0(FP)
        POPQ AX
        RET

// func vmProtect(targetTask uint32, address uintptr, size int, setMaximum, newProtection uint32) (ret uint32)
TEXT ·vmProtect(SB), $56-40
        PUSHQ AX
        PUSHQ DI
        PUSHQ SI
        PUSHQ DX
        PUSHQ R10
        PUSHQ R8
        PUSHQ R9
        MOVL  targetTask+0(FP), DI
        MOVL  address+4(FP), SI
        MOVL  size+8(FP), DX
        MOVL  set_maximum+16(FP), R10
        MOVL  new_protection+20(FP), R8
        XORL  R9, R9
        XORQ  $(0x1000000+14), AX // mach_vm_protect
        SYSCALL
        MOVL AX, ret+24(FP)
        POPQ R9
        POPQ R8
        POPQ R10
        POPQ DX
        POPQ SI
        POPQ DI
        POPQ AX
        RET
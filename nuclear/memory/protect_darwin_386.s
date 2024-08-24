//go:build darwin && 386
// +build darwin,386

#include "go_asm.h"
#include "textflag.h"

// func taskSelfTrap() (ret uint32)
TEXT ·taskSelfTrap(SB), $4-0
    PUSHL AX
    MOVL  $(0x1000000+28), AX // task_self_trap
    SYSCALL
    MOVL AX, ret+0(FP)
    POPQ AX
    RET

// func vmProtect(targetTask uint32, address uintptr, size int, setMaximum, newProtection uint32) (ret uint32)
// 4 + 4 + 8 + 4 + 4 = +4
TEXT ·vmProtect(SB), $28-24
        PUSHL AX
        PUSHL DI
        PUSHL SI
        PUSHL DX
        PUSHL R10
        PUSHL R8
        PUSHL R9
        MOVL  targetTask+0(FP), DI
        MOVL  address+4(FP), SI
        MOVL  size+8(FP), DX
        MOVL  set_maximum+16(FP), R10
        MOVL  new_protection+20(FP), R8
        XORL  R9, R9
        MOVL  $(0x1000000+14), AX // mach_vm_protect
        SYSCALL
        MOVL AX, ret+24(FP)
        POPL R9
        POPL R8
        POPL R10
        POPL DX
        POPL SI
        POPL DI
        POPL AX
        RET
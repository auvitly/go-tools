package function

/*
	https://www.cs.uaf.edu/2009/fall/cs441/lecture/09_15_arguments.html
	https://www.felixcloutier.com/x86/jmp
*/

// moveDX: MOV RDX, to
// OpCode: 0x48, 0xBA
func moveDX(to uintptr) []byte {
	return []byte{
		0x48, 0xBA, // "0x48 0xBA <64-bit constant>" loads a 64-bit constant into register 2.
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56), // movabs rdx, to
	}
}

// jumpFar: JMP [RIP]
// OpCode: 0xFF, 0x25, 0x00, 0x00, 0x00, 0x00
func jumpFar(to uintptr) []byte {
	return []byte{
		0xFF, 0x25, 0x00, 0x00, 0x00, 0x00, // ff 25 00 00 00 00   jmp    QWORD PTR
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56),
	}
}

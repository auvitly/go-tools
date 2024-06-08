package function

func moveDX(to uintptr) []byte {
	return []byte{
		0x48, 0xBA,
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

func jumpFar(to uintptr) []byte {
	return []byte{
		0xFF, 0x25, 0, 0, 0, 0, // JMP qword ptr [6]
		byte(to), // movabs rdx, to
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56),
	}
}

package function

// moveDX: MOV EDX, to
func moveDX(to uintptr) []byte {
	return []byte{
		0xBA, // ?
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24), // mov edx, to
	}
}

// jumpFar: JMP to
func jumpFar(to uintptr) []byte {
	return []byte{
		0xE9,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24), // JMP to
	}
}

package lab

import "time"

var (
	// Now - timestamp within the current locale.
	// The method for generating timestamps is assumed relative to this object.
	// If necessary, you can use the Add method to define time boundaries.
	Now = time.Now().Local()
)

package general

import "os"

// the operator
func Operation() {
	err := general()
	if err != nil {
		fatalError(err)
		os.Exit(1)
	}
}

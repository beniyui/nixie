package general

import "os"

// the main entry point
func Operation() {
	err := general()
	if err != nil {
		fatalError(err)
		return
	}
	os.Exit(0)
}

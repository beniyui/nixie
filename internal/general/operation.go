package general

import "os"

// the main entry point
func Operation() {
	defer fatalError("君がこれをみているということはプログラムがパニックを起こした？")
	err := general()
	if err != nil {
		fatalError(err)
		os.Exit(1)
		return
	}
	os.Exit(0)
}

package general

// the main entry point
func Entry() {
	err := main()
	defer func() {
		if err != nil {
			fatalError(err)
		}
	}()
}

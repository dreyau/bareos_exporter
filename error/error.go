package error

// Check panics if an error occurred
func Check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

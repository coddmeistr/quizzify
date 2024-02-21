package errs

func NoErrors(errs ...error) bool {
	for _, err := range errs {
		if err != nil {
			return false
		}
	}
	return true
}

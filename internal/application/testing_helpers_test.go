package application

// OverrideHashFuncForTests temporarily replaces the password hashing function.
// The returned restore callback must be invoked to put the original function
// back in place.
func OverrideHashFuncForTests(fn func([]byte, int) ([]byte, error)) func() {
	original := generateFromPassword
	generateFromPassword = fn
	return func() {
		generateFromPassword = original
	}
}

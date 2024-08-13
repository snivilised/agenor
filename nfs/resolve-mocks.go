package nfs

// AbsFunc signature of function used to obtain the absolute representation of
// a path.
type AbsFunc func(path string) (string, error)

// Abs function invoker, allows a function to be used in place where
// an instance of an interface would be expected.
func (f AbsFunc) Abs(path string) (string, error) {
	return f(path)
}

// HomeUserFunc signature of function used to obtain the user's home directory.
type HomeUserFunc func() (string, error)

// Home function invoker, allows a function to be used in place where
// an instance of an interface would be expected.
func (f HomeUserFunc) Home() (string, error) {
	return f()
}

// ResolveMocks, used to override the internal functions used
// to resolve the home path (os.UserHomeDir) and the abs path
// (filepath.Abs). In normal usage, these do not need to be provided,
// just used for testing purposes.
type ResolveMocks struct {
	HomeFunc HomeUserFunc
	AbsFunc  AbsFunc
}

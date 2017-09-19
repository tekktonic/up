package main;

type upError struct {
	s string;
}

func (e *upError) Error() string {
	return e.s;
}
func NewUpError(s string) error {
	return &upError{s}
}

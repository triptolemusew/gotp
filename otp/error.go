package otp

type NotFound string

func (e NotFound) Error() string {
	return string(e)
}

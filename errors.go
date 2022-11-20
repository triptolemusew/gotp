package main

type GotpError struct{}

func (e *GotpError) Error() string {
	return "Wrong password"
}

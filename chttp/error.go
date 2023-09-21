package chttp

import "fmt"

type UnexpectedStatusError struct {
	StatusCode int
	Status     string
	Data       []byte
}

func (e UnexpectedStatusError) Error() string {
	return fmt.Sprintf("unexpected HTTP status: %d %s body:%s", e.StatusCode, e.Status, string(e.Data))
}

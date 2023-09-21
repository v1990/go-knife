package chttp

import (
	"io"
	"net/http"

	json "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"golang.org/x/exp/slices"
)

type Response struct {
	err error
	*http.Response
}

func makeResponse(resp *http.Response, err error) *Response {
	return &Response{
		Response: resp,
		err:      err,
	}
}

type ErrChecker interface {
	Err() error
}

func (r *Response) Err() error {
	if r == nil {
		return errors.New("response is nil")
	}
	if r.err != nil {
		return r.err
	}
	if r.Response == nil {
		return errors.New("response is nil")
	}
	return nil
}

func (r *Response) Close() {
	if r == nil || r.Response == nil || r.Body == nil {
		return
	}
	_ = r.Body.Close()
}

func (r *Response) CheckStatus(expectStatuses ...int) *Response {
	if r.Err() != nil {
		return r
	}

	if slices.Index(expectStatuses, r.StatusCode) < 0 {
		defer cleanResponse(r.Response)
		body, _ := io.ReadAll(io.LimitReader(r.Body, 4096))

		r.err = &UnexpectedStatusError{
			StatusCode: r.StatusCode,
			Status:     r.Status,
			Data:       body,
		}
	}
	return r
}

//func (r *Response) UnmarshalBody(status int,out any) error {
//	if err := r.Err(); err != nil {
//		return err
//	}
//	// TODO: support more content types
//	return r.UnmarshalBodyJSON(status,out)
//}

func (r *Response) UnmarshalBodyJSON(out any) error {
	if err := r.Err(); err != nil {
		return err
	}

	defer cleanResponse(r.Response)

	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return errors.Wrap(err, "unmarshal response err")
	}
	switch t := out.(type) {
	case ErrChecker:
		return errors.WithStack(t.Err())
	}
	return nil
}

package chttp

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	json "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
)

// HTTPEntity represents a payload that can be included in an outgoing HTTP request.
type HTTPEntity interface {
	Bytes() ([]byte, error)
	Mime() string
}

// Form represents a form payload that can be included in an outgoing HTTP request.
type Form struct {
	url.Values
}

type jsonEntity struct {
	Val interface{}
}

// NewJSONEntity creates a new HTTPEntity that will be serialized into JSON.
func NewJSONEntity(v interface{}) HTTPEntity {
	return &jsonEntity{Val: v}
}

func (e *jsonEntity) Bytes() ([]byte, error) {
	return json.Marshal(e.Val)
}

func (e *jsonEntity) Mime() string {
	return "application/json"
}

func marshalRequestData(in any) (body io.Reader, mime string, err error) {
	switch v := in.(type) {
	case nil:
		return nil, "", nil
	case Form:
		body = strings.NewReader(v.Encode())
		mime = "application/x-www-form-urlencoded"
		return body, mime, nil
	case HTTPEntity:
		data, err := v.Bytes()
		mime = v.Mime()
		return bytes.NewReader(data), mime, err
	case proto.Message:
		data, err := proto.Marshal(v)
		mime = "application/x-protobuf"
		return bytes.NewReader(data), mime, err
	case []byte:
		mime = http.DetectContentType(v)
		return bytes.NewReader(v), mime, nil
	case string:
		return strings.NewReader(v), "", nil
	case io.Reader: // *os.File, *bytes.Buffer, *bytes.Reader, *strings.Reader...
		return v, "", nil
	}
	// default:

	// TODO: support more content types
	data, err := json.Marshal(in)
	if err == nil {
		body = bytes.NewReader(data)
		mime = "application/json"
	}
	return body, mime, err
}

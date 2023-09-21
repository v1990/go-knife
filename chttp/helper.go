package chttp

import (
	"io"
	"net/http"
)

// CheckHeadStatusOK 发送 HEAD 请求，检查 reqURL 的 status code 是否为 200
func CheckHeadStatusOK(reqURL string) bool {
	resp, err := http.DefaultClient.Head(reqURL)
	if err != nil {
		return false
	}
	defer cleanResponse(resp)

	return resp.StatusCode == http.StatusOK
}

func cleanResponse(r *http.Response) {
	if r == nil || r.Body == nil {
		return
	}
	defer r.Body.Close()
	_, _ = io.Copy(io.Discard, r.Body)
}

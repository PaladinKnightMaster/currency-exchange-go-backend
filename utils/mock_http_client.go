package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type MockClient struct {
	ResponseBody string
	StatusCode   int
	Err          error
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	resp := &http.Response{
		StatusCode: m.StatusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(m.ResponseBody)),
	}
	return resp, nil
}

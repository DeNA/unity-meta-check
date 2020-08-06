package github

import "net/http"

type HttpFunc func(request *http.Request) (*http.Response, error)

func NewHttp() HttpFunc {
	return func(request *http.Request) (*http.Response, error) {
		return http.DefaultClient.Do(request)
	}
}

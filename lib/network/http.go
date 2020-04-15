package network

import (
	"net/http"
)

func HttpStatusCode(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	resp.Body.Close()
	return resp.StatusCode, nil
}

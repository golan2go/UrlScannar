package reader

import (
	"fmt"
)

const ERROR = -1

type responseResult struct {
	url        string
	statusCode int
	length     int
	error      error
}

func errorResult(url string, error error) *responseResult {
	return &responseResult{
		url:        url,
		error:      error,
		statusCode: ERROR,
	}
}

func okResult(url string, statusCode int, length int) *responseResult {
	return &responseResult{
		url:        url,
		statusCode: statusCode,
		length:     length,
	}
}

func (res *responseResult) String() string {
	return fmt.Sprintf("%s\t|\t%d\t|\t%s\n", res.url, res.statusCode, res.text())
}

func (res *responseResult) text() string {
	if res.statusCode == ERROR {
		return res.error.Error()
	} else {
		return fmt.Sprintf("%d", res.length)
	}

}

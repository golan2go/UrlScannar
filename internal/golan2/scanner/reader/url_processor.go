package reader

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type urlProcessor struct {
	index   int
	urlChan chan string
	resChan chan responseResult
	wg      *sync.WaitGroup
}

func newUrlProcessor(index int, urlChan chan string, resChan chan responseResult, wg *sync.WaitGroup) *urlProcessor {
	fmt.Println("urlProcessor - ", index)

	return &urlProcessor{
		index:   index,
		urlChan: urlChan,
		resChan: resChan,
		wg:      wg,
	}
}

func (up *urlProcessor) asyncProcess() {
	go up.process()
}

func (up *urlProcessor) process() {
	defer up.done()
	count := 0
	for url := range up.urlChan {
		result := up.invokeUrl(url)
		up.resChan <- result
		count++
		if count%10 == 0 {
			fmt.Printf("[P][%d] %d messages\n", up.index, count)
		}
	}
	fmt.Printf("[%d] DONE PROCESSING: %d URLs were processed\n", up.index, count)
}

func (up *urlProcessor) invokeUrl(url string) responseResult {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 5,
	}

	response, err := client.Get(url)
	if err != nil {
		return *errorResult(url, err)
	}

	length, err := up.responseLength(response)
	if err != nil {
		return *errorResult(url, err)
	}
	return *okResult(url, response.StatusCode, length)
}

func (up *urlProcessor) responseLength(response *http.Response) (int, error) {
	res := 0
	buf := make([]byte, 1000)
	body := response.Body
	var read int
	var err error

	for err == nil {
		begin := time.Now()
		read, err = body.Read(buf)
		res += read
		if time.Now().UnixNano()-begin.UnixNano() > 5000000000 {
			return -1, errors.New("BodyReadTimeout")
		}
	}

	if err.Error() != "EOF" {
		return -1, err
	}
	return res, nil
}

func (up *urlProcessor) done() {
	up.wg.Done()
	fmt.Printf("urlProcessor - process - END %d    wg=[%v]\n", up.index, up.wg)
}

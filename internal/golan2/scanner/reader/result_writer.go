package reader

import (
	"fmt"
	"os"
	"sync"
)

type resultWriter struct {
	fileName string
	resChan  chan responseResult
	wg       *sync.WaitGroup
}

func newResultWriter(fileName string, resChan chan responseResult, wg *sync.WaitGroup) *resultWriter {
	return &resultWriter{
		fileName: fileName,
		resChan:  resChan,
		wg:       wg,
	}
}

func (rw *resultWriter) asyncWrite() {
	rw.wg.Add(1)
	go rw.write()
}

func (rw *resultWriter) write() {
	defer rw.wg.Done()

	file, err := os.Create(rw.fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	count := 0
	for res := range rw.resChan {
		_, err := file.Write([]byte(res.String()))
		if err != nil {
			fmt.Println(err)
		}
		count++
		if count%10 == 0 {
			fmt.Printf("[W] %d messages\n", count)
		}
	}
	fmt.Printf("DONE WRITING TO FILE [%d messages]\n", count)
}

package reader

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type urlFeeder struct {
	fileName string
	urlChan  chan string
	done     *sync.WaitGroup
}

func newUrlFeeder(fileName string, outChannel chan string, done *sync.WaitGroup) *urlFeeder {
	return &urlFeeder{
		fileName: fileName,
		urlChan:  outChannel,
		done:     done,
	}
}

func (uf *urlFeeder) asyncFeed() {
	go uf.feed()
}

func (uf *urlFeeder) feed() {
	file, err := os.Open(uf.fileName)
	if err != nil {
		fmt.Println("ERROR: ", uf, err)
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		uf.urlChan <- text
		count++
		if count%10 == 0 {
			fmt.Printf("[F] %d messages\n", count)
		}
	}

	fmt.Printf("DONE FEEDING: %d messages were ingested to the [urlChan]\n", count)

	uf.done.Done()
}

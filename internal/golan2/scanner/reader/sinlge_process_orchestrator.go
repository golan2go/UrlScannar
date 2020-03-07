package reader

import (
	"sync"
)

//
//
// urlFeeder[1]                             urlProcessor[*]                                resultWriter[1]   ====> FILE
//            \\                        //                   \\                          //
//             \\                      //                     \\                        //
//              \\                    //                       \\                      //
//                >==== urlChan  ====>                           >=====  resChan  ====>
//
//
//
//

//
// SingleProcessOrchestrator will read from an input file all URLs and write the results to an output file
// - A single worker reads from the file and populates the [urlChan] channel.
// - Several workers read from the [urlChan] channel and write the results to the [resChan] channel.
// - A single worker reads the [resChan] channel and writes to the output file.
//
type SingleProcessOrchestrator struct {
	inputFile   string
	outputFile  string
	parallelism int
}

func NewSingleProcessOrchestrator(inputFile string, outputFile string, parallelism int) *SingleProcessOrchestrator {
	return &SingleProcessOrchestrator{
		inputFile:   inputFile,
		outputFile:  outputFile,
		parallelism: parallelism,
	}
}

func (sp *SingleProcessOrchestrator) Run() {
	urlChan := make(chan string, sp.parallelism)
	resChan := make(chan responseResult, sp.parallelism)

	waitFeeder := sync.WaitGroup{}
	waitFeeder.Add(1)
	newUrlFeeder(sp.inputFile, urlChan, &waitFeeder).asyncFeed()
	go closeUrlChannelWhenDone(urlChan, &waitFeeder)

	var waitProcessors sync.WaitGroup
	for i := 0; i < sp.parallelism; i++ {
		waitProcessors.Add(1)
		newUrlProcessor(i, urlChan, resChan, &waitProcessors).asyncProcess()
	}

	var waitWriter sync.WaitGroup
	newResultWriter(sp.outputFile, resChan, &waitWriter).asyncWrite()

	waitProcessors.Wait()
	close(resChan)

	waitWriter.Wait()
}

func closeUrlChannelWhenDone(urlChan chan string, waitFeeder *sync.WaitGroup) {
	waitFeeder.Wait()
	close(urlChan)
	println("urlChan CLOSED!")
}

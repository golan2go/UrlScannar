package main

import (
	"../internal/golan2/scanner/reader"
)

const inputFile = "/tmp/urls.txt"
const outputFile = "/tmp/results.txt"
const parallelism int = 2

func main() {
	reader.NewSingleProcessOrchestrator(inputFile, outputFile, parallelism).Run()
}

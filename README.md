# Go Queue
Multi backend queues for Golang


[![](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/nickalie/go-queue)
[![](https://circleci.com/gh/nickalie/go-queue.png?circle-token=12e613097830e35d6a5426361f2783cd4331f709)](https://circleci.com/gh/nickalie/go-queue)
[![codecov](https://codecov.io/gh/nickalie/go-queue/branch/master/graph/badge.svg)](https://codecov.io/gh/nickalie/go-queue)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/3ab58d8ce189430cac752b26465350d2)](https://www.codacy.com/app/nickalie/go-queue?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nickalie/go-queue&amp;utm_campaign=Badge_Grade)[![Go Report Card](https://goreportcard.com/badge/github.com/nickalie/go-queue)](https://goreportcard.com/report/github.com/nickalie/go-queue)

## Install

```go get -u github.com/nickalie/go-queue```

## Example of usage

```go
package main

import (
	"fmt"
	"github.com/nickalie/go-queue"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		go consumer(i + 1)
	}

	producer()
}

func producer() {
	i := 0
	for {
		i++
		queue.Put("messages", fmt.Sprintf("message %d", i))
		time.Sleep(time.Second)
	}
}

func consumer(index int) {
	for {
		var message string
		queue.Get("messages", &message)

		fmt.Printf("Consumer %d got a message: %s\n", index, message)
		time.Sleep(2 * time.Second)
	}
}
```

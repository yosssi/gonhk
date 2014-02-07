# gonhk - NHK API Client in Go

[![Build Status](https://travis-ci.org/yosssi/gonhk.png?branch=master)](https://travis-ci.org/yosssi/gonhk)
[![GoDoc](https://godoc.org/github.com/yosssi/gonhk?status.png)](https://godoc.org/github.com/yosssi/gonhk)

gonhk implements a NHK API client.

## Installation

```sh
$ go get github.com/yosssi/gonhk
```

## Example

```go
package main

import (
	"fmt"
	"github.com/yosssi/gonhk"
	"time"
)

func main() {
	// Get an NHK API client. You have to pass your NHK API key to this function.
	client := gonhk.NewClient("Your NHK API Key")

	today := time.Now().Format("2006-01-02")

	// Call NHK Program List API.
	// http://api-portal.nhk.or.jp/doc_list-v1_con
	result1, err := client.ProgramList("v1", "130", "g1", today)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", result1)

	// Call NHK Program Genre API.
	// http://api-portal.nhk.or.jp/doc_genre-v1_con
	result2, err := client.ProgramGenre("v1", "130", "g1", "0000", today)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", result2)

	// Call NHK Program Info API.
	// http://api-portal.nhk.or.jp/doc_info-v1_con
	result3, err := client.ProgramInfo("v1", "130", "g1", "2014020702065")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", result3)

	// Call NHK Now On Air API.
	// http://api-portal.nhk.or.jp/doc_now-v1_con
	result4, err := client.NowOnAir("v1", "130", "g1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", result4)
}
```

## Documentation

[GoDoc](http://godoc.org/github.com/yosssi/gonhk)

## Note
情報提供:ＮＨＫ  
This program uses [NHK番組表API](http://api-portal.nhk.or.jp/).

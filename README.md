# Quickwit Go SDK
The Go SDK specifically for Quickwit search engine

## Installation

```bash
go get github.com/akbariandev/quickwit-gosdk
```

## How to use

#### Search Query

```go
package main

import (
    "fmt"
    quickwitgo "github.com/akbariandev/quickwit-gosdk"
)

func main() {
	client := quickwitgo.NewClient("http://localhost:7280")
	queryRequest := quickwitgo.SearchRequest{Query: "events:error"}
	response, err := client.Search("otel-traces-v0_7", queryRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
```
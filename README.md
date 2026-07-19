# Quickwit Go SDK

A Go SDK for the [Quickwit](https://quickwit.io) search engine.

## Installation

```bash
go get github.com/akbariandev/quickwit-gosdk
```

## Features

- [x] Search
- [x] Search Stream
- [x] Scroll Search
- [x] Delete Query
- [x] Index CRUD operations
- [x] Ingest API
- [x] Delete Tasks API

## Usage

### Creating a Client

```go
client := quickwitgosdk.NewClient("http://localhost:7280")

// With API key authentication
client := quickwitgosdk.NewClient("http://localhost:7280",
    quickwitgosdk.WithAPIKey("your-api-key"),
)

// With custom timeout
client := quickwitgosdk.NewClient("http://localhost:7280",
    quickwitgosdk.WithTimeout(30 * time.Second),
)

// With custom transport (e.g. for custom TLS)
client := quickwitgosdk.NewClient("http://localhost:7280",
    quickwitgosdk.WithTransport(customTransport),
)
```

### Search

```go
resp, err := client.Search("my-index", quickwitgosdk.SearchRequest{
    Query:           "events:error",
    DefaultOperator: "AND",
    MaxHits:         10,
    SearchFields:    []string{"message", "title"},
    SortByField:     &quickwitgosdk.SortByField{FieldName: "timestamp", Order: "desc"},
    SnippetFields:   &quickwitgosdk.SnippetRequest{FieldName: "message", MaxNumCharsPerFragment: 200},
    TagFilters:      []string{"tag1:value1", "tag2:value2"},
    Filter:          "status_code >= 400",
})

fmt.Printf("Found %d hits\n", resp.NumHits)
for _, hit := range resp.Hits {
    fmt.Printf("  %v\n", hit.JSON)
}
```

### Search Stream (Channel-based)

```go
ctx := context.Background()
respCh, errCh := client.SearchStream(ctx, "my-index", quickwitgosdk.SearchStreamRequest{
    SearchRequest: quickwitgosdk.SearchRequest{Query: "streaming"},
    FastField:     "timestamp",
})

for resp := range respCh {
    fmt.Printf("Batch: %d hits\n", resp.NumHits)
}
if err := <-errCh; err != nil {
    log.Fatal(err)
}
```

### Search Stream (Callback-based)

```go
err := client.SearchStreamWithCallback(ctx, "my-index", quickwitgosdk.SearchStreamRequest{
    SearchRequest: quickwitgosdk.SearchRequest{Query: "streaming"},
}, func(resp quickwitgosdk.SearchResponse) error {
    fmt.Printf("Batch: %d hits\n", resp.NumHits)
    return nil
})
```

### Scroll Search

```go
// Channel-based
respCh, errCh := client.SearchScroll(ctx, "my-index", quickwitgosdk.SearchScrollRequest{
    SearchRequest: quickwitgosdk.SearchRequest{Query: "all logs"},
    ScrollTTLSecs: 60,
})

// Callback-based
err := client.SearchScrollWithCallback(ctx, "my-index", req, func(resp quickwitgosdk.SearchResponse) error {
    fmt.Printf("Batch: %d hits\n", resp.NumHits)
    return nil
})
```

### Ingest

```go
// Ingest a batch of documents
resp, err := client.Ingest("my-index", []interface{}{
    map[string]interface{}{"title": "doc1", "body": "hello world"},
    map[string]interface{}{"title": "doc2", "body": "quickwit sdk"},
})
fmt.Printf("Persisted: %d, Failed: %d\n", resp.NumPersisted, resp.NumFailed)

// Ingest from a reader (NDJSON)
reader := strings.NewReader(`{"title":"doc1"}
{"title":"doc2"}`)
resp, err := client.IngestFromReader("my-index", reader)

// Force merge
err = client.ForceMerge("my-index")
```

### Index CRUD

```go
// Create
index, err := client.CreateIndex(quickwitgosdk.CreateIndexRequest{IndexID: "my-index"})

// Get
index, err := client.GetIndex("my-index")

// List all
indexes, err := client.ListIndexes()

// Delete (dry run)
resp, err := client.DeleteIndex("my-index", true)

// Delete (for real)
resp, err := client.DeleteIndex("my-index", false)

// Clear (remove all data, keep index)
err = client.ClearIndex("my-index")
```

### Delete by Query

```go
// Submit a delete task
resp, err := client.DeleteByQuery("my-index", quickwitgosdk.DeleteQueryRequest{
    Query:   "old_logs:true",
    Filter:  "timestamp < 1704067200",
})

// Check task status
task, err := client.GetDeleteTask("my-index", resp.TaskID)
fmt.Printf("Status: %s\n", task.Status) // "running", "success", "error", "cancelled"
```

## License

MIT

# Quickwit Go SDK

A Go SDK for the [Quickwit](https://quickwit.io) search engine.

## Installation

```bash
go get github.com/akbariandev/quickwit-gosdk
```

## Features

- [x] Search (with scroll support)
- [x] Delete Query
- [x] Index CRUD operations
- [x] Ingest API
- [x] Delete Tasks API

> **Note:** Quickwit removed the `/search/stream` and `/search/scroll` endpoints
> in v0.9 (PR #5886). Scroll is now a field (`scroll_ttl_secs`) on the regular
> search request. For large result sets, paginate using `start_offset`/`max_hits`.

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
    Format:          quickwitgosdk.FormatPrettyJSON,
})

fmt.Printf("Found %d hits\n", resp.NumHits)
for _, hit := range resp.Hits {
    fmt.Printf("  %v\n", hit.Fields)
}
```

### Scroll Search

Scroll is enabled by setting `ScrollTTLSecs` on a regular search request:

```go
scrollTTL := uint64(60) // seconds
resp, err := client.Search("my-index", quickwitgosdk.SearchRequest{
    Query:         "all logs",
    MaxHits:       1000,
    ScrollTTLSecs: &scrollTTL,
})
```

### Paginating Large Result Sets

Since the stream endpoint was removed, paginate using `start_offset`/`max_hits`:

```go
const pageSize = 100
offset := uint64(0)
for {
    resp, err := client.Search("my-index", quickwitgosdk.SearchRequest{
        Query:       "body:error",
        MaxHits:     pageSize,
        StartOffset: offset,
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, hit := range resp.Hits {
        fmt.Printf("  %v\n", hit.Fields)
    }
    if len(resp.Hits) < pageSize {
        break // no more results
    }
    offset += pageSize
}
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
// Create — version, index_id, and doc_mapping are required.
index, err := client.CreateIndex(quickwitgosdk.CreateIndexRequest{
    Version: "0.9",
    IndexID: "my-index",
    DocMapping: json.RawMessage(`{
        "mode": "lenient",
        "field_mappings": [
            {"name": "title",   "type": "text",   "stored": true, "indexed": true},
            {"name": "body",    "type": "text",   "stored": true, "indexed": true},
            {"name": "ts",      "type": "datetime","stored": true, "indexed": true, "fast": true}
        ],
        "tag_fields": [],
        "timestamp_field": "ts",
        "store_source": true
    }`),
    // Optional fields below — omitted when empty:
    IndexURI: "", // custom index storage URI
    IndexingSettings: json.RawMessage(`{
        "commit_timeout_secs": 60,
        "split_num_docs_target": 10000000,
        "merge_policy": {"type": "stable_log"},
        "resources": {"heap_size": "2 GB"}
    }`),
    IngestSettings: json.RawMessage(`{
        "min_shards": 1,
        "validate_docs": true
    }`),
    SearchSettings: json.RawMessage(`{
        "default_search_fields": ["title", "body"]
    }`),
    Retention: json.RawMessage(`{
        "period": "90 days",
        "schedule": "daily"
    }`),
})

// Get
index, err := client.GetIndex("my-index")

// List all
indexes, err := client.ListIndexes()

// Access index fields via IndexConfig:
//   indexes[0].IndexConfig.IndexID   // "my-index"
//   indexes[0].IndexConfig.IndexURI  // "file:///..."
//   indexes[0].IndexUID              // "my-index:01XYZ"
//   indexes[0].CreateTimestamp.Time  // creation time

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

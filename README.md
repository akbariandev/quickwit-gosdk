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

// With custom timeout
client := quickwitgosdk.NewClient("http://localhost:7280",
    quickwitgosdk.WithTimeout(30 * time.Second),
)

// With custom transport (e.g. for custom TLS or an auth-injecting reverse proxy)
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
    DocMapping: quickwitgosdk.DocMapping{
        Mode: "lenient",
        FieldMappings: []quickwitgosdk.FieldMapping{
            {Name: "title", Type: "text", Stored: true, Indexed: true},
            {Name: "body",  Type: "text", Stored: true, Indexed: true},
            {Name: "ts",    Type: "datetime", Stored: true, Indexed: true, Fast: &quickwitgosdk.FastField{Enabled: true}},
        },
        TagFields:      []string{},
        TimestampField: "ts",
        StoreSource:    true,
    },
    // Optional fields below — omitted when empty:
    IndexingSettings: quickwitgosdk.IndexingSettings{
        CommitTimeoutSecs:  60,
        SplitNumDocsTarget: 10000000,
        MergePolicy: &quickwitgosdk.MergePolicy{
            Type:        "stable_log",
            MergeFactor: 10,
        },
        Resources: &quickwitgosdk.IndexingResources{
            HeapSize: json.Number("2000000000"), // or "2 GB"
        },
    },
    IngestSettings: quickwitgosdk.IngestSettings{
        MinShards:    1,
        ValidateDocs: true,
    },
    SearchSettings: quickwitgosdk.SearchSettings{
        DefaultSearchFields: []string{"title", "body"},
    },
    Retention: &quickwitgosdk.Retention{
        Period:   "90 days",
        Schedule: "daily",
    },
})

// Get
index, err := client.GetIndex("my-index")

// List all
indexes, err := client.ListIndexes()

// Access typed index fields via IndexConfig:
//   index.IndexConfig.IndexID                    // "my-index"
//   index.IndexConfig.IndexURI                   // "file:///..."
//   index.IndexConfig.DocMapping.Mode            // "lenient"
//   index.IndexConfig.DocMapping.FieldMappings[0].Name  // "title"
//   index.IndexConfig.SearchSettings.DefaultSearchFields
//   index.IndexUID                                // "my-index:01XYZ"
//   index.CreateTimestamp.Time                    // creation time
//   index.Sources[0].SourceID                     // "_ingest-source"

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

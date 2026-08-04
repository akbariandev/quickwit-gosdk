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

## Package Structure

```
quickwit-gosdk/
├── client/           # HTTP client + options
├── types/            # Shared data types (doc mapping, hits, errors)
├── internal/httputil # Shared HTTP utilities (error hook, helpers)
├── search/           # Search API
├── index/            # Index management API
├── ingest/           # Document ingestion API
└── delete/           # Delete-by-query task API
```

## Usage

### Creating a Client

```go
import "github.com/akbariandev/quickwit-gosdk/client"

// Basic
c := client.New("http://localhost:7280")

// With custom timeout
c := client.New("http://localhost:7280",
    client.WithTimeout(30 * time.Second),
)

// With custom transport (e.g. for custom TLS or an auth-injecting reverse proxy)
c := client.New("http://localhost:7280",
    client.WithTransport(customTransport),
)
```

### Search

```go
import (
    "github.com/akbariandev/quickwit-gosdk/search"
    "github.com/akbariandev/quickwit-gosdk/types"
)

resp, err := search.Do(c, "my-index", search.Request{
    Query:           "events:error",
    DefaultOperator: "AND",
    MaxHits:         10,
    SearchFields:    []string{"message", "title"},
    SortByField:     &types.SortByField{FieldName: "timestamp", Order: "desc"},
    SnippetFields:   &types.SnippetRequest{FieldName: "message", MaxNumCharsPerFragment: 200},
    TagFilters:      []string{"tag1:value1", "tag2:value2"},
    Filter:          "status_code >= 400",
    Format:          types.FormatPrettyJSON,
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
resp, err := search.Do(c, "my-index", search.Request{
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
    resp, err := search.Do(c, "my-index", search.Request{
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
import "github.com/akbariandev/quickwit-gosdk/ingest"

// Send a batch of documents
resp, err := ingest.Send(c, "my-index", []interface{}{
    map[string]interface{}{"title": "doc1", "body": "hello world"},
    map[string]interface{}{"title": "doc2", "body": "quickwit sdk"},
})
fmt.Printf("Persisted: %d, Failed: %d\n", resp.NumPersisted, resp.NumFailed)

// Send from a reader (NDJSON)
reader := strings.NewReader(`{"title":"doc1"}
{"title":"doc2"}`)
resp, err := ingest.SendFromReader(c, "my-index", reader)

// Force merge
err = ingest.ForceMerge(c, "my-index")
```

### Index CRUD

```go
import (
    "encoding/json"

    "github.com/akbariandev/quickwit-gosdk/index"
    "github.com/akbariandev/quickwit-gosdk/types"
)

// Create — version, index_id, and doc_mapping are required.
metadata, err := index.Create(c, index.CreateRequest{
    Version: "0.9",
    IndexID: "my-index",
    DocMapping: types.DocMapping{
        Mode: "lenient",
        FieldMappings: []types.FieldMapping{
            {Name: "title", Type: "text", Stored: true, Indexed: true},
            {Name: "body",  Type: "text", Stored: true, Indexed: true},
            {Name: "ts",    Type: "datetime", Stored: true, Indexed: true, Fast: &types.FastField{Enabled: true}},
        },
        TagFields:      []string{},
        TimestampField: "ts",
        StoreSource:    true,
    },
    // Optional fields below — omitted when empty:
    IndexingSettings: types.IndexingSettings{
        CommitTimeoutSecs:  60,
        SplitNumDocsTarget: 10000000,
        MergePolicy: &types.MergePolicy{
            Type:        "stable_log",
            MergeFactor: 10,
        },
        Resources: &types.IndexingResources{
            HeapSize: json.Number("2000000000"), // or "2 GB"
        },
    },
    IngestSettings: types.IngestSettings{
        MinShards:    1,
        ValidateDocs: true,
    },
    SearchSettings: types.SearchSettings{
        DefaultSearchFields: []string{"title", "body"},
    },
    Retention: &types.Retention{
        Period:   "90 days",
        Schedule: "daily",
    },
})

// Get
metadata, err := index.Get(c, "my-index")

// List all
indexes, err := index.List(c)

// Access typed index fields via Config:
//   metadata.Config.IndexID                      // "my-index"
//   metadata.Config.IndexURI                     // "file:///..."
//   metadata.Config.DocMapping.Mode              // "lenient"
//   metadata.Config.DocMapping.FieldMappings[0].Name  // "title"
//   metadata.Config.SearchSettings.DefaultSearchFields
//   metadata.IndexUID                            // "my-index:01XYZ"
//   metadata.CreateTimestamp.Time                // creation time
//   metadata.Sources[0].SourceID                 // "_ingest-source"

// Delete (dry run — returns affected splits)
splits, err := index.Delete(c, "my-index", true)

// Delete (for real — returns empty array)
splits, err = index.Delete(c, "my-index", false)

// Clear (remove all data, keep index)
err = index.Clear(c, "my-index")
```

### Delete by Query

```go
import "github.com/akbariandev/quickwit-gosdk/delete"

// Submit a delete task
resp, err := delete.Submit(c, "my-index", delete.Request{
    Query:  "old_logs:true",
    Filter: "timestamp < 1704067200",
})

// Check task status
task, err := delete.GetTask(c, "my-index", resp.TaskID)
fmt.Printf("Status: %s\n", task.Status) // "running", "success", "error", "cancelled"
```

### Error Handling

All API errors are returned as `*types.Error` with the HTTP status code and
Quickwit's error message:

```go
import "github.com/akbariandev/quickwit-gosdk/types"

resp, err := search.Do(c, "my-index", search.Request{Query: "test"})
if err != nil {
    var apiErr *types.Error
    if errors.As(err, &apiErr) {
        fmt.Printf("Quickwit API error %d: %s\n", apiErr.StatusCode, apiErr.Message)
    }
    return
}
```

## License

MIT

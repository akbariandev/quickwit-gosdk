package main

import (
	"fmt"

	quickwitgosdk "github.com/akbariandev/quickwit-gosdk"
)

func main() {
	client := quickwitgosdk.NewClient("http://localhost:7280")

	resp, err := client.Search("stackoverflow", quickwitgosdk.SearchRequest{
		Query: "body:'*Microsoft SQL Server Profiler*'",
		// DefaultOperator: "AND",
		MaxHits: 10,
		// SearchFields:    []string{"body"},
		// SortByField:     &quickwitgosdk.SortByField{FieldName: "timestamp", Order: "desc"},
		// SnippetFields: &quickwitgosdk.SnippetRequest{FieldName: "message", MaxNumCharsPerFragment: 200},
		// TagFilters:    []string{"tag1:value1", "tag2:value2"},
		// Filter: "status_code >= 400",
	})

	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	// fmt.Printf("%v\n", resp)

	fmt.Printf("Found %d hits\n", resp.NumHits)
	for _, hit := range resp.Hits {
		fmt.Printf("  %v\n", hit.Fields)
	}

	// With API key authentication
	// client := quickwitgosdk.NewClient("http://localhost:7280",
	// 	quickwitgosdk.WithAPIKey("your-api-key"),
	// )

	// // With custom timeout
	// client := quickwitgosdk.NewClient("http://localhost:7280",
	// 	quickwitgosdk.WithTimeout(30*time.Second),
	// )

	// // With custom transport (e.g. for custom TLS)
	// client := quickwitgosdk.NewClient("http://localhost:7280",
	// 	quickwitgosdk.WithTransport(customTransport),
	// )
}

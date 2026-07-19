package quickwitgosdk

// SortByField defines a sort by a numeric field.
type SortByField struct {
	FieldName     string `json:"field_name"`
	Order         string `json:"order,omitempty"`          // "asc" or "desc"
	MissingValues string `json:"missing_values,omitempty"` // "first", "last", or a literal value
}

// SortByFieldDocOrder defines a sort by doc order.
type SortByFieldDocOrder struct {
	FieldName string `json:"field_name"`
	Order     string `json:"order,omitempty"` // "asc" or "desc"
}

// SnippetRequest defines snippet extraction for a field.
type SnippetRequest struct {
	FieldName              string `json:"field_name"`
	MaxNumCharsPerFragment int    `json:"max_num_chars_per_fragment,omitempty"`
}

// HighlightRequest defines highlight extraction for a field.
type HighlightRequest struct {
	FieldName     string   `json:"field_name"`
	Fragmenter    string   `json:"fragmenter,omitempty"`     // "plain" or "sentence"
	MaxNumChars   int      `json:"max_num_chars,omitempty"`
	NumFragments  int      `json:"num_fragments,omitempty"`
	PreTags       []string `json:"pre_tags,omitempty"`
	PostTags      []string `json:"post_tags,omitempty"`
}

// Hit represents a single search result hit.
type Hit struct {
	JSON       interface{} `json:"json"`
	Fragment   string      `json:"fragment,omitempty"`
	Score      *float64    `json:"score,omitempty"`
	PartialHit *PartialHit `json:"partial_hit,omitempty"`
}

// PartialHit represents a partial hit with a sorting value.
type PartialHit struct {
	DocID    string    `json:"doc_id"`
	Segment  string    `json:"segment_id,omitempty"`
	Shard    *Shard    `json:"shard,omitempty"`
	Sorting  []float64 `json:"sorting,omitempty"`
}

// Shard identifies a split and node.
type Shard struct {
	LeaderID   string `json:"leader_id"`
	ShardID    string `json:"shard_id"`
	Source     string `json:"source,omitempty"`
	Offset     uint64 `json:"offset,omitempty"`
	Position   uint32 `json:"position,omitempty"`
}

// SearchError represents an error returned within a search response.
type SearchError struct {
	Message string `json:"message"`
	Kind    string `json:"kind,omitempty"`
}

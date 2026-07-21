package quickwitgosdk

import "encoding/json"

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
// Quickwit returns hit fields flat in the JSON object alongside optional
// fragment, score, and partial_hit fields. Use the Fields map to access
// document fields.
type Hit struct {
	// Fields contains all document fields of the hit.
	// Use type assertions to access individual fields, e.g.:
	//   hit.Fields["body"].(string)
	//   hit.Fields["tags"].([]interface{})
	Fields     map[string]interface{} `json:"-"`
	Fragment   string                 `json:"fragment,omitempty"`
	Score      *float64               `json:"score,omitempty"`
	PartialHit *PartialHit            `json:"partial_hit,omitempty"`
}

// UnmarshalJSON implements custom unmarshaling for Hit.
// Quickwit returns document fields flat alongside fragment/score/partial_hit.
// This method separates known metadata fields from document fields.
func (h *Hit) UnmarshalJSON(data []byte) error {
	// First, unmarshal the full hit using a raw intermediate to preserve all fields.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract known metadata fields.
	if v, ok := raw["fragment"]; ok {
		if err := json.Unmarshal(v, &h.Fragment); err != nil {
			return err
		}
	}
	if v, ok := raw["score"]; ok {
		if err := json.Unmarshal(v, &h.Score); err != nil {
			return err
		}
	}
	if v, ok := raw["partial_hit"]; ok {
		if err := json.Unmarshal(v, &h.PartialHit); err != nil {
			return err
		}
	}

	// Everything else is a document field.
	fields := make(map[string]interface{})
	for k, v := range raw {
		switch k {
		case "fragment", "score", "partial_hit":
			// already handled above
		default:
			var val interface{}
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			fields[k] = val
		}
	}
	h.Fields = fields
	return nil
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

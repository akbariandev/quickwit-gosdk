package quickwitgosdk

import "encoding/json"

// DocMapping defines how documents are mapped to index fields.
type DocMapping struct {
	DocMappingUID      string          `json:"doc_mapping_uid,omitempty"`
	Mode               string          `json:"mode,omitempty"`               // "strict", "lenient", or "dynamic"
	DynamicMapping     *DynamicMapping `json:"dynamic_mapping,omitempty"`
	FieldMappings      []FieldMapping  `json:"field_mappings"`
	TimestampField     string          `json:"timestamp_field,omitempty"`
	TagFields          []string        `json:"tag_fields,omitempty"`
	MaxNumPartitions   int             `json:"max_num_partitions,omitempty"`
	IndexFieldPresence bool            `json:"index_field_presence,omitempty"`
	StoreDocumentSize  bool            `json:"store_document_size,omitempty"`
	StoreSource        bool            `json:"store_source,omitempty"`
	PartitionKey       string          `json:"partition_key,omitempty"`
	Tokenizers         []Tokenizer     `json:"tokenizers,omitempty"`
}

// DynamicMapping defines the mapping applied to fields not explicitly listed in field_mappings.
type DynamicMapping struct {
	Description string       `json:"description,omitempty"`
	ExpandDots  bool         `json:"expand_dots,omitempty"`
	Fast        *FastField   `json:"fast,omitempty"`
	Indexed     bool         `json:"indexed,omitempty"`
	Record      string       `json:"record,omitempty"` // "basic", "freq", or "position"
	Stored      bool         `json:"stored,omitempty"`
	Tokenizer   string       `json:"tokenizer,omitempty"`
}

// FieldMapping defines a single field in a doc mapping.
// Common fields across all field types; type-specific options are captured in Options.
type FieldMapping struct {
	Name          string                 `json:"name"`
	Type          string                 `json:"type"` // "text", "u64", "i64", "f64", "datetime", "json", "bool", "ip", "bytes"
	Fast          *FastField             `json:"fast,omitempty"`
	Fieldnorms    *bool                  `json:"fieldnorms,omitempty"`
	Indexed       bool                   `json:"indexed,omitempty"`
	Record        string                 `json:"record,omitempty"` // "basic", "freq", "position"
	Stored        bool                   `json:"stored,omitempty"`
	Tokenizer     string                 `json:"tokenizer,omitempty"`
	InputFormats  []string               `json:"input_formats,omitempty"`  // datetime only
	OutputFormat  string                 `json:"output_format,omitempty"`  // datetime only
	FastPrecision string                 `json:"fast_precision,omitempty"` // datetime only
	Coerce        bool                   `json:"coerce,omitempty"`         // numeric only
	ExpandDots    bool                   `json:"expand_dots,omitempty"`    // json only
	Options       map[string]interface{} `json:"-"`                        // catch-all for type-specific fields not modeled above
}

// FastField configures fast field behavior.
// It can be a bool (disabled/enabled) or an object with a normalizer.
type FastField struct {
	Enabled     bool   `json:"-"`
	Normalizer  string `json:"normalizer,omitempty"`
	isObj       bool
}

// MarshalJSON implements custom marshaling for FastField.
// Serializes as a bare false, bare true, or {"normalizer": "..."} object.
func (f FastField) MarshalJSON() ([]byte, error) {
	if !f.isObj {
		if f.Enabled {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	}
	// Object form with normalizer
	type alias FastField
	return json.Marshal(struct {
		Normalizer string `json:"normalizer,omitempty"`
	}{
		Normalizer: f.Normalizer,
	})
}

// UnmarshalJSON implements custom unmarshaling for FastField.
// Accepts a bare bool or an object with a normalizer field.
func (f *FastField) UnmarshalJSON(data []byte) error {
	// Try bool first.
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		f.Enabled = b
		f.isObj = false
		return nil
	}

	// Object form: {"normalizer": "raw"} or {"normalizer": {"name": "lowercase"}}
	var obj struct {
		Normalizer string `json:"normalizer"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	f.Normalizer = obj.Normalizer
	f.Enabled = true
	f.isObj = true
	return nil
}

// Tokenizer defines a custom tokenizer configuration.
type Tokenizer struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"` // "ngram", "regex", "simple", "source_code"
	Filters     []string `json:"filters,omitempty"`
	MinGram     int      `json:"min_gram,omitempty"`     // ngram only
	MaxGram     int      `json:"max_gram,omitempty"`     // ngram only
	PrefixOnly  bool     `json:"prefix_only,omitempty"`  // ngram only
	Pattern     string   `json:"pattern,omitempty"`      // regex only
}

// IndexingSettings configures how documents are indexed.
type IndexingSettings struct {
	CommitTimeoutSecs        int               `json:"commit_timeout_secs,omitempty"`
	DocstoreCompressionLevel int               `json:"docstore_compression_level,omitempty"`
	DocstoreBlocksize        int               `json:"docstore_blocksize,omitempty"`
	SplitNumDocsTarget       int               `json:"split_num_docs_target,omitempty"`
	MergePolicy              *MergePolicy      `json:"merge_policy,omitempty"`
	Resources                *IndexingResources `json:"resources,omitempty"`
}

// MergePolicy defines the merge policy for index splits.
type MergePolicy struct {
	Type             string `json:"type"` // "stable_log", "limit_interval", or "no_merge"
	MinLevelNumDocs  int    `json:"min_level_num_docs,omitempty"`
	MergeFactor      int    `json:"merge_factor,omitempty"`
	MaxMergeFactor   int    `json:"max_merge_factor,omitempty"`
	MaturationPeriod string `json:"maturation_period,omitempty"`
}

// IndexingResources defines resource limits for indexing.
type IndexingResources struct {
	// HeapSize accepts both numeric bytes (e.g. 2000000000) and human-readable strings (e.g. "2 GB").
	HeapSize json.Number `json:"heap_size,omitempty"`
}

// IngestSettings configures ingestion behavior.
type IngestSettings struct {
	MinShards    int  `json:"min_shards,omitempty"`
	ValidateDocs bool `json:"validate_docs,omitempty"`
}

// SearchSettings configures search behavior.
type SearchSettings struct {
	DefaultSearchFields []string `json:"default_search_fields,omitempty"`
}

// Retention defines the retention policy for an index.
type Retention struct {
	Period   string `json:"period,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

// Source represents an ingest source for an index.
type Source struct {
	Version     string `json:"version"`
	SourceID    string `json:"source_id"`
	NumPipelines int   `json:"num_pipelines"`
	Enabled     bool   `json:"enabled"`
	SourceType  string `json:"source_type"`
	InputFormat string `json:"input_format,omitempty"`
}

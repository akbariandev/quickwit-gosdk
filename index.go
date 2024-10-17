package quickwitgosdk

type MergePolicyType string

const (
	NoMerge    MergePolicyType = "no_merge"
	LimitMerge MergePolicyType = "limit_merge"
	StableLog  MergePolicyType = "stable_log"
)

type MergePolicy struct {
	Type           MergePolicyType `json:"type"`
	MaxMergeOps    int             `json:"max_merge_ops"`
	MergeFactor    int             `json:"merge_factor"`
	MaxMergeFactor int             `json:"max_merge_factor"`
}

type Resource struct {
	HeapSize string `json:"heap_size"`
}

type IndexingSettings struct {
	MergePolicy              MergePolicy `json:"merge_policy"`
	Resources                Resource    `json:"resources,omitempty"`
	CommitTimeoutSecs        int         `json:"commit_timeout_secs,omitempty"`
	DocstoreBlocksize        int         `json:"docstore_blocksize,omitempty"`
	DocstoreCompressionLevel int         `json:"docstore_compression_level,omitempty"`
	SplitNumDocsTarget       int         `json:"split_num_docs_target,omitempty"`
}

type Retention struct {
	Period   string `json:"period"`
	Schedule string `json:"schedule"`
}

type FieldMapping struct {
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Fast          bool     `json:"fast,omitempty"`
	InputFormats  []string `json:"input_formats,omitempty"`
	FastPrecision string   `json:"fast_precision,omitempty"`
	Record        string   `json:"record,omitempty"`
}

type DynamicMapping struct {
	Description string `json:"description,omitempty"`
	ExpandDots  bool   `json:"expand_dots,omitempty"`
	Fast        string `json:"fast,omitempty"`
	Indexed     bool   `json:"indexed,omitempty"`
	Record      string `json:"record,omitempty"`
	Stored      bool   `json:"stored,omitempty"`
	Tokenizer   string `json:"tokenizer,omitempty"`
}

type DocMapping struct {
	FieldMappings      []FieldMapping  `json:"field_mappings,omitempty"`
	PartitionKey       string          `json:"partition_key,omitempty"`
	MaxNumPartitions   int             `json:"max_num_partitions,omitempty"`
	TagFields          []string        `json:"tag_fields,omitempty"`
	TimestampField     string          `json:"timestamp_field,omitempty"`
	DynamicMapping     *DynamicMapping `json:"dynamic_mapping,omitempty"`
	IndexFieldPresence bool            `json:"index_field_presence,omitempty"`
	Mode               string          `json:"mode,omitempty"`
	StoreSource        bool            `json:"store_source,omitempty"`
	Tokenizers         []Tokenizers    `json:"tokenizers"`
}

type Transform struct {
	Script   string `json:"script"`
	Timezone string `json:"timezone"`
}

type Params struct {
	Filepath string `json:"filepath"`
}

type Sources struct {
	DesiredNumPipelines       int       `json:"desired_num_pipelines,omitempty"`
	Enabled                   bool      `json:"enabled"`
	InputFormat               string    `json:"input_format"`
	MaxNumPipelinesPerIndexer int       `json:"max_num_pipelines_per_indexer,omitempty"`
	SourceId                  string    `json:"source_id"`
	Transform                 Transform `json:"transform"`
	Version                   string    `json:"version"`
	Params                    Params    `json:"params"`
	SourceType                string    `json:"source_type"`
	NumPipelines              int       `json:"num_pipelines,omitempty"`
}

type Tokenizers struct {
	Type       string   `json:"type"`
	Filters    []string `json:"filters"`
	Name       string   `json:"name"`
	MaxGram    int      `json:"max_gram,omitempty"`
	MinGram    int      `json:"min_gram,omitempty"`
	PrefixOnly bool     `json:"prefix_only,omitempty"`
	Pattern    string   `json:"pattern,omitempty"`
}

type IndexConfig struct {
	DocMapping       DocMapping       `json:"doc_mapping"`
	IndexId          string           `json:"index_id"`
	IndexUri         string           `json:"index_uri"`
	IndexingSettings IndexingSettings `json:"indexing_settings"`
	Retention        Retention        `json:"retention"`
	SearchSettings   SearchSettings   `json:"search_settings"`
	Version          string           `json:"version"`
}

type CreateIndexRequest struct {
	Version          string           `json:"version"`
	IndexId          string           `json:"index_id"`
	SearchSettings   SearchSettings   `json:"search_settings,omitempty"`
	IndexingSettings IndexingSettings `json:"indexing_settings,omitempty"`
	Retention        Retention        `json:"retention,omitempty"`
	DocMapping       DocMapping       `json:"doc_mapping"`
}

type CreateIndexResponse struct {
	Checkpoint      struct{}    `json:"checkpoint"`
	CreateTimestamp int         `json:"create_timestamp"`
	IndexConfig     IndexConfig `json:"index_config"`
	IndexUid        string      `json:"index_uid"`
	Sources         []Sources   `json:"sources"`
	Version         string      `json:"version"`
}

func (c *Client) CreateIndex(request CreateIndexRequest) (response CreateIndexResponse, err error) {
	_, err = c.client.R().
		SetBody(request).
		SetHeader("Accept", "application/json").
		SetResult(&response).
		SetError(&ErrorMessage{}).
		Post("/api/v1/indexes")

	return response, err
}

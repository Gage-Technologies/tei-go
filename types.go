package tei

// EmbedRequest
// Mock of the embed request for the embedding server
type EmbedRequest struct {
	Inputs   string `json:"inputs"`
	Truncate bool   `json:"truncate"`
}

// EmbedResponse
// Mock of the embed response for the embedding server
type EmbedResponse [][]float32

// InfoResponse
// Mock of the info response for the embedding server
type InfoResponse struct {
	DockerLabel           *string `json:"docker_label"`
	MaxBatchRequests      *int    `json:"max_batch_requests"`
	MaxBatchTokens        int     `json:"max_batch_tokens"`
	MaxClientBatchSize    int     `json:"max_client_batch_size"`
	MaxConcurrentRequests int     `json:"max_concurrent_requests"`
	MaxInputLength        int     `json:"max_input_length"`
	ModelDtype            string  `json:"model_dtype"`
	ModelID               string  `json:"model_id"`
	ModelSha              string  `json:"model_sha"`
	Sha                   *string `json:"sha"`
	TokenizationWorkers   int     `json:"tokenization_workers"`
	Version               string  `json:"version"`
}

type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeTokenizer  ErrorType = "tokenizer"
	ErrorTypeBackend    ErrorType = "backend"
	ErrorTypeOverloaded ErrorType = "overloaded"
	ErrorTypeUnhealthy  ErrorType = "unhealthy"
)

// EmbedError
// Mock of the embed error for the embedding server
type EmbedError struct {
	Error     string    `json:"error"`
	ErrorType ErrorType `json:"error_type"`
}

// PredictRequest
// Represents a request to the predict endpoint
type PredictRequest struct {
	Inputs string `json:"inputs"`
}

// Prediction
// Represents a single prediction response with a score and label
type Prediction struct {
	Label string  `json:"label"`
	Score float32 `json:"score"`
}

// RerankRequest
// Represents a request to the rerank endpoint with a query and list of texts
type RerankRequest struct {
	Query string   `json:"query"`
	Texts []string `json:"texts"`
}

// Rank
// Represents a ranked response for reranking functionality
type Rank struct {
	Index int     `json:"index"`
	Score float32 `json:"score"`
	Text  string  `json:"text,omitempty"` // optional field
}

// TokenizeRequest
// Represents a request to the tokenize endpoint
type TokenizeRequest struct {
	Inputs           []string `json:"inputs"` // Can be a single string or an array of strings
	AddSpecialTokens bool     `json:"add_special_tokens"`
}

// SimpleToken
// Represents a single tokenized output
type SimpleToken struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Special bool   `json:"special"`
	Start   int    `json:"start"`
	Stop    int    `json:"stop"`
}

// DecodeRequest
// Represents a request to the decode endpoint
type DecodeRequest struct {
	Ids               [][]int `json:"ids"` // Can be a single array of integers or a nested array
	SkipSpecialTokens bool    `json:"skip_special_tokens"`
}

// DecodeResponse
// Represents a response from the decode endpoint, which is an array of strings
type DecodeResponse []string

// EmbedSparseRequest
// Represents a request to the embed_sparse endpoint
type EmbedSparseRequest struct {
	Inputs   []string `json:"inputs"` // Can be a single string or an array of strings
	Truncate bool     `json:"truncate"`
}

// SparseValue
// Represents a single sparse embedding output
type SparseValue struct {
	Index int     `json:"index"`
	Value float32 `json:"value"`
}

// EmbedAllRequest
// Represents a request to the embed_all endpoint
type EmbedAllRequest struct {
	Inputs   []string `json:"inputs"` // Can be a single string or an array of strings
	Truncate bool     `json:"truncate"`
}

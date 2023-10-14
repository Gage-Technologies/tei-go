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
	ErrorTypeValidation = "validation"
	ErrorTypeTokenizer  = "tokenizer"
	ErrorTypeBackend    = "backend"
	ErrorTypeOverloaded = "overloaded"
)

// EmbedError
// Mock of the embed error for the embedding server
type EmbedError struct {
	Error     string    `json:"error"`
	ErrorType ErrorType `json:"error_type"`
}

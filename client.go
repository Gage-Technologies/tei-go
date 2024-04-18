package tei

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client
// Wrapper client around the http api of the embedding server. Example usage:
//
//	client := embedding_server.NewClient(http://localhost:8080, nil, nil, time.Second * 30)
//	res, err := client.Embed("Hi there!")
//	if err != nil {
//	    panic(err)
//	}
//	fmt.Println("Embedding: ", res[0])
type Client struct {
	baseURL string
	client  *http.Client
	headers map[string]string
	cookies map[string]string
}

// NewClient
// Create a new client for the embedding server
func NewClient(baseURL string, headers map[string]string, cookies map[string]string, timeout time.Duration) *Client {
	// create custom http client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// trim trailing slashes from url if any
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		baseURL: baseURL,
		client:  client,
		headers: headers,
		cookies: cookies,
	}
}

// prepareRequest
// Prepare a request for the embedding server
func (c *Client) prepareRequest(req *http.Request) *http.Request {
	// set content type
	req.Header.Set("Content-Type", "application/json")

	// set the headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// set the cookies
	for k, v := range c.cookies {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	return req
}

// Info
// Get information about the embedding server
func (c *Client) Info() (*InfoResponse, error) {
	// create the http request
	httpReq, err := http.NewRequest(
		"GET",
		c.baseURL+"/info",
		nil,
	)

	// prepare the request
	httpReq = c.prepareRequest(httpReq)

	// execute the request
	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}

	// parse the response
	var response InfoResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// HealthCheck
// Check the health of the embedding server
func (c *Client) HealthCheck() (bool, error) {
	httpReq, err := http.NewRequest("GET", c.baseURL+"/health", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq = c.prepareRequest(httpReq)
	res, err := c.client.Do(httpReq)
	if err != nil {
		return false, fmt.Errorf("failed to execute http request: %w", err)
	}

	return res.StatusCode == 200, nil
}

// Metrics
// Retrieve Prometheus metrics from the server
func (c *Client) Metrics() (string, error) {
	httpReq, err := http.NewRequest("GET", c.baseURL+"/metrics", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq = c.prepareRequest(httpReq)
	res, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to execute http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", c.handleErrorResponse(res)
	}

	metricsData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read metrics data: %w", err)
	}

	return string(metricsData), nil
}

// Embed
// Embed a string using the embedding server
func (c *Client) Embed(inputs string, truncate bool) (EmbedResponse, error) {
	// ensure inputs is not empty
	if inputs == "" {
		return nil, ErrEmptyInputs
	}

	// create the request
	req, err := json.Marshal(EmbedRequest{
		Inputs:   inputs,
		Truncate: truncate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// create the http request
	httpReq, err := http.NewRequest(
		"POST",
		c.baseURL+"/embed",
		bytes.NewBuffer(req),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// prepare the request
	httpReq = c.prepareRequest(httpReq)

	// execute the request
	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer res.Body.Close()

	// handle non-200 status codes
	if res.StatusCode != 200 {
		return nil, c.handleErrorResponse(res)
	}

	// parse the response
	var response EmbedResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// EmbedAll
// Retrieve embeddings for all input tokens without pooling
func (c *Client) EmbedAll(inputs []string, truncate bool) ([][][]float32, error) {
	req, err := json.Marshal(EmbedAllRequest{
		Inputs:   inputs,
		Truncate: truncate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/embed_all", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq = c.prepareRequest(httpReq)
	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, c.handleErrorResponse(res)
	}

	var response [][][]float32
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// EmbedSparse
// Retrieve sparse embeddings using SPLADE pooling if supported by the model
func (c *Client) EmbedSparse(inputs []string, truncate bool) ([][]SparseValue, error) {
	req, err := json.Marshal(EmbedSparseRequest{
		Inputs:   inputs,
		Truncate: truncate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/embed_sparse", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq = c.prepareRequest(httpReq)
	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, c.handleErrorResponse(res)
	}

	var response [][]SparseValue
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// Predict
// Make predictions using the provided inputs
func (c *Client) Predict(inputs string) ([]Prediction, error) {
	req, err := json.Marshal(PredictRequest{
		Inputs: inputs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/predict", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq = c.prepareRequest(httpReq)
	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, c.handleErrorResponse(res)
	}

	var response []Prediction
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// Rerank
// Rerank a list of texts based on a query
func (c *Client) Rerank(query string, texts []string) ([]Rank, error) {
	req, err := json.Marshal(RerankRequest{
		Query: query,
		Texts: texts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/rerank", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq = c.prepareRequest(httpReq)
	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, c.handleErrorResponse(res)
	}

	var response []Rank
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	for i, rank := range response {
		response[i].Text = texts[rank.Index]
	}

	return response, nil
}

// Tokenize
// Tokenize input text
func (c *Client) Tokenize(inputs []string, addSpecialTokens bool) ([][]SimpleToken, error) {
	req, err := json.Marshal(TokenizeRequest{
		Inputs:           inputs,
		AddSpecialTokens: addSpecialTokens,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/tokenize", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq = c.prepareRequest(httpReq)
	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, c.handleErrorResponse(res)
	}

	var response [][]SimpleToken
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// Decode
// Decode input ids into text
func (c *Client) Decode(ids [][]int, skipSpecialTokens bool) ([]string, error) {
	req, err := json.Marshal(DecodeRequest{
		Ids:               ids,
		SkipSpecialTokens: skipSpecialTokens,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/decode", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq = c.prepareRequest(httpReq)
	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, c.handleErrorResponse(res)
	}

	var response []string
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// handleErrorResponse
// Handle API error responses uniformly
func (c *Client) handleErrorResponse(res *http.Response) error {
	var response EmbedError
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("error decoding error response: %w", err)
	}

	switch response.ErrorType {
	case ErrorTypeValidation:
		return fmt.Errorf("%w: %v", ErrValidation, response.Error)
	case ErrorTypeTokenizer:
		return fmt.Errorf("%w: %v", ErrTokenizer, response.Error)
	case ErrorTypeOverloaded:
		return ErrOverloaded
	case ErrorTypeBackend:
		return fmt.Errorf("%w: %v", ErrBackend, response.Error)
	default:
		return fmt.Errorf("unhandled error type %v: %v", response.ErrorType, response.Error)
	}
}

package tei

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		// parse the error body
		var response EmbedError
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", res.StatusCode, response.Error)
		}

		// select a known error to wrap if we have a valid type
		switch response.ErrorType {
		case ErrorTypeValidation:
			return nil, fmt.Errorf("%w: %v", ErrValidation, response.Error)
		case ErrorTypeTokenizer:
			return nil, fmt.Errorf("%w: %v", ErrTokenizer, response.Error)
		case ErrorTypeOverloaded:
			return nil, ErrOverloaded
		case ErrorTypeBackend:
			return nil, fmt.Errorf("%w: %v", ErrBackend, response.Error)
		default:
			return nil, fmt.Errorf("embed request failed with unkown status code %v: %v", res.StatusCode, response.Error)
		}

	}

	// parse the response
	var response EmbedResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

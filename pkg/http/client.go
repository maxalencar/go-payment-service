package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ResilientHTTPClient wraps http.Client and includes Circuit Breaker and Exponential Backoff
type ResilientHTTPClient struct {
	client  *http.Client
	breaker *gobreaker.CircuitBreaker
}

// NewResilientHTTPClient initializes the HTTP client with circuit breaker and timeout settings
func NewResilientHTTPClient() *ResilientHTTPClient {
	// Configure the circuit breaker
	cbSettings := gobreaker.Settings{
		Name:        "HTTP Client Circuit Breaker",
		MaxRequests: 1,                // Allow 1 request in Half-Open state
		Interval:    60 * time.Second, // Reset failure count every 60 seconds
		Timeout:     5 * time.Second,  // Circuit stays open for 5 seconds before testing
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3 // Trip the breaker after 3 failures
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			slog.Info("Circuit breaker state changed", slog.Any("from", from), slog.Any("to", to))
		},
	}

	return &ResilientHTTPClient{
		client:  &http.Client{Timeout: 10 * time.Second}, // Set an overall request timeout
		breaker: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

// Do makes an HTTP request, applies exponential backoff retries, and integrates the circuit breaker
func (hc *ResilientHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response

	// Retry logic wrapped with circuit breaker
	operation := func() error {
		// Execute the HTTP request within the circuit breaker context
		result, err := hc.breaker.Execute(func() (interface{}, error) {
			var err error

			resp, err = hc.client.Do(req)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode >= 400 && resp.StatusCode <= 499 { // Treat 4xx HTTP responses as failures
				return nil, fmt.Errorf("received client error: %d", resp.StatusCode)
			}

			if resp.StatusCode >= 500 { // Treat 5xx HTTP responses as failures
				return nil, fmt.Errorf("received server error: %d", resp.StatusCode)
			}

			return resp, nil
		})
		if err != nil {
			return err
		}

		resp = result.(*http.Response)

		return nil
	}

	// Use exponential backoff for retrying the request
	backOff := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)

	// Retry the operation with the backoff strategy
	err := backoff.Retry(operation, backOff)
	if err != nil {
		return nil, fmt.Errorf("http request failed after retries: %w", err)
	}

	return resp, nil
}

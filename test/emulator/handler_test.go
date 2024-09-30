package emulator

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	paymenthttp "go-payment-service/pkg/http"
	"go-payment-service/pkg/model"
)

type TestSuite struct {
	suite.Suite
	handler               *handler
	callbackServer        *httptest.Server
	callbackHandlerCalled bool
}

// SetupSuite runs before all tests
func (suite *TestSuite) SetupSuite() {
	client := paymenthttp.NewResilientHTTPClient()
	repository := newMemoryRepository()
	service := newService(client, repository)
	suite.handler = newHandler(service)
	suite.callbackServer = suite.newCallbackHTTPTestServer()
}

func (suite *TestSuite) TestProcessHandler() {
	testCases := []struct {
		name          string
		given         ProcessRequest
		givenMIMEType string
		expected      ProcessResponse
		expectedCode  int
	}{
		{
			name: "success json",
			given: ProcessRequest{
				OrderID: "order-123",
				Amount: model.Money{
					Amount:   1000,
					Currency: "USD",
				},
				CardDetails: model.CardDetails{
					Number:      "4111111111111111",
					Name:        "John Doe",
					ExpiryMonth: 12,
					ExpiryYear:  2023,
					CVV:         "123",
				},
				Type: model.Deposit,
			},
			givenMIMEType: paymenthttp.MIMETypeJSON,
			expected: ProcessResponse{
				Status: model.Succeeded,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "success json with callback",
			given: ProcessRequest{
				OrderID: "order-123",
				Amount: model.Money{
					Amount:   1000,
					Currency: "USD",
				},
				CardDetails: model.CardDetails{
					Number:      "4111111111111111",
					Name:        "John Doe",
					ExpiryMonth: 12,
					ExpiryYear:  2023,
					CVV:         "123",
				},
				CallbackURL: suite.callbackServer.URL + "/callback",
				Type:        model.Deposit,
			},
			givenMIMEType: paymenthttp.MIMETypeJSON,
			expected: ProcessResponse{
				Status: model.Succeeded,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "success xml",
			given: ProcessRequest{
				OrderID: "order-1234",
				Amount: model.Money{
					Amount:   1005,
					Currency: "USD",
				},
				CardDetails: model.CardDetails{
					Number:      "4111111111111111",
					Name:        "John Doe",
					ExpiryMonth: 12,
					ExpiryYear:  2023,
					CVV:         "123",
				},
				Type: model.Deposit,
			},
			givenMIMEType: paymenthttp.MIMETypeXML,
			expected: ProcessResponse{
				Status: model.Succeeded,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := paymenthttp.Marshal(tc.givenMIMEType, tc.given)
			suite.Require().NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/process", bytes.NewReader(b))
			r.Header.Add(paymenthttp.HeaderContentType, tc.givenMIMEType)

			suite.handler.mux.ServeHTTP(w, r)

			suite.Equal(tc.expectedCode, w.Code)

			var resp ProcessResponse
			err = paymenthttp.Decode(w.Body, tc.givenMIMEType, &resp)
			suite.Require().NoError(err)

			suite.NotEmpty(resp.TransactionID)
			suite.NotEmpty(resp.ProcessedAt)
			suite.NotEmpty(resp.Data)

			suite.Equal(tc.expected.Status, resp.Status)
			suite.Equal(tc.given, resp.Data)

			if tc.given.CallbackURL != "" {
				waitFor, err := time.ParseDuration("5s")
				suite.Require().NoError(err)

				tick, err := time.ParseDuration("100ms")
				suite.Require().NoError(err)

				// Wait for callback
				suite.Eventually(func() bool {
					return suite.callbackHandlerCalled
				}, waitFor, tick)
			}
		})
	}
}

func (suite *TestSuite) newCallbackHTTPTestServer() *httptest.Server {
	// Initialize HTTP request multiplexer
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("POST /callback", func(w http.ResponseWriter, r *http.Request) {
		suite.callbackHandlerCalled = true

		w.WriteHeader(http.StatusOK)
	})

	return httptest.NewServer(mux)

}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

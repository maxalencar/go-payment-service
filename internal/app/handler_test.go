package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	paymenthttp "go-payment-service/pkg/http"
	"go-payment-service/pkg/model"
	"go-payment-service/test/emulator"
)

type TestHandlerSuite struct {
	suite.Suite
	handler *handler
}

// SetupSuite runs before all tests
func (suite *TestHandlerSuite) SetupSuite() {
	wg := &sync.WaitGroup{}

	// Initialize HTTP client
	resilientHTTPClient := paymenthttp.NewResilientHTTPClient()

	// Initialize payment gateways
	gatewayEmulator := emulator.Start()

	gateways := map[string]PaymentGateway{
		"gatewayA": newGatewayAAdapter(resilientHTTPClient, gatewayEmulator.URL),
		"gatewayB": newGatewayBAdapter(resilientHTTPClient, gatewayEmulator.URL),
	}

	repository := newMemoryTransactionRepository()
	service := newTransactionService(wg, gateways, repository)
	suite.handler = newHandler(service)
}

func (suite *TestHandlerSuite) TestDeposit() {
	testCases := []struct {
		name          string
		given         model.DepositRequest
		givenMIMEType string
		expected      model.Transaction
		expectedCode  int
	}{
		{
			name: "pending json",
			given: model.DepositRequest{
				BaseRequest: model.BaseRequest{
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
					GatewayDetails: model.GatewayDetails{
						ID:   "gatewayA",
						Name: "Gateway A",
					},
				},
			},
			givenMIMEType: paymenthttp.MIMETypeJSON,
			expected: model.Transaction{
				Status: model.Pending,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "pending xml",
			given: model.DepositRequest{
				BaseRequest: model.BaseRequest{
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
					GatewayDetails: model.GatewayDetails{
						ID:   "gatewayB",
						Name: "Gateway B",
					},
				},
			},
			givenMIMEType: paymenthttp.MIMETypeXML,
			expected: model.Transaction{
				Status: model.Pending,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := paymenthttp.Marshal(tc.givenMIMEType, tc.given)
			suite.Require().NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewReader(b))
			r.Header.Add(paymenthttp.HeaderContentType, tc.givenMIMEType)

			suite.handler.mux.ServeHTTP(w, r)

			suite.Equal(tc.expectedCode, w.Code)

			var resp model.DepositResponse
			err = paymenthttp.Decode(w.Body, tc.givenMIMEType, &resp)
			suite.Require().NoError(err)

			suite.NotEmpty(resp.TransactionID)
			suite.Empty(resp.ProcessedAt)

			suite.Equal(tc.expected.Status, resp.Status)
		})
	}
}

func (suite *TestHandlerSuite) TestWithdrawal() {
	testCases := []struct {
		name          string
		given         model.WithdrawalRequest
		givenMIMEType string
		expected      model.Transaction
		expectedCode  int
	}{
		{
			name: "pending json",
			given: model.WithdrawalRequest{
				BaseRequest: model.BaseRequest{
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
					GatewayDetails: model.GatewayDetails{
						ID:   "gatewayA",
						Name: "Gateway A",
					},
				},
			},
			givenMIMEType: paymenthttp.MIMETypeJSON,
			expected: model.Transaction{
				Status: model.Pending,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "pending xml",
			given: model.WithdrawalRequest{
				BaseRequest: model.BaseRequest{
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
					GatewayDetails: model.GatewayDetails{
						ID:   "gatewayB",
						Name: "Gateway B",
					},
				},
			},
			givenMIMEType: paymenthttp.MIMETypeXML,
			expected: model.Transaction{
				Status: model.Pending,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := paymenthttp.Marshal(tc.givenMIMEType, tc.given)
			suite.Require().NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/withdrawal", bytes.NewReader(b))
			r.Header.Add(paymenthttp.HeaderContentType, tc.givenMIMEType)

			suite.handler.mux.ServeHTTP(w, r)

			suite.Equal(tc.expectedCode, w.Code)

			var resp model.DepositResponse
			err = paymenthttp.Decode(w.Body, tc.givenMIMEType, &resp)
			suite.Require().NoError(err)

			suite.NotEmpty(resp.TransactionID)
			suite.Empty(resp.ProcessedAt)

			suite.Equal(tc.expected.Status, resp.Status)
		})
	}
}

func TestTestHandlerSuite(t *testing.T) {
	suite.Run(t, new(TestHandlerSuite))
}

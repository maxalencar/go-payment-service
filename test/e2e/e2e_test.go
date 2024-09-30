package e2e

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"go-payment-service/internal/app"
	paymenthttp "go-payment-service/pkg/http"
	"go-payment-service/pkg/model"
)

type TestE2ESuite struct {
	suite.Suite
	client paymenthttp.HTTPClient
	server *httptest.Server
}

// It runs before all tests
func (suite *TestE2ESuite) SetupSuite() {
	srv := app.NewServer()
	suite.server = srv.StartTest()
	suite.client = paymenthttp.NewResilientHTTPClient()
}

// It runs after all tests
func (suite *TestE2ESuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *TestE2ESuite) TestDeposit() {
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
		{
			name: "success json",
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
						ID:          "gatewayA",
						Name:        "Gateway A",
						CallbackURL: suite.server.URL + "/callback",
					},
				},
			},
			givenMIMEType: paymenthttp.MIMETypeJSON,
			expected: model.Transaction{
				Status: model.Succeeded,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "success xml",
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
						ID:          "gatewayB",
						Name:        "Gateway B",
						CallbackURL: suite.server.URL + "/callback",
					},
				},
			},
			givenMIMEType: paymenthttp.MIMETypeXML,
			expected: model.Transaction{
				Status: model.Succeeded,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			b, err := paymenthttp.Marshal(tc.givenMIMEType, tc.given)
			suite.NoError(err)

			req, err := http.NewRequest(http.MethodPost, suite.server.URL+"/deposit", bytes.NewBuffer(b))
			suite.NoError(err)
			req.Header.Add(paymenthttp.HeaderContentType, tc.givenMIMEType)

			resp, err := suite.client.Do(req)
			suite.NoError(err)
			defer resp.Body.Close()

			suite.Equal(tc.expectedCode, resp.StatusCode)

			var dr model.DepositResponse
			err = paymenthttp.Decode(resp.Body, tc.givenMIMEType, &dr)
			suite.NoError(err)

			suite.NotEmpty(dr.TransactionID)
			suite.Empty(dr.ProcessedAt)

			expectedStatus := tc.expected.Status

			// If the expected status is Succeeded, then we need to check the transaction status
			if expectedStatus == model.Succeeded {
				url := suite.server.URL + "/transactions/" + dr.TransactionID
				req, err := http.NewRequest(http.MethodGet, url, nil)
				suite.NoError(err)
				req.Header.Add(paymenthttp.HeaderContentType, tc.givenMIMEType)

				resp, err := suite.client.Do(req)
				suite.NoError(err)
				defer resp.Body.Close()

				var tr model.Transaction
				err = paymenthttp.Decode(resp.Body, tc.givenMIMEType, &tr)
				suite.NoError(err)

				suite.Equal(expectedStatus, tr.Status)

				jsonBytes, err := paymenthttp.Marshal(tc.givenMIMEType, tr)
				suite.NoError(err)
				slog.Info(string(jsonBytes))
				return
			}

			suite.Equal(expectedStatus, dr.Status)
		})
	}
}

func (suite *TestE2ESuite) TestWithdrawal() {
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
		{
			name: "success json",
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
						ID:          "gatewayA",
						Name:        "Gateway A",
						CallbackURL: suite.server.URL + "/callback",
					},
				},
			},
			givenMIMEType: paymenthttp.MIMETypeJSON,
			expected: model.Transaction{
				Status: model.Succeeded,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "success xml",
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
						ID:          "gatewayB",
						Name:        "Gateway B",
						CallbackURL: suite.server.URL + "/callback",
					},
				},
			},
			givenMIMEType: paymenthttp.MIMETypeXML,
			expected: model.Transaction{
				Status: model.Succeeded,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			b, err := paymenthttp.Marshal(tc.givenMIMEType, tc.given)
			suite.NoError(err)

			req, err := http.NewRequest(http.MethodPost, suite.server.URL+"/withdrawal", bytes.NewBuffer(b))
			suite.NoError(err)
			req.Header.Add(paymenthttp.HeaderContentType, tc.givenMIMEType)

			resp, err := suite.client.Do(req)
			suite.NoError(err)
			defer resp.Body.Close()

			suite.Equal(tc.expectedCode, resp.StatusCode)

			var dr model.WithdrawalResponse
			err = paymenthttp.Decode(resp.Body, tc.givenMIMEType, &dr)
			suite.NoError(err)

			suite.NotEmpty(dr.TransactionID)
			suite.Empty(dr.ProcessedAt)

			expectedStatus := tc.expected.Status

			// If the expected status is Succeeded, then we need to check the transaction status
			if expectedStatus == model.Succeeded {
				url := suite.server.URL + "/transactions/" + dr.TransactionID
				req, err := http.NewRequest(http.MethodGet, url, nil)
				suite.NoError(err)
				req.Header.Add(paymenthttp.HeaderContentType, tc.givenMIMEType)

				resp, err := suite.client.Do(req)
				suite.NoError(err)
				defer resp.Body.Close()

				var tr model.Transaction
				err = paymenthttp.Decode(resp.Body, tc.givenMIMEType, &tr)
				suite.NoError(err)

				suite.Equal(expectedStatus, tr.Status)
				return
			}

			suite.Equal(expectedStatus, dr.Status)
		})
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTestE2ESuite(t *testing.T) {
	suite.Run(t, new(TestE2ESuite))
}

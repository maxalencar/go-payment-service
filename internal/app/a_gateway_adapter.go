package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	paymenthttp "go-payment-service/pkg/http"
	"go-payment-service/pkg/model"
)

type GatewayA struct {
	client   paymenthttp.HTTPClient
	endpoint string
}

func newGatewayAAdapter(client paymenthttp.HTTPClient, endpoint string) *GatewayA {
	return &GatewayA{
		client:   client,
		endpoint: endpoint,
	}
}

func (g *GatewayA) ProcessTransaction(tx model.Transaction) (model.GatewayResponse, error) {
	jsonData, err := json.Marshal(g.buildGatewayRequest(tx))
	if err != nil {
		return model.GatewayResponse{}, fmt.Errorf("failed to marshal gateway request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, g.endpoint+"/process", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.GatewayResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(paymenthttp.HeaderContentType, paymenthttp.MIMETypeJSON)

	resp, err := g.client.Do(req)
	if err != nil {
		return model.GatewayResponse{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	var gr model.GatewayResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return model.GatewayResponse{}, fmt.Errorf("failed to decode gateway response: %w", err)
	}

	return gr, nil
}

func (g *GatewayA) buildGatewayRequest(tx model.Transaction) model.GatewayRequest {
	return model.GatewayRequest{
		OrderID:     tx.ID,
		Amount:      tx.Amount,
		CardDetails: tx.CardDetails,
		CallbackURL: tx.GatewayDetails.CallbackURL,
		Type:        tx.Type,
	}
}

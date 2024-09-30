package app

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"

	paymenthttp "go-payment-service/pkg/http"
	"go-payment-service/pkg/model"
)

type GatewayB struct {
	client   paymenthttp.HTTPClient
	endpoint string
}

func newGatewayB(client paymenthttp.HTTPClient, endpoint string) *GatewayB {
	return &GatewayB{
		client:   client,
		endpoint: endpoint,
	}
}

func (g *GatewayB) ProcessTransaction(tx model.Transaction) (model.GatewayResponse, error) {
	xmlData, err := xml.Marshal(g.buildGatewayRequest(tx))
	if err != nil {
		return model.GatewayResponse{}, fmt.Errorf("failed to marshal gateway request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, g.endpoint+"/process", bytes.NewBuffer(xmlData))
	if err != nil {
		return model.GatewayResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(paymenthttp.HeaderContentType, paymenthttp.MIMETypeXML)

	resp, err := g.client.Do(req)
	if err != nil {
		return model.GatewayResponse{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.GatewayResponse{}, fmt.Errorf("gateway returned non-200 status code: %d", resp.StatusCode)
	}

	var gr model.GatewayResponse
	if err := xml.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return model.GatewayResponse{}, fmt.Errorf("failed to decode gateway response: %w", err)
	}

	return gr, nil
}

func (g *GatewayB) buildGatewayRequest(tx model.Transaction) model.GatewayRequest {
	return model.GatewayRequest{
		OrderID:     tx.ID,
		Amount:      tx.Amount,
		CardDetails: tx.CardDetails,
		CallbackURL: tx.GatewayDetails.CallbackURL,
		Type:        tx.Type,
	}
}

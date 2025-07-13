package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DataInput struct {
	CorrelationId string
	Amount        float64
	RequestedAt   time.Time
}

type healthResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

func HealthCkeckDefaultRequest() (*healthResponse, error) {
	url := fmt.Sprintf("%s%s", defaultProcessorURL, "/payments/service-health")
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		return nil, err
	}

	var result healthResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Erro ao fazer unmarshal:", err)
		return nil, err
	}

	return &result, nil
}

func PaymentProcesorDefaultRequest(correlationId string, amount float64, requestedAt time.Time) error {
	url := fmt.Sprintf("%s%s", defaultProcessorURL, "/payments")
	data := DataInput{
		CorrelationId: correlationId,
		Amount:        amount,
		RequestedAt:   requestedAt,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error ao serializar json data", err)
		return err
	}

	response, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		fmt.Println("Error ao realizar requisição payments", err)
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("requisição retornou status %d", response.StatusCode)
	}

	return nil
}

// =================================================== Fallback =======================================

func HealthCkeckFallbackRequest() (*healthResponse, error) {
	url := fmt.Sprintf("%s%s", defaultProcessorURL, "/payments/service-health")
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		return nil, err
	}

	var result healthResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Erro ao fazer unmarshal:", err)
		return nil, err
	}

	return &result, nil
}

func PaymentProcesorFallbackRequest(correlationId string, amount float64, requestedAt time.Time) error {
	url := fmt.Sprintf("%s%s", fallbackProcessorURL, "/payments")

	data := DataInput{
		CorrelationId: correlationId,
		Amount:        amount,
		RequestedAt:   requestedAt,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error ao serializar json data", err)
		return err
	}

	response, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		fmt.Println("Error ao realizar requisição payments", err)
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("requisição retornou status %d", response.StatusCode)
	}

	return nil
}

package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"time"
)

type MessageParser struct {
	CorrelationId string
	Amount        float64
}

type Payment struct {
	CorrelationId string
	Amount        float64
	RequestedAt   time.Time
	Origin        *string
}

type SummaryDetail struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type PaymentSummary map[string]SummaryDetail

func unmarshall(reader *bufio.Reader) (*MessageParser, error) {
	var msg MessageParser

	correlationIdBytes := make([]byte, 36)
	_, err := io.ReadFull(reader, correlationIdBytes)
	if err != nil {
		return nil, fmt.Errorf("error ao ler o correlationId: %w", err)
	}

	msg.CorrelationId = string(correlationIdBytes)

	err = binary.Read(reader, binary.BigEndian, &msg.Amount)
	if err != nil {
		return nil, fmt.Errorf("error ao ler amount: %w", err)
	}

	return &msg, nil
}

func SaveMessage(correlationId string, amount float64, requestedAt time.Time, origin string) error {
	payment := Payment{
		CorrelationId: correlationId,
		Amount:        amount,
		RequestedAt:   requestedAt,
		Origin:        &origin,
	}

	_, err := dbInstance.Exec(`
		INSERT INTO payments (correlationid, amount, requested_at, origin)
		VALUES ($1, $2, $3, $4)
	`, payment.CorrelationId, payment.Amount, payment.RequestedAt, payment.Origin)

	return err
}

func GetPaymentsSummary(from string, to string) PaymentSummary {
	query := `SELECT amount, origin FROM payments WHERE requested_at BETWEEN $1 AND $2 `
	rows, err := dbInstance.Query(query, from, to)
	if err != nil {
		log.Printf("Erro ao executar a query: %v", err)
		return PaymentSummary{}
	}

	defer rows.Close()

	summary := make(PaymentSummary)

	summary[DEFAULT] = SummaryDetail{}
	summary[FALLBACK] = SummaryDetail{}

	for rows.Next() {
		var amount float64
		var origin string

		err := rows.Scan(&amount, &origin)
		if err != nil {
			log.Printf("Erro ao escanear linha: %v", err)
			continue
		}

		detail := summary[origin]
		detail.TotalRequests++
		detail.TotalAmount += amount
		summary[origin] = detail
	}

	for origin, detail := range summary {
		detail.TotalAmount = math.Round(detail.TotalAmount*100) / 100
		summary[origin] = detail
	}

	if err = rows.Err(); err != nil {
		log.Printf("Erro após iterar sobre as linhas: %v", err)
		return PaymentSummary{}
	}

	return summary
}

func Purge() error {
	_, err := dbInstance.Exec(`DELETE FROM payments`)
	return err
}

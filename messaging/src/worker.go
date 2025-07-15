package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var channelMsg = make(chan MessageParser, 50000)

var selectedOrigin atomic.Uint32

func RunHealthCkeck() {
	for {
		respFallback, err := HealthCkeckFallbackRequest()
		if err != nil {
			fmt.Println(err)
		}
		respDefault, err := HealthCkeckDefaultRequest()
		if err != nil {
			fmt.Println(err)
		}

		if respDefault != nil && respFallback != nil {
			if respDefault.MinResponseTime > respFallback.MinResponseTime {
				selectedOrigin.Store(2)
			} else {
				selectedOrigin.Store(1)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func GetCurrentOrigin() string {
	if selectedOrigin.Load() == 1 {
		return DEFAULT
	} else if selectedOrigin.Load() == 2 {
		return FALLBACK
	}
	return ""
}

func RunWorker() {
	for {
		message := <-channelMsg
		now := time.Now().UTC()
		switch GetCurrentOrigin() {
		case DEFAULT:
			err := PaymentProcesorDefaultRequest(message.CorrelationId, message.Amount, now)
			if err != nil {
				channelMsg <- message
				return
			} else {
				SaveMessage(message.CorrelationId, message.Amount, now, DEFAULT)
			}
		case FALLBACK:
			err := PaymentProcesorFallbackRequest(message.CorrelationId, message.Amount, now)
			if err != nil {
				channelMsg <- message
				return
			} else {
				SaveMessage(message.CorrelationId, message.Amount, now, FALLBACK)
			}
		default:
		}
	}
}

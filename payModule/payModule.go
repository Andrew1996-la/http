package payModule

import (
	"fmt"
	"net/http"
	"sync"
)

type Purchase struct {
	Object string  `json:"object"`
	Amount float64 `json:"amount"`
}

type PurchaseResponse struct {
	History []Purchase `json:"history"`
}

type PaymentModule struct {
	Money   float64
	History []Purchase
	mtx     sync.Mutex
}

func PayModule() {
	pay := &PaymentModule{
		Money:   1000,
		History: make([]Purchase, 0),
	}

	http.HandleFunc("/pay", pay.handlePay)
	http.HandleFunc("/history", pay.getHandleHistory)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err)
		return
	}
}

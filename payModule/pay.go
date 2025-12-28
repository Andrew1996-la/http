package payModule

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *PaymentModule) handlePay(w http.ResponseWriter, req *http.Request) {
	purchase := Purchase{}

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&purchase); err != nil {
		http.Error(w, "Не валидный JSON", http.StatusBadRequest)
		return
	}

	p.mtx.Lock()
	if p.Money-purchase.Amount < 0 {
		http.Error(w, "Недостаточно денег", http.StatusBadRequest)
		return
	}

	p.Money = p.Money - purchase.Amount
	p.History = append(p.History, purchase)
	w.Header().Set("Content-Type", "application/json")

	resp := ResponsePay{
		Message: fmt.Sprintf(
			"совершена покупа %s на сумму: %.2f",
			purchase.Object,
			purchase.Amount,
		),
		Balance: p.Money,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.mtx.Unlock()
}

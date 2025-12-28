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
	_, err := fmt.Fprintf(
		w,
		"совершена покупа %s на сумму: %.2f\nТекущий баланс: %.2f",
		purchase.Object,
		purchase.Amount,
		p.Money,
	)
	p.mtx.Unlock()

	if err != nil {
		fmt.Println(err)
		return
	}
}

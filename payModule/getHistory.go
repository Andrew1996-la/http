package payModule

import (
	"encoding/json"
	"net/http"
)

func (p *PaymentModule) getHandleHistory(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res := ResponseHistory{
		History: p.History,
	}

	p.mtx.Lock()
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		return
	}
	p.mtx.Unlock()
}

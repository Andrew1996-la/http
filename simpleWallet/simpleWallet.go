package simpleWallet

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	mtx   sync.Mutex
	bank  atomic.Int64
	money atomic.Int64
)

func handlePay(w http.ResponseWriter, r *http.Request) {
	httpBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Произошла ошибка при чтении данных", err)
		return
	}

	paymentAmount, err := strconv.Atoi(string(httpBody))
	if err != nil {
		fmt.Fprintf(w, "Ошибка конвертацииn данных\n", err)
		return
	}

	// пропускаю через мьютекс критическую секцию потому что из-за одновременного запроса
	// баланс может стать отрицательным, из-за одновременного прохождения условия
	mtx.Lock()
	if money.Load()-int64(paymentAmount) >= 0 {
		newBalance := money.Add(int64(-paymentAmount))
		fmt.Fprintf(w, "Операция успешна. Остаток - %d\n", newBalance)
	} else {
		fmt.Fprintf(w, "Недостаточно средств\n")
	}
	mtx.Unlock()
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	httpBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Произошла ошибка при чтении данных", err)
		return
	}

	paymentSave, err := strconv.Atoi(string(httpBody))

	if err != nil {
		fmt.Fprintf(w, "Ошибка конвертацииn данных\n", err)
		return
	}

	// пропускаю через мьютекс критическую секцию потому что из-за одновременного запроса
	// баланс может стать отрицательным, из-за одновременного прохождения условия
	mtx.Lock()
	if int64(paymentSave) <= money.Load() {
		bank.Add(int64(paymentSave))
		money.Add(int64(-paymentSave))

		fmt.Fprintf(w, "На баковский счет поступило - %d\n", paymentSave)
		fmt.Fprintf(w, "Баланс банка - %d\n", bank.Load())

		fmt.Fprintf(w, "Наличных осталось - %d\n", money.Load())
	} else {
		fmt.Fprintf(w, "Недостаточно средств\n")
	}
	mtx.Unlock()
}

func SimpleWallet() {
	money.Add(50)

	http.HandleFunc("/pay", handlePay)
	http.HandleFunc("/save", handleSave)

	fmt.Println("Сервер запускается")

	http.ListenAndServe(":8080", nil)
}

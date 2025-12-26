package simpleWallet

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

var (
	mtx   sync.Mutex
	bank  int
	money = 1000
)

func handlePay(w http.ResponseWriter, r *http.Request) {
	httpBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "Произошла ошибка при чтении данных")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	paymentAmount, err := strconv.Atoi(string(httpBody))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, "Ошибка конвертации данных")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// пропускаю через мьютекс критическую секцию потому что из-за одновременного запроса
	// баланс может стать отрицательным, из-за одновременного прохождения условия
	mtx.Lock()
	if money-paymentAmount >= 0 {
		money -= paymentAmount
		_, err := fmt.Fprintf(w, "Операция успешна. Остаток - %d\n", money)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err := fmt.Fprintf(w, "Недостаточно средств\n")
		if err != nil {
			fmt.Println(err)
		}
	}
	mtx.Unlock()
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	httpBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "Произошла ошибка при чтении данных")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	paymentSave, err := strconv.Atoi(string(httpBody))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, "Ошибка конвертации данных\n")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// пропускаю через мьютекс критическую секцию потому что из-за одновременного запроса
	// баланс может стать отрицательным, из-за одновременного прохождения условия
	mtx.Lock()
	if paymentSave <= money {
		bank += paymentSave
		money -= paymentSave

		_, err := fmt.Fprintf(
			w,
			"На баковский счет поступило - %d\nБаланс банка - %d\nНаличных осталось - %d\n",
			paymentSave,
			bank,
			money,
		)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err := fmt.Fprintf(w, "Недостаточно средств\n")
		if err != nil {
			fmt.Println(err)
		}
	}
	mtx.Unlock()
}

func SimpleWallet() {
	http.HandleFunc("/pay", handlePay)
	http.HandleFunc("/save", handleSave)

	fmt.Println("Сервер запускается")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

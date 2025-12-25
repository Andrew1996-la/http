package main

import (
	"fmt"
	"net/http"
)

func handleCancelPay(w http.ResponseWriter, r *http.Request) {
	str := "Произошла успешная отмена опреции"

	_, err := w.Write([]byte(str))

	if err != nil {
		fmt.Println("Ошибка при отмене операции:", err.Error())
	} else {
		fmt.Println("Отмена операции успешно проведена")
	}
}

func handlePay(w http.ResponseWriter, r *http.Request) {
	str := "Оплата проведена"

	_, err := w.Write([]byte(str))

	if err != nil {
		fmt.Println("Ошибка при оплате:", err.Error())
	} else {
		fmt.Println("Запрос на оплату успешно прошел")
	}
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) {
	str := "Hello World from http!"

	_, err := w.Write([]byte(str))

	if err != nil {
		fmt.Println("При приветствии произошла ошибка:", err.Error())
	} else {
		fmt.Println("Запрос на приветствие выполнился успешно")
	}
}

func main() {
	fmt.Println("Запуск http сервиса")

	http.HandleFunc("/", handleHelloWorld)
	http.HandleFunc("/pay", handlePay)
	http.HandleFunc("/cancel", handleCancelPay)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса", err.Error())
	}

	fmt.Println("Программа завершена")
}

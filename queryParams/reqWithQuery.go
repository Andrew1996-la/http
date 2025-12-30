package queryParams

import (
	"fmt"
	"net/http"
)

// example - localhost:9091?foo=x&boo=y
func handler(w http.ResponseWriter, r *http.Request) {
	firstQuery := r.URL.Query().Get("foo")
	secondQuery := r.URL.Query().Get("boo")

	fmt.Println("firstQuery:", firstQuery)
	fmt.Println("secondQuery:", secondQuery)
}

func ReqWithQuery() {
	http.HandleFunc("/default", handler)
	fmt.Println("Запуск query")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}

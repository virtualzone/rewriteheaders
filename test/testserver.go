package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Set-Cookie", "unifises=abcdefg; Path=/; Secure; HttpOnly")
	w.Header().Add("Set-Cookie", "csrf_token=123456; Path=/; Secure")
	fmt.Fprintf(w, "hello\n")
}

func main() {
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8081", nil)
}

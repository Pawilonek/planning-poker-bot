package main

import (
	"io"
	"net/http"

	"github.com/Pawilonek/scrumpoke/internal/bots"
)

func testHttp(w http.ResponseWriter, r *http.Request) {
    _, err := io.WriteString(w, "This is my website!\n"+ bots.TestMessage())
    if err != nil {
        panic(err)
    }
}

func main() {
	http.HandleFunc("/", testHttp)

	err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}


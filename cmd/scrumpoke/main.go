package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Pawilonek/scrumpoke/internal/bots"
)

func testHttp(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "This is my website!")
}

func main() {
    fmt.Println(bots.TestMessage())

	http.HandleFunc("/", testHttp)

	err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}


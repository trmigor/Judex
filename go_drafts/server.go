package main

import (
    "fmt"
    "log"
    "net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, world!\n")
}

func main() {
    http.HandleFunc("/reg_submit", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

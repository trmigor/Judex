package main

import (
	"log"
	"github.com/dpapathanasiou/go-recaptcha"
)

// ProcessRequest accepts the http.Request object
func ProcessRequest(request string, ip string) (result bool) {
	result, err := recaptcha.Confirm(ip, request)
	if err != nil {
		log.Println("recaptcha server error", err)
	}
	return result
}
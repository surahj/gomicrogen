package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	goutils "github.com/mudphilo/go-utils"
)

func GetLanguage() string {

	language := os.Getenv("LANGUAGE")

	if len(language) == 0 {
		language = "fr"
	}

	return language
}

func GetCurrency() string {

	currency := os.Getenv("CURRENCY")

	if len(currency) == 0 {
		currency = "XOF"
	}

	return currency
}

func HTTPGet(remoteURL string, headers map[string]string, payload map[string]string) (httpStatus int, response string) {

	if payload != nil {
		var fields []string

		for key, value := range payload {
			val := fmt.Sprintf("%s=%v", key, url.QueryEscape(value))

			fields = append(fields, val)
		}

		params := strings.Join(fields, "&")
		remoteURL = fmt.Sprintf("%s?%s", remoteURL, params)
	}

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {
		log.Printf("Wants to GET data to URL... %s", remoteURL)
	}

	req, err := http.NewRequest("GET", remoteURL, nil)
	if err != nil {
		log.Printf("Error making HTTP request... %s", err.Error())

		return 0, ""
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := goutils.NewNetClient().Do(req)
	if err != nil {
		log.Printf("Error making HTTP request... %s", err.Error())

		return 0, ""
	}

	st := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error making HTTP request... %s", err.Error())

		return st, ""
	}

	return st, string(body)
}

package main

import (
	"errors"
	"io"
	"log"
	"net/http"
)

// This file will handle everything related to requests.

// This function will make a get request and return the body.
func getRequestBody(url *string) (body string, err error) {
	res, err := http.Get(*url)
	if err != nil {
		return "", errors.New("error sending request")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("error deferring request body")
		}
	}(res.Body)

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalln("error reading body of request")
	}

	return string(resBody), nil
}

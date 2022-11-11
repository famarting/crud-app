package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

var client dapr.Client

func main() {
	fmt.Println("starting service A app")

	generateCalls(context.Background())

}

func generateCalls(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("generating call")
			makeRequest()
		}
	}
}

func makeRequest() {
	var DAPR_HOST, DAPR_HTTP_PORT string
	var okHost, okPort bool
	if DAPR_HOST, okHost = os.LookupEnv("DAPR_HOST"); !okHost {
		DAPR_HOST = "http://localhost"
	}
	if DAPR_HTTP_PORT, okPort = os.LookupEnv("DAPR_HTTP_PORT"); !okPort {
		DAPR_HTTP_PORT = "3500"
	}

	order := "{\"orderId\":\"" + time.Now().String() + "\"}"

	// TODO TRY MAKE MATHOD API/HELLO
	requestURL := DAPR_HOST + ":" + DAPR_HTTP_PORT + "/v1.0/invoke/" + "service-b" + "/method/" + "hello"
	client := &http.Client{}
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(order))
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	// Adding app id as part of th header
	// req.Header.Add("dapr-app-id", "service-b")

	// Invoking a service
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		return
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Println("Order passed: ", string(result))
}

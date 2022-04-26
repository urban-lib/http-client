### HTTP CLIENT

`go get github.com/urban-lib/http-client`

```go
package main

import (
	httpClient "github.com/urban-lib/http-client"
	"net/http"
	"sync"
	"log"
)

var CL httpClient.Client

func getGoogle(uri string, data interface{}, response chan httpClient.Response, wg *sync.WaitGroup) {
	CL.Request(http.MethodGet, uri, data, http.Header{}, response)
	wg.Done()
}

func main() {
	countRequest := 100
	response := make(chan httpClient.Response, countRequest)
	CL = httpClient.NewClient("host", "proxyType", "proxyHost")

	wg := sync.WaitGroup{}
	for i := 0; i < countRequest; i++ {
		wg.Add(1)
		go getGoogle("https://google.com", nil, response, wg)
	}
	wg.Wait()
	for res := range response {
		log.Printf("=====================================\n")
		log.Printf("\n")
		log.Printf("Response: %#v\n", res)
		log.Printf("\n")
		log.Printf("=====================================")
	}
	close(response)

}
```
package endpoints

import (
	"net/http"
	"fmt"
	"time"
	"io/ioutil"
)

func Dummy(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "https://storage-broker-processor.ingress-stage.svc.cluster.local:8800", nil)
	if err != nil {
		fmt.Print("could not build request")
	}
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Client failed in making request")
	}

	fmt.Printf("response successful: %d", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Print("could not read body")
	}
	fmt.Printf("body response: %s", resBody)

}

package httpwalker

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HandleRequest(w http.ResponseWriter, method string, url string, body io.Reader) (err error) {
	if method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		log.Println("here")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var client = &http.Client{Timeout: 15 * time.Second}

	request := new(http.Request)
	response := new(http.Response)

	request, err = http.NewRequest(method, url, body)
	if err != nil {
		log.Print(err)
		return
	}

	response, err = client.Do(request)
	if err != nil {
		log.Print(err)
		return
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return
	}



	_, err = w.Write(responseBody)
	if err != nil {
		log.Print(err)
		return
	}

	return
}


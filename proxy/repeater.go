package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"proxy/utils"
)

type Repeater struct {
	Proxy *Proxy
}

func (rp *Repeater) HandleRepeater(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		rp.HandleSendToRepeater(w, r)
	} else if r.Method == http.MethodPost {
		rp.HandleSendRequest(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (rp *Repeater) HandleSendToRepeater(w http.ResponseWriter, r *http.Request) {
	request, ok := r.URL.Query()["request"]
	if !ok || len(request[0]) < 1 {
		log.Println("Url Param 'request' is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repeaterFile, err := os.OpenFile("proxy/repeater.txt", os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("Invalid repeater file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(request[0])
	defer repeaterFile.Close()

	bytes, err := os.ReadFile(request[0])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = repeaterFile.Write(bytes)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rp *Repeater) HandleSendRequest(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest("proxy/repeater.txt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	io.Copy(w, resp.Body)
}

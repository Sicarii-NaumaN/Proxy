package proxy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Repeater struct {
	Proxy *Proxy
}

func NewRepeater(proxy *Proxy) *Repeater {
	return &Repeater{Proxy: proxy}
}


func (rp *Repeater) HandleRepeater(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		rp.HandleSendToRepeater(w, r)
	} else {
		rp.HandleSendRequest(w, r)
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
	reqBytes, err := os.Open("proxy/repeater.txt")
	if err != nil {
		log.Println("Invalid repeater file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reader io.Reader = reqBytes
	req, err := http.ReadRequest(bufio.NewReader(reader))
	if err != nil {
		log.Println("Repeater file broken")
		return
	}

	req.URL, err = url.Parse("http://" + req.Host + req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
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
			fmt.Println(key,": ", value)
		}
	}

	io.Copy(w, resp.Body)
}

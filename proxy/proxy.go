package proxy

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Proxy struct {
	// CA specifies the root CA for generating leaf certs for each incoming
	// TLS request.
	CA *tls.Certificate

	// TLSServerConfig specifies the tls.Config to use when generating leaf
	// cert using CA.
	TLSServerConfig *tls.Config

	// TLSClientConfig specifies the tls.Config to use when establishing
	// an upstream connection for proxying.
	TLSClientConfig *tls.Config

	// FlushInterval specifies the flush interval
	// to flush to the client while copying the
	// response body.
	// If zero, no periodic flushing is done.
	FlushInterval time.Duration
}

func (p *Proxy) HandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.HandleHTTPS(w, r)
	} else {
		p.HandleHTTP(w, r)
	}
}

func (p *Proxy) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	// It is an error to set this field in an HTTP client request.
	r.RequestURI = ""
	for key := range r.Header {
		if key == "Proxy-Connection" {
			r.Header.Del(key)
		}
	}

	file, err := os.Create("proxy/history/last_request_" + r.Host + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = r.Write(file)
	if err != nil {
		log.Fatal(err)
	}

	r.Close = true
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	bodyBytes := new(bytes.Buffer)
	if resp.ContentLength == -1 {
		io.Copy(bodyBytes, resp.Body)

		fmt.Println(string(bodyBytes.Bytes()), "<-----------------------")
		resp.ContentLength = int64(len(bodyBytes.Bytes()))
		resp.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes.Bytes())))
		resp.Header.Del("Transfer-Encoding")
	}


	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(key, value)
			//fmt.Println(key,": ", value)
		}
	}

	io.Copy(w, bodyBytes)
}

func (p *Proxy) HandleHTTPS(w http.ResponseWriter, r *http.Request) {
	for key, values := range r.Header {
		for _, value := range values {
			w.Header().Set(key, value)
			fmt.Println(key,": ", value)
		}
	}

	fmt.Println(r.URL, r.Method, r.RequestURI, r.RemoteAddr)

	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

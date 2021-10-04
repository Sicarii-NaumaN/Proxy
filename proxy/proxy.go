package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
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
	//for key, values := range r.Header {
	//	for _, value := range values {
	//		w.Header().Set(key, value)
	//		fmt.Println(key,": ", value)
	//	}
	//}
	for key := range r.Header {
		if key == "Proxy-Connection" {
			r.Header.Del(key)
		}
	}
	r.Close = true
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	//err = resp.Write(w)
	//if err != nil {
	//	fmt.Print(err)
	//}

	for key, values := range r.Header {
		for _, value := range values {
			w.Header().Set(key, value)
			fmt.Println(key,": ", value)
		}
	}

	io.Copy(w, resp.Body)
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

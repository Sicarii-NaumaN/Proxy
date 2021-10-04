package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"proxy/proxy"
)

func init() {

}

func main() {
	var isHttps bool
	flag.BoolVar(&isHttps, "https", false, "Proxy protocol (http or https)")
	flag.Parse()

	proxyServ := proxy.Proxy{}
	proxyServe := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(proxyServ.HandleFunc),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	repeaterServ := &http.Server {
		Addr: ":8081",
		Handler: http.HandlerFunc(proxyServ.HandleRepeater),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Println("Proxy server started on localhost:8080, repeater started 0n localhost:8081")
	go repeaterServ.ListenAndServe()
	if !isHttps {
		log.Fatal(proxyServe.ListenAndServe())
	} else {
		log.Fatal(proxyServe.ListenAndServeTLS("cert.key", "ca.key"))
	}


}
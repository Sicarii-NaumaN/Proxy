package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	param_miner "proxy/param-miner"

	"proxy/proxy"
)

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

	repeaterServ := proxy.Repeater{}
	repeaterServe := &http.Server {
		Addr: ":8081",
		Handler: http.HandlerFunc(repeaterServ.HandleRepeater),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	paramMinerServ := param_miner.ParamMiner{}
	paramMinerServe := &http.Server {
		Addr: ":8082",
		Handler: http.HandlerFunc(paramMinerServ.HandleParamMiner),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Println("Proxy server started on localhost: 8080. " +
				"Repeater started on localhost: 8081. " +
				"Param-miner started on localhost: 8082.")

	go paramMinerServe.ListenAndServe()
	go repeaterServe.ListenAndServe()
	if !isHttps {
		log.Fatal(proxyServe.ListenAndServe())
	} else {
		log.Fatal(proxyServe.ListenAndServeTLS("cert.key", "ca.key"))
	}
}
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
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

	getHistoryServe := &http.Server{
		Addr: ":8081",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				d, e := os.Create("repeater/last.txt")
				if e != nil {
					log.Fatal(e)
				}
				d.Write([]byte("kekekk"))
				d.Close()
			} else {
				d, e := os.Create("repeater/last.txt")
				if e != nil {
					log.Fatal(e)
				}
				d.Write([]byte("lololo"))
				d.Close()
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}


	log.Println("Proxy server started on localhost:8080")
	go getHistoryServe.ListenAndServe()
	if !isHttps {
		log.Fatal(proxyServe.ListenAndServe())
	} else {
		log.Fatal(proxyServe.ListenAndServe())
	}


}
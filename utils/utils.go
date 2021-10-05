package utils

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func ParseRequest(path string) (req *http.Request, err error) {
	reqBytes, err := os.Open(path)
	if err != nil {
		log.Println("Invalid repeater file")
		return
	}

	var reader io.Reader = reqBytes
	req, err = http.ReadRequest(bufio.NewReader(reader))
	if err != nil {
		log.Println("Repeater file broken")
		return
	}

	req.URL, err = url.Parse("http://" + req.Host + req.URL.Path)
	if err != nil {
		return
	}

	return
}

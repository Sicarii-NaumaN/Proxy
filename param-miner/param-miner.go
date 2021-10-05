package param_miner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"proxy/utils"
	"time"
)

type ParamMiner struct {
	Value string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// letterRunes variable uses to generate random challenge
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (pm *ParamMiner) HandleParamMiner(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		pm.HandleGetParams(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (pm *ParamMiner) HandleGetParams(w http.ResponseWriter, r *http.Request) {
	err := pm.guessGetParams()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func (pm *ParamMiner) guessGetParams() error {
	pm.hello()

	queriesFile, err := os.Open("param-miner/resourses/params")
	if err != nil {
		log.Println("Invalid query params file")
		return err
	}
	defer queriesFile.Close()

	var queries []string
	scanner := bufio.NewScanner(queriesFile)
	for scanner.Scan() {
		queries = append(queries, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		log.Println(err)
		return err
	}

	exampleValue := make([]rune, 10)
	for i := range exampleValue {
	exampleValue[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	repeaterRequest, err := utils.ParseRequest("proxy/repeater.txt")
	if err != nil {
		log.Println(err)
		return err
	}

	pm.Value = string(exampleValue)
	for _, query := range queries {
		req, err := http.NewRequest(http.MethodGet, repeaterRequest.URL.String() + "?" + query + pm.Value, nil)
		if err != nil {
			log.Println(err)
		}
		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			log.Println(err)
		}

		b := new(bytes.Buffer)
		io.Copy(b, resp.Body)

		//if strings.Contains(b.String(), pm.Value) {
			fmt.Printf("status: %d ----- length: %d ----- param: { %s }\n", resp.StatusCode, resp.ContentLength, query)
		//}
		resp.Body.Close()
	}
	pm.bye()
	return err
}

func (pm *ParamMiner) hello() {
	fmt.Println("\n-----------------PARAM-MINER PROCESSING STARTED-----------------")
}

func (pm *ParamMiner) bye() {
	fmt.Println("\n-----------------PARAM-MINER PROCESSING FINISHED-----------------")
}
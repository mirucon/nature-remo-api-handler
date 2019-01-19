package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type IRSignal struct {
	Freq   int    `json:"freq"`
	Data   []int  `json:"data"`
	Format string `json:"format"`
}

type Signal struct {
	Name    string   `json:"name"`
	Content IRSignal `json:"content"`
}

type Signals struct {
	Signals []Signal `json:"signals"`
}

type SignalRequest struct {
	Name    string    `json:"name"`
	Content *IRSignal `json:"content"`
}

type Config struct {
	RemoIP string `json:"remo_ip"`
	Port   int    `json:"port"`
}

func loadConfig() Config {
	jsonFile, err := os.Open(`config/config.json`)
	if err != nil {
		fmt.Println("config/config.json is not found.")
		return Config{RemoIP: "", Port: 8000}
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config

	err = json.Unmarshal(byteValue, &config)

	if err != nil {
		fmt.Println("Could not parse config/config.json")
		return Config{RemoIP: "", Port: 8000}
	}

	return config
}

func getApiUrl() string  {
	config := loadConfig()
	var url strings.Builder
	if !strings.HasPrefix(config.RemoIP, "http://") {
		url.WriteString("http://")
	}
	url.WriteString(config.RemoIP)
	url.WriteString("/messages/")

	return url.String()
}

func getSignalsList(writer http.ResponseWriter) Signals {
	jsonFile, err := os.Open(`config/signals.json`)
	if err != nil {
		http.Error(writer, "config/signals.json is not found.", 400)
		return Signals{}
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var signals Signals

	err = json.Unmarshal(byteValue, &signals)

	if err != nil {
		http.Error(writer, "Could not parse config/signals.json", 400)
		return Signals{}
	}

	return signals
}

func main() {
	config := loadConfig()
	var port strings.Builder
	port.WriteString(":")
	port.WriteString(strconv.Itoa(config.Port))

	router := mux.NewRouter()
	router.HandleFunc("/signal/", sendSignal).Methods("POST")
	router.HandleFunc("/signal/", getSignal).Methods("GET")
	log.Fatal(http.ListenAndServe(port.String(), router))
}

func sendSignal(writer http.ResponseWriter, request *http.Request) {

	var signalRequest SignalRequest
	var matchedSignal string

	_ = json.NewDecoder(request.Body).Decode(&signalRequest)

	if signalRequest.Content != nil {
		// "content" が存在したらそちら優先
		byteContent, err := json.Marshal(signalRequest.Content)
		matchedSignal = string(byteContent)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// key のみ指定の時は JSON ファイルより探す
		for _, v := range getSignalsList(writer).Signals {
			if v.Name == signalRequest.Name {
				byteContent, err := json.Marshal(v.Content)
				if err != nil {
					http.Error(writer, "Invalid IR signal in your JSON file.", 400)
					return
				}
				matchedSignal = string(byteContent)
				break
			}
		}
	}

	// Return if no key found from the JSON.
	if len(matchedSignal) <= 0 {
		http.Error(writer, "No matching signal.", 400)
		return
	}

	client := http.Client{}

	req, err := http.NewRequest("POST", getApiUrl(), bytes.NewBuffer([]byte(matchedSignal)))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(writer, err.Error(), 400)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(writer, "Cloud not connect to Nature Remo.", 400)
		fmt.Println(resp)
	}
	defer resp.Body.Close()

	fmt.Println(bytes.NewBuffer([]byte(matchedSignal)))

	json.NewEncoder(writer).Encode(resp.Body)
}

func getSignal(writer http.ResponseWriter, request *http.Request) {
	client := http.Client{}

	req, err := http.NewRequest("GET", getApiUrl(), nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(writer, err.Error(), 400)
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		http.Error(writer, "No signal detected on Nature Remo. Send a signal to it again.", 404)
		return
	} else if resp.StatusCode != http.StatusOK {
		http.Error(writer, "Cloud not connect to Nature Remo.", 400)
		return
	}
	defer resp.Body.Close()

	var irSignal IRSignal
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &irSignal)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(writer).Encode(irSignal)
}

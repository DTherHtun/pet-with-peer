package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const dataFile string = "/var/data/hola.txt"
const serviceName string = "hola-pet.default.svc.cluster.local"

func makeRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func pet(w http.ResponseWriter, r *http.Request) {
	if !(r.URL.Path == "/" || r.URL.Path == "/data") {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(dataFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	switch r.Method {
	case "GET":
		file, err := os.Stat(dataFile)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadFile(dataFile)
		if err != nil {
			log.Fatal(err)
		}
		filedata := string(data)
		if r.URL.Path == "/data" {
			if file.Size() == 0 {
				fmt.Fprintf(w, "No data posted yet.")
			} else {
				fmt.Fprintf(w, filedata)
			}
		} else {
			fmt.Fprintf(w, "Current Pod Name => %s \n", hostName)
			fmt.Fprintf(w, "Data stored in the cluster:\n")
			_, srvs, err := net.LookupSRV("", "", serviceName)
			if err != nil {
				log.Fatalf("Could not look up DNS SRV records: %v", err)
			}
			for _, srv := range srvs {
				host := fmt.Sprintf(strings.Trim(srv.Target, ".$"))
				url := "http://" + host + ":8080/data"
				fmt.Println(url)
				content := makeRequest(url)
				fmt.Fprintf(w, "- "+host+": "+string(content)+"\n")
			}
		}

	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		bodyData := string(body)
		datawriter := bufio.NewWriter(file)
		datawriter.WriteString(bodyData + "\n\r")
		fmt.Println("New data has been received and stored.")
		datawriter.Flush()
		fmt.Fprintf(w, "Data stored on pod %s \n", hostName)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.\n")
	}
}

func main() {
	http.HandleFunc("/", pet)
	fmt.Printf("Starting pet server at 8080.......... \n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

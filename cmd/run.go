/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var FilePath string
var wg sync.WaitGroup

type initData struct {
	Url      string `json:"Url"`
	Requests string `json:"Requests"`
	Pacing   string `json:"Pacing"`
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "For a run the test",
	Long:  `HUETA`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonfile, err := os.Open(FilePath)

		if err != nil {
			fmt.Println(err)
		}

		defer jsonfile.Close()

		byteValue, _ := ioutil.ReadAll(jsonfile)

		var dataForRun initData

		json.Unmarshal(byteValue, &dataForRun)

		requestInt, _ := strconv.Atoi(dataForRun.Requests)
		pacing, _ := strconv.Atoi(dataForRun.Pacing)
		runVusers(dataForRun.Url, requestInt, pacing)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&FilePath, "filePath", "f", "", "Path to file")
}

func runVusers(url string, parallelRequests int, pacing int) {

	var count int
	fmt.Println(parallelRequests)
	if parallelRequests > 0 && pacing > 0 {
		wg.Add(parallelRequests)
		for i := 0; i < parallelRequests; i++ {
			time.Sleep(time.Duration(pacing) * time.Second)
			go sendGetRequest(url)
			count++
		}
	} else if parallelRequests > 0 {
		wg.Add(parallelRequests)
		for i := 0; i < parallelRequests; i++ {
			time.Sleep(time.Duration(pacing) * time.Second)
			go sendGetRequest(url)
			count++
		}
	} else {
		log.Fatal("Error! Cant run test!")
	}

	fmt.Println("Sending requests: ", count)
	wg.Wait()
}

func sendGetRequest(url string) string {
	defer wg.Done()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Errored when sending request to the server")
		os.Exit(1)
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	strBody := string(responseBody)
	timeNow := time.Now()

	result := timeNow.Format("2006-01-02 03:04:05") + ", " + req.Method + ", " + resp.Status + ", " + resp.Request.Host
	fmt.Println(result)

	return strBody
}

//Эта функция нужна для расчета времени ответа от HTTP запроса. Требует доработки!
//Был спизжен из stackOverFlow
func timeGet(url string) {
	req, _ := http.NewRequest("GET", url, nil)

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Done: %v\n", time.Since(dns))
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("Connect time: %v\n", time.Since(connect))
		},

		GotFirstResponseByte: func() {
			fmt.Printf("Time from start to first byte: %v\n", time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}

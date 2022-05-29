/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var FilePath string

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

		vusersInt, _ := strconv.Atoi(dataForRun.Vusers)

		fmt.Println(vusersInt)
		runVusers(dataForRun.Url, vusersInt)

		var s int
		fmt.Print("Enter anything for exit:")
		fmt.Scanln(&s)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&FilePath, "filePath", "f", "", "Path to file")
}

type initData struct {
	Url    string `json:"Url"`
	Vusers string `json:"Vusers"`
}

func runVusers(url string, vusers int) {

	var count int
	if vusers > 0 {
		for i := 0; i < vusers; i++ {
			go sendRequest(url)
			count++
		}
	} else {
		log.Fatal("Cant run test!")
	}

	fmt.Println("Sending requests: ", count)
}

func sendRequest(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	stringBody := string(body)

	log.Printf(stringBody)

	return stringBody
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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

		fmt.Println(dataForRun.Url)
		fmt.Println(dataForRun.Vusers)
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

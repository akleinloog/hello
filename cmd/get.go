/*
Copyright © 2020 Arnoud Kleinloog

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var clientPort int

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets an hello from the HTTP Server",
	Long: `Calls the localhost and dumps the answer on the standard out.
If no port is specified, it will use port 80 as default.`,
	Run: func(cmd *cobra.Command, args []string) {
		get()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().IntVarP(&clientPort, "port", "p", 80, "port number")
}

func get() {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d", clientPort))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

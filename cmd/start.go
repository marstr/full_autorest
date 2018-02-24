// Copyright © 2018 Microsoft Corporation and contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/full_autorest/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	startPortLong        = "port"
	startPortShort       = "p"
	startPortDefault     = 80
	startPortDescription = "The port that should be used to listen for requests."
)

var startFlags = viper.New()

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/generate", handleGenerate)

		log.Print("starting full_autorest server on port ", startFlags.GetInt(startPortLong))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", startFlags.GetInt(startPortLong)), nil); err == nil {
			log.Print("unable to start server: ", err)
		}
	},
}

func handleGenerate(resp http.ResponseWriter, req *http.Request) {
	log.Print("request received: generate")

	const autorestTimeout = 5 * 60 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), autorestTimeout)
	defer cancel()

	outputLocation, err := ioutil.TempDir("", "full_autorest")
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, err)
		return
	}
	defer os.RemoveAll(outputLocation)

	options := new(model.AutoRestOptions).SetOutputFolder(outputLocation).SetStdout(resp).SetStderr(resp)

	log.Print(
		"err: ",
		model.InvokeAutoRest(
			ctx,
			model.AutoRestLanguageGo,
			[]string{
				"https://github.com/Azure/azure-rest-api-specs/blob/27c79e5cf0a222441b18828ae81551308e84c758/specification/batch/resource-manager/readme.md",
			},
			*options))
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	startCmd.Flags().IntP(startPortLong, startPortShort, startPortDefault, startPortDescription)
	startFlags.BindPFlags(startCmd.Flags())
}

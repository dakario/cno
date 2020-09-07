/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)


// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {

		config , err := LoadCnoConfig()
		if err != nil && os.IsNotExist(err){
			config = &CnoConfig{}
			SaveConfigOnFileSystem(*config)
		}else if err != nil{
			fmt.Println(err)
			return
		}
		if config.ServerUrl, err = read("serveur URL", config.ServerUrl); err!=nil{
			fmt.Println(err)
			return
		}
		if config.OidcUrl, err = read("oidc realm URL", config.OidcUrl); err!=nil{
			fmt.Println(err)
			return
		}
		if config.OidcClientId, err = read("oidc client-id", config.OidcClientId); err!=nil{
			fmt.Println(err)
			return
		}
		if config.OidcClientSecret, err = read("oidc client-secret", config.OidcClientSecret); err!=nil{
			fmt.Println(err)
			return
		}

		err = SaveConfigOnFileSystem(*config)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
	},
}

func read(label, defaultValue string) (string, error){
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(label+" ["+ defaultValue +"]: ")
	input, err := reader.ReadString('\n')
	if err!=nil {
		fmt.Println("typing error!")
		return defaultValue, err
	}
	input = strings.ReplaceAll(input, "\n", "")
	if len(input)>0 {
		defaultValue = input
	}
	return defaultValue, nil
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func SaveConfigOnFileSystem(config CnoConfig) error{
	user, err := user.Current()
	if err != nil {
		return err
	}
	workspace := filepath.Join(user.HomeDir, "/.cno")
	if _, err := os.Stat(workspace); os.IsNotExist(err){
		err := os.Mkdir(workspace,0700)
		_, err = os.Create(filepath.Join(workspace, "/config"))
		if err != nil {
			return err
		}
		err = nil
	}
	os.Remove(filepath.Join(workspace, "/config"))
	configFile, err := os.OpenFile(filepath.Join(workspace, "/config"), os.O_RDWR|os.O_CREATE, 0666)
	defer configFile.Close()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	_, err = configFile.Write(data)
	return err
}
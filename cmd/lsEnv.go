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
	"fmt"

	"github.com/spf13/cobra"
)

// lsEnvCmd represents the lsEnv command
var lsEnvCmd = &cobra.Command{
	Use:   "env",
	Aliases: []string{"envs","environment","environments"},
	Short: "A brief description of your command",
	Long: `A longer description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadCnoConfig()
		if err!=nil {
			fmt.Println(err)
			return
		}
		projectIdFlag, err := cmd.Flags().GetString("project")
		if err != nil {
			fmt.Println(err)
			return
		}
		if projectIdFlag == "" {
			fmt.Println("Project ID is required, set it with --project flag")
			return
		}
		envs, err := getEnvByProject(projectIdFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, env := range envs {
			if env.ID == config.Environment {
				fmt.Println(" * "+env.Name)
			}else{
				fmt.Println("   "+env.Name)
			}
		}
	},
}

func init() {
	lsCmd.AddCommand(lsEnvCmd)
	lsEnvCmd.Flags().String("project", "", "project-id of the environments")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsEnvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsEnvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

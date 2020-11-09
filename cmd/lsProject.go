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

// lsProjectCmd represents the lsProject command
var lsProjectCmd = &cobra.Command{
	Use:   "project",
	Aliases: []string{"projects"},
	Short: "A brief description of your command",
	Long: `A longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadCnoConfig()
		if err!=nil {
			fmt.Println(err)
			return
		}
		projects, err := getAllProjects()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, project := range projects {
			if project.ID == config.Project {
				fmt.Println(" * "+project.Name)
			}else{
				fmt.Println("   "+project.Name)
			}
		}

	},
}

func init() {
	lsCmd.AddCommand(lsProjectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsProjectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsProjectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

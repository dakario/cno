/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not env this file except in compliance with the License.
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


// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "environment",
	Aliases: []string{"env"},
	Short: "Select an environment",
	Long: `This command allows you to configure an environment of the selected project as the default namespace of you kubeconfig.`,
	Run: func(cmd *cobra.Command, args []string) {

		config, err := LoadCnoConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		envIdFlag, err := cmd.Flags().GetString("env-id")
		if err != nil {
			fmt.Println(err)
			return
		}

		if envIdFlag!="" {
			err = setDefaultNamespace(envIdFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
			config.EnvironmentId = envIdFlag
		}else{
			env, err := chooseEnv(config.ProjectId)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = setDefaultNamespace(env.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
			config.EnvironmentId = env.ID
		}

		err = SaveConfigOnFileSystem(*config)
		if err != nil {
			fmt.Println("WARNING error to save data on $HOME/.cno/config. Cause: "+ err.Error())
		}

		fmt.Println("Project selected successfully!")
		fmt.Println("CNO context is generated and setted as the current context of yo kubeConfig!")
		fmt.Println("Execute 'cno select env' to select an environment as your default namespace")
	},
}


func init() {
	selectCmd.AddCommand(envCmd)
	envCmd.Flags().String("env-id", "", "id of the environment you want to select")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not project this file except in compliance with the License.
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


// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Aliases: []string{""},
	Short: "Select a project",
	Long: `This command allows you to have a valid kubeconfig allowing you to interact with the cluster on which the project is deployed.
The generated kubeconfig contains a certificate with your username as CN signed by the cluster k8s.Which will allow the k8s cluster to identify you.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadCnoConfig()
		if err!=nil {
			fmt.Println(err)
			return
		}

		projectIdFlag, err := cmd.Flags().GetString("project-id")
		if err != nil {
			fmt.Println(err)
			return
		}

		var project *Project
		if projectIdFlag!="" {
			project, err = getProject(projectIdFlag)
			if err != nil {
				fmt.Println(err)
				return
			}

		}else {
			company, err := chooseCompany()
			if err != nil {
				fmt.Println(err)
				return
			}

			org, err := chooseOrganization(company.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
			project, err = chooseProject(org.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if len(project.UidAgent)==0 {
			fmt.Println("This project is not yet deployed on a cluster!")
			return
		}

		err = GenerateKubeConfig(project.UidAgent, project.CompanyID, "default")
		if err != nil {
			fmt.Println(err)
			return
		}
		config.CompanyId = project.CompanyID
		config.OrganizationId = project.OrganizationID
		config.ProjectId = project.ID
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
	selectCmd.AddCommand(projectCmd)
	projectCmd.Flags().String("project-id", "", "id of the project you want to select")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

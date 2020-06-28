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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

type node struct {
	Type  interface{}
	Name  string
	Nodes nodes
}

type nodes []*node

func readDirectory(dir string) nodes {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var nodes nodes
	for _, file := range files {
		if switchType(file) == "directory" {
			fullPath := filepath.Join(dir, file.Name())
			nodes = append(nodes, &node{Type: switchType(file), Name: file.Name(), Nodes: readDirectory(fullPath)})
		} else {
			nodes = append(nodes, &node{Type: switchType(file), Name: file.Name()})
		}
	}
	return nodes
}

func switchType(file os.FileInfo) string {
	if file.IsDir() {
		return "directory"
	}
	return "file"
}

func format(nodes nodes, prefix string) string {
	var result string
	for index, node := range nodes {
		edge := index == (len(nodes) - 1)
		var guide string
		var next string
		if edge {
			guide = prefix + "`--"
			next = prefix + "  "
		} else {
			guide = prefix + "|--"
			next = prefix + "|  "
		}

		result += fmt.Sprintf("%s %s\n", guide, node.Name)
		if node.Type == "directory" {
			result += format(node.Nodes, next)
		}
	}
	return result
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "go-tree",
	Run: func(c *cobra.Command, args []string) {
		fmt.Println(args[0])
		nodes := readDirectory(args[0])
		result := format(nodes, "")
		fmt.Println(result)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-tree.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-tree" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-tree")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

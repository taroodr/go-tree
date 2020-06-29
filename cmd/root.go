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
)

var cfgFile string

type node struct {
	Type  interface{}
	Name  string
	Nodes nodes
}

type nodes []*node

type options struct {
	Level uint16
}

var ops = options{}

func readDirectory(dir string, depth uint16, ops options) nodes {
	var nodes nodes
	if ops.Level < depth {
		return nodes
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if switchType(file) == "directory" {
			fullPath := filepath.Join(dir, file.Name())
			nodes = append(nodes, &node{Type: switchType(file), Name: file.Name(), Nodes: readDirectory(fullPath, depth+1, ops)})
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
		fmt.Println("level:", ops.Level)
		searchDirectoryName := "./"
		if len(args) != 0 {
			searchDirectoryName = args[0]
		}
		nodes := readDirectory(searchDirectoryName, 1, ops)
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
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().Uint16VarP(&ops.Level, "level", "L", 65535, "Descend only level directories deep.")
}

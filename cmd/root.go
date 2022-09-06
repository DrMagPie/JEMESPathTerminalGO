/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/DrMagPie/JEMESPathTerminalGO/jpterm"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var Version string
var OUTPUT_MODES = []string{"result", "expression", "quiet"}
var rootCmd = &cobra.Command{
	Use:   "jpterm",
	Short: "JMES Path Terminal",
	Long:  "JMESPath is an expression language for manipulating JSON documents. If you've never heard of JMESPath before, you write a JMESPath expression that when applied to an input JSON document will produces an output JSON document based on the expression you've provided.\n\nYou can check out the JMESPath site for more information https://jmespath.org/.\n\nOne of the best ways to learn the JMESPath language is to experiment by creating your own JMESPath expressions. The JMESPath Terminal makes it easy to see the results of your JMESPath expressions immediately as you type.",
	Run:   run,
}

func Execute(version string) {
	Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("expression", "e", "", "Specify starting expression.")
	rootCmd.Flags().StringP("mode", "m", "expression", fmt.Sprintf("Specify what's printed to stdout when jpterm exits.\nThis can also be changed when jpterm is running using Tab\nOptions: %s", OUTPUT_MODES))
	// rootCmd.Flags().StringP("output", "o", "", "By default, the output is printed to stdout when jpterm exits.\nYou can instead direct the output to a file.")
	rootCmd.Flags().BoolP("version", "v", false, "Prints version.")
}

func scan(s []string, str string) (int, bool) {
	for i, v := range s {
		if v == str {
			return i, true
		}
	}
	return 0, false
}

func in(s []string, str string) bool {
	_, res := scan(s, str)
	return res
}

func idx(s []string, str string) int {
	idx, _ := scan(s, str)
	return idx
}

func run(cmd *cobra.Command, args []string) {
	version, _ := cmd.Flags().GetBool("version")
	if version != false {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	mode, _ := cmd.Flags().GetString("mode")
	idx, isAPartOf := scan(OUTPUT_MODES, mode)
	if !isAPartOf {
		fmt.Printf("Unknown output mode.\nPlease pick one of the following modes: %s\n", OUTPUT_MODES)
		os.Exit(1)
	}
	expression, _ := cmd.Flags().GetString("expression")

	var content []byte
	var err error
	if len(args) != 0 {
		content, err = ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
	} else if len(args) == 0 {
		stat, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
			fmt.Println("Try piping in some text.")
			os.Exit(1)
		}
		content, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Error while reading from stdin: ", err)
		}
	}
	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	// tea.WithAltScreen()
	if err := tea.NewProgram(jpterm.NewModel(data, idx, expression)).Start(); err != nil {
		log.Fatal(err)
	}

}

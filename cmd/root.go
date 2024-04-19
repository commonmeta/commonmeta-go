/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"commonmeta/commonmeta"
	"commonmeta/crossref"
	"commonmeta/datacite"
	"commonmeta/types"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "commonmeta",
	Short: "Convert scholarly metadata from one format to another",
	Long: `Convert scholarly metadata between formats. Currently
supported input formats are Crossref and DataCite DOIs, currently
the only supported output format is Commonmeta. Example usage:

commonmeta 10.5555/12345678`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide an input DOI")
			return
		}
		input := args[0]
		from, _ := cmd.Flags().GetString("from")

		var data types.Data
		var err error
		if from == "crossref" {
			data, err = crossref.FetchCrossref(input)
		} else if from == "datacite" {
			data, err = datacite.FetchDatacite(input)
		}

		if err != nil {
			fmt.Println(err)
		}
		output, jsErr := commonmeta.WriteCommonmeta(data)
		var out bytes.Buffer
		json.Indent(&out, output, "=", "\t")
		fmt.Println(out.String())

		if jsErr != nil {
			fmt.Println(jsErr)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("from", "f", "crossref", "the format to convert from")
	rootCmd.PersistentFlags().StringP("to", "t", "commonmeta", "the format to convert to")
}
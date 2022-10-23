package main

import (
	"github.com/spf13/cobra"
	"github.com/svetlyi/dircompar/compar"
	"github.com/svetlyi/dircompar/dump"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	var (
		searchIn     string
		saveDumpTo   string
		removePrefix string
	)

	var cmdDump = &cobra.Command{
		Use:   "dump",
		Short: "Creates a dump",
		Long:  `Creates a dump.`,
		Run: func(cmd *cobra.Command, args []string) {
			dump.DumpRun(searchIn, saveDumpTo, removePrefix)
		},
	}
	cmdDump.Flags().StringVarP(&searchIn, "searchIn", "s", "", "Location to dump")
	cmdDump.Flags().StringVarP(&saveDumpTo, "saveDumpTo", "t", "", "Location to save the dump")
	cmdDump.Flags().StringVarP(&removePrefix, "removePrefix", "p", "", "Remove prefix (if searchIn=../folder/, than you can pass removePrefix equals '../')")

	var (
		dump1 string
		dump2 string
	)

	var cmdCompare = &cobra.Command{
		Use:   "compare",
		Short: "Compares two dumps",
		Long:  `Compares two dumps.`,
		Run: func(cmd *cobra.Command, args []string) {
			compar.Compare(dump1, dump2)
		},
	}
	cmdCompare.Flags().StringVarP(&dump1, "dump1", "1", "", "Location of the first dump")
	cmdCompare.Flags().StringVarP(&dump2, "dump2", "2", "", "Location of the second dump")

	rootCmd.AddCommand(cmdDump)
	rootCmd.AddCommand(cmdCompare)
	rootCmd.Execute()
}

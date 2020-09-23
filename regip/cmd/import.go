package cmd

import (
	"bufio"
	"github.com/spf13/cobra"
	"os"
	"regip"
)

var recordType string
var recordPlace string

// importCmd takes care of importing records from text files
var importCmd = &cobra.Command{
	Use:   "import [text file to import]",
	Args:  cobra.ExactArgs(1),
	Short: "imports records",
	Long: `imports records into a running instance.
Add usage examples here f.ex (:2020)`,
	Run: func(cmd *cobra.Command, args []string) {
		lgr := CreateLogger("import", regip.CLR_cli)

		// Check inputs

		rt, ok := regip.ParseRecordType(recordType)
		if !ok {
			lgr.Error("Unable to parse record type ", recordType)
			return
		}

		rp, err := regip.ParseHex(recordPlace)
		if err != nil {
			lgr.Error("Unable to parse record place ", recordPlace, " with error ", err)
			return
		}

		// Check local file before connecting to server

		file, err := os.Open(args[0])
		if err != nil {
			lgr.Error("Error opening file: ", err)
			return
		}
		defer file.Close()

		// File is good, check server

		c, ok := CreateClient(lgr)
		if !ok {
			return
		}

		// Verify connectivity
		ok = c.Ping()
		if !ok {
			lgr.Error("pinging server failed")
			return
		}
		// OK, file is good and we're connected

		ch, qch, err := c.Add(regip.RT_record)
		if err != nil {
			lgr.Error("Error getting add channel ", err)
			return
		}

		scan := bufio.NewScanner(file)
		for scan.Scan() { // Iterate by line
			line := scan.Text()
			cleaned := regip.NormalizeString(line)
			rec := regip.NewRecord(cleaned, rt, rp)
			select {
			case ch <- rec:
				lgr.Print("Added record ", rec.ID().String())
				continue
			case mes := <-qch:
				if mes != nil {
					lgr.Error("Got error: ", mes)
				}
				return
			}
		}

		// We're done adding records
		close(ch)
		// Wait for it to be done
		c.Wait()
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&recordType, "type", "t", "", "type of record")
	importCmd.Flags().StringVarP(&recordPlace, "place", "p", "", "id of place")
	importCmd.MarkFlagRequired("type")
	importCmd.MarkFlagRequired("place")
}

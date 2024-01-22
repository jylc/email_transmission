package cmd

import (
	"email_transmission/pkg"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"log"
)

var (
	inputFileName  string
	configFileName string
	inputFilePath  string
	sendTo         string
	prefix         string
	body           string
	sizeLimit      string
)

var rootCmd = &cobra.Command{
	Use: "trans",
	Run: func(cmd *cobra.Command, args []string) {
		if configFileName == "" || !pkg.Exists(configFileName) {
			_ = cmd.Help()
			log.Fatalln("[ERROR] input file does not exists")
		}

		parseBytes, err := humanize.ParseBytes(sizeLimit)
		if err != nil {
			log.Fatalf("[ERROR] parse file size limit failed, %v\n", err)
		}
		transporter := pkg.NewTransporter(configFileName)
		transporter.Transmit(inputFilePath, inputFileName, sendTo, prefix, body, parseBytes)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFileName, "input", "i", "", "input file name")
	rootCmd.Flags().StringVarP(&inputFilePath, "path", "p", "", "input file path")
	rootCmd.Flags().StringVarP(&configFileName, "config", "", "", "config file name")
	rootCmd.Flags().StringVarP(&sendTo, "sendTo", "", "", "send to somebody")
	rootCmd.Flags().StringVarP(&prefix, "prefix", "", "", "mail title prefix")
	rootCmd.Flags().StringVarP(&body, "body", "", "", "mail body")
	rootCmd.Flags().StringVarP(&sizeLimit, "sizeLimit", "", "20MiB", "file size limit")
}

func Execute() error {
	return rootCmd.Execute()
}

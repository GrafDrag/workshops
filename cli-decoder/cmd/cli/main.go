package main

import (
	"cli-decoder/internal/app"
	"cli-decoder/internal/decode"
	"cli-decoder/internal/store"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	configPath string
	json       bool
	xml        bool
	cmd        = &cobra.Command{
		Use:   "parser",
		Short: "CLI application for parse json and xml string",
		Run: func(cmd *cobra.Command, args []string) {
			if !xml && !json {
				log.Fatal("problem no decoding type set\n")
			}

			config := store.NewConfig()
			_, err := toml.DecodeFile(configPath, config)
			if err != nil {
				log.Fatal("problem failed read configuration file\n")
			}

			store, closeFunc, err := store.NewStore(config)
			if err != nil {
				log.Fatal("problem failed get hash store\n")
			}

			var decoder decode.Decoder
			switch {
			case json:
				decoder = &decode.JSONDecoder{
					Store: store,
				}
			case xml:
				decoder = &decode.XMLDecoder{
					Store: store,
				}
			}

			defer closeFunc()

			app.NewCLI(
				decoder,
				os.Stdin,
				os.Stdout,
			).Run()

		},
	}
)

func init() {
	cmd.PersistentFlags().StringVar(&configPath, "config", "configs/cli.toml", "path to config file")

	cmd.PersistentFlags().BoolVar(&json, "json", false, "waiting json string for execute")
	cmd.PersistentFlags().BoolVar(&xml, "xml", false, "waiting xml string for execute")
}

func main() {
	cmd.Execute()
}

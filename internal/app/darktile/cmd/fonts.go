package cmd

import (
	"fmt"

	darktilefont "github.com/kyamel/goatty/internal/app/darktile/font"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listFontsCmd)
}

var listFontsCmd = &cobra.Command{
	Use:          "list-fonts",
	Short:        "List fonts on your system which are compatible with darktile",
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {

		fonts, err := darktilefont.ListFamilies()
		if err != nil {
			return err
		}

		for _, family := range fonts {
			fmt.Println(family)
		}
		return nil
	},
}

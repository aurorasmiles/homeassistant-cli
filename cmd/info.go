package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Provides a general Home Assistant information overview",
	Long: `
This command provide a general information about your Home Assistant system.
The information provide can be useful for sharing when you are encountering
issues or when reporting one on GitHub.
	`,
	Example: `
  ha info
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("info")

		section := "info"
		command := ""
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

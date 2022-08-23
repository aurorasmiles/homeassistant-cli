package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dockerRegistriesDeleteCmd = &cobra.Command{
	Use:     "delete [host]",
	Aliases: []string{"del", "remove"},
	Short:   "Delete docker registry login for specific host",
	Long: `
Remove login for the Docker OCI registry server.
`,
	Example: `
  ha docker registries delete my-docker.example.com"
`,
	ValidArgsFunction: dockerRegistriesDeleteCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("registries delete")

		section := "docker"
		command := "registries/{host}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		host := args[0]

		request.SetPathParams(map[string]string{
			"host": host,
		})

		resp, err := request.Delete(url)

		// returns 200 OK or 400, everything else is wrong
		if err == nil {
			if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
				err = errors.New("Unexpected server response")
				log.Error(err)
			} else if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
				err = errors.New("API did not return a JSON response")
				log.Error(err)
			}
		}

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dockerRegistriesCmd.AddCommand(dockerRegistriesDeleteCmd)
}

func dockerRegistriesDeleteCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("docker", "registries")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["registries"] != nil {
		if registries, ok := data.Data["registries"].(map[string]interface{}); ok {
			for k := range registries {
				ret = append(ret, k)
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}

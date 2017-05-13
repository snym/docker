package weight

import (
"github.com/spf13/cobra"

"github.com/docker/docker/cli"
"github.com/docker/docker/cli/command"
)

// NewServiceCommand returns a cobra command for `service` subcommands
func NewWeightCommand(dockerCli *command.DockerCli) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "weight",
        Short: "Manage weight",
        Args:  cli.NoArgs,
        RunE:  dockerCli.ShowHelp,
        Tags:  map[string]string{"version": "1.24"},
    }
    cmd.AddCommand(
        setWeight(dockerCli),
        showWeight(dockerCli),
        showPcInfo(dockerCli),
    )
    return cmd
}

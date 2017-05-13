package weight

import (
    "golang.org/x/net/context"

    "github.com/docker/docker/cli"
    "github.com/docker/docker/cli/command"
    "github.com/spf13/cobra"
    "github.com/docker/docker/opts"
    "github.com/hashicorp/memberlist"
    "github.com/docker/swarmkit/manager/scheduler"
    "fmt"
    "text/tabwriter"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/swarm"
    "io"
)

type listOptions struct {
    quiet  bool
    filter opts.FilterOpt
}

const (
    listMenuFmt = "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n"
    listItemFmt = "%s\t%s\t%s\t%.2f\t%d\t%d\t%.2f\t%.2f\n"
)

func showPcInfo(dockerCli *command.DockerCli) *cobra.Command {
    opts := listOptions{filter: opts.NewFilterOpt()}

    cmd := &cobra.Command{
        Use:   "showinfo ",
        Short: "show weight in all directions",
        Args:  cli.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            if err := runShowInfo(dockerCli, opts); err != nil {
                fmt.Println("send err")
                return err
            }
            return nil
        },
    }
    return cmd
}

func runShowInfo(dockerCli command.Cli, opts listOptions) error {
    client := dockerCli.Client()
    out := dockerCli.Out()
    ctx := context.Background()

    nodes, err := client.NodeList(
        ctx,
        types.NodeListOptions{Filters: opts.filter.Value()})
    if err != nil {
        return err
    }
    pcInfo, err := memberlist.GetPcInfo(nodes)
    if err != nil {
        return err
    }
    nodesPcInfo, err := scheduler.PcInfoTranslate(pcInfo, "")
    if err != nil {
        return err
    }

    fmt.Println("nodesPcInfo->", nodesPcInfo)
    printTable(out, nodes, nodesPcInfo)
    return nil
}

func printTable(out io.Writer, nodes []swarm.Node, nodesPcInfo map[string] scheduler.NodePcInfo)  {
    writer := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)

    // Ignore flushing errors
    defer writer.Flush()

    fmt.Fprintf(writer, listMenuFmt, "ID", "HOSTNAME", "STATUS", "CPUUSAGE", "MEMFREE", "MEMALL" ,"MEMUSAGE", "WEIGHT")
    for _, node := range nodes {
        if node.Status.State == swarm.NodeStateReady {
            ID := node.ID
            HOSTNAME := node.Description.Hostname
            STATUS := node.Status.State
            CPUUSAGE := nodesPcInfo[ID].CpuUsageRate
            MEMFREE := nodesPcInfo[ID].MemFree
            MEMALL := nodesPcInfo[ID].MemAll
            MEMUSAGE := nodesPcInfo[ID].MemUsageRete
            WEIGHT := nodesPcInfo[ID].Weight

            fmt.Fprintf(
                writer,
                listItemFmt,
                ID,
                HOSTNAME,
                STATUS,
                CPUUSAGE,
                MEMFREE,
                MEMALL,
                MEMUSAGE,
                WEIGHT,
            )
        }
    }
}
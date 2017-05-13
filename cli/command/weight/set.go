package weight

import (

    "github.com/docker/docker/cli"
    "github.com/docker/docker/cli/command"
    "github.com/spf13/cobra"

    "fmt"
    "os"
)

func setWeight(dockerCli *command.DockerCli) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "set [ARG...]",
        Short: "Set CPU MEM NET NMEM",
        Args:  cli.RequiresMinArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            if err := writeInfo(); err != nil {
                return err
            }
            return nil
        },
    }
    return cmd
}

func writeInfo() error {
    if len(os.Args[3:]) != 5 {
        fmt.Println("info wrong")
        return nil
    }

    pcInfo, err := os.Create("/home/mqy/.docker/dockerWeight.txt")

    if err != nil {
        fmt.Println("create file error")
        return err
    }

    defer pcInfo.Close()

    _, err2 := pcInfo.WriteString(os.Args[3]+" "+os.Args[4]+" "+os.Args[5]+" "+os.Args[6]+" "+os.Args[7]+"\n")
    if err2 != nil {
        fmt.Println("write info error")
        return err2
    }
    return nil
}
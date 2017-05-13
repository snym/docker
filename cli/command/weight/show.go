package weight

import (
    "github.com/spf13/cobra"
    "github.com/docker/docker/cli"
    "github.com/docker/docker/cli/command"
    "os"
    "fmt"
    "bufio"
    "strings"
    "strconv"
)

type pcInfoList struct {
    CPU int
    MEM int
    NET int
    NMEM    int
    K   int
}

func showWeight(dockerCli *command.DockerCli) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "show ",
        Short: "show weight in all directions",
        Args:  cli.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            infoList := []string{"CPU", "MEM", "NET", "NMEM", "K"}
            infoArray, err := readInfo();
            if err != nil {
                return err
            }
            for n:=0; n<len(infoArray); n++ {
                fmt.Printf("%s: %d%% \n", infoList[n], infoArray[n])
            }
            return nil
        },
    }
    return cmd
}

func readInfo() ([5]int, error) {
    var List[5] int
    pcInfo, err := os.Open("/home/mqy/.docker/dockerWeight.txt")
    if err != nil {
        fmt.Println("open error")
        return List, err
    }
    defer pcInfo.Close()

    pcBuff := bufio.NewReader(pcInfo)
    pcRead, err2 := pcBuff.ReadString('\n')

    if err2 != nil {
        fmt.Println("read error")
        return List, err
    }

    infoArray := strings.Fields(pcRead)

    for n:=0; n<len(infoArray); n++ {
        t, _ := strconv.Atoi(infoArray[n])
        List[n] = t
    }

    return List, nil
}
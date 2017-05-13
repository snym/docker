package memberlist

import (
    "bufio"
    "io"
    "strings"
    "strconv"
    "os"
    "fmt"
    "time"
    "math"
)

func ReadPcInfo() (CPUUsage float64, MemFree int64, MemAll int64, NETUsage float64) {
    //Mem info
    MemAll, MemFree = ReadMemInfo()
    fmt.Printf("%.2f%%\n", float64(MemAll-MemFree)*100/float64(MemAll))

    //CPU info Net info
    cpuTime1, idle1 := ReadCpuInfo()
    receive1, transmit1 := ReadNetInfo()
    time.Sleep(1000000000)
    cpuTime2, idle2 := ReadCpuInfo()
    receive2, transmit2 := ReadNetInfo()
    CPUUsage = float64((cpuTime2-cpuTime1)-(idle2-idle1)) / float64(cpuTime2-cpuTime1)
    NETUsage = float64((receive2-receive1)+(transmit2-transmit1)) / 1000

    fmt.Printf("%.2f%% \n", CPUUsage*100)
    fmt.Printf("%.2fKB/s \n", NETUsage)

    return CPUUsage, MemFree, MemAll, NETUsage
    return math.Trunc(CPUUsage*1e4+0.5)*1e-4, MemFree, MemAll, math.Trunc(NETUsage*1e4+0.5)*1e-4
}

func ReadCpuInfo() (cpuTime int, idleTime int) {
    infoFile, err := os.Open("/proc/stat")
    if err != nil {
        fmt.Printf("%s \n", "cpuinfo read error")
        return
    }
    defer infoFile.Close()

    infoBuff := bufio.NewReader(infoFile)
    infoRead, _ := infoBuff.ReadString('\n')

    infoArray := strings.Split(infoRead, " ")[2:11]
    cpuTime = 0
    idleTime, _ = strconv.Atoi(infoArray[3])
    for n:=0; n<len(infoArray); n++ {
        t, _ := strconv.Atoi(infoArray[n])
        cpuTime += t
    }


    return cpuTime, idleTime

}


func ReadMemInfo() (memTotal int64, memFree int64) {

    infoFile, err := os.Open("/proc/meminfo")
    if err != nil {
        fmt.Printf("%s \n", "cpuinfo read error")
        return memTotal, memFree
    }
    defer infoFile.Close()

    infoBuff := bufio.NewReader(infoFile)
    infoRead, _ := infoBuff.ReadString('\n')
    infoArray := strings.Fields(infoRead)
    memTotal, _ = strconv.ParseInt(infoArray[1], 10, 64)

    for n:=0; n<4; n++ {
        infoRead, _ := infoBuff.ReadString('\n')
        infoArray := strings.Fields(infoRead)
        t, _ := strconv.ParseInt(infoArray[1], 10, 64)
        if infoArray[0] != "MemAvailable:" {
            memFree += t
        }
    }

    return memTotal, memFree
}

func ReadNetInfo() (receive int, transmit int) {
    infoFile, err := os.Open("/proc/net/dev")
    if err != nil {
        fmt.Printf("%s \n", "cpuinfo read error")
        return receive, transmit
    }
    defer infoFile.Close()
    infoBuff := bufio.NewReader(infoFile)

    for n := 0; ; n++ {
        infoRead, err := infoBuff.ReadString('\n')
        if err == io.EOF {
            return
        }
        infoArray := strings.Fields(infoRead)
        //fmt.Println(infoArray)
        if n > 1 {
            r, _ := strconv.Atoi(infoArray[1])
            receive += r
            t, _ := strconv.Atoi(infoArray[9])
            transmit += t

        }
    }

    return receive, transmit
}

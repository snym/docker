package scheduler

import (
    "strings"
    "strconv"
    "fmt"
    "github.com/hashicorp/memberlist"
    "github.com/docker/swarmkit/manager/getpcinfo"
    "github.com/docker/swarmkit/api"
    "math"
)

func (s *Scheduler)schedulePcInfo() error {
    var nodesIP []getpcinfo.NodeIp

    for _, node := range s.nodeSet.nodes {
        if node.Status.State == api.NodeStatus_READY {
            nodeIp := getpcinfo.NodeIp{}

            addrStr := strings.Split(node.Status.Addr, ".")

            var addr []byte
            for n := 0; n < len(addrStr); n++ {
                t, err := strconv.Atoi(addrStr[n])
                if err != nil {
                    return err
                }
                addr = append(addr, byte(t))
            }

            nodeIp.ID = node.ID
            nodeIp.Addr = addr

            nodesIP = append(nodesIP, nodeIp)
        }
    }

    msg, err := memberlist.SendInfoRequst(nodesIP)
    if err != nil {
        return  err
    }

    nodePcInfo, err := PcInfoTranslate(msg, "c1")
    if err != nil {
        return err
    }

    //fmt.Println("nodePcInfo1->")
    //fmt.Printf("%+v", nodePcInfo)
    //var ids string
    //for id, node := range nodePcInfo {
    //    fmt.Println("nodePcInfo", id, node)
    //    fmt.Println("node", nodePcInfo[id])
    //    fmt.Println("node12", nodePcInfo["16eg6ae17ojvot7jirxv9sbk4"])
    //    fmt.Println(id,  id == "16eg6ae17ojvot7jirxv9sbk4")
    //    fmt.Println(len(id),  len("16eg6ae17ojvot7jirxv9sbk4"))
    //    fmt.Println(reflect.TypeOf(id),  reflect.TypeOf("16eg6ae17ojvot7jirxv9sbk4"))
    //    fmt.Printf("%s %s", id, "16eg6ae17ojvot7jirxv9sbk4")
    //
    //    ids = id
    //}

    for _, node := range s.nodeSet.nodes {
        fmt.Println("id", node.ID)
        if nowResources, err := nodePcInfo[node.ID]; err == true {
            fmt.Println("nowResources->", nowResources)
            node.NowResources = nowResources
            nodeInfo := s.nodeSet.nodes[node.ID]
            nodeInfo.NowResources = nowResources
            s.nodeSet.nodes[node.ID] = nodeInfo
        }
        fmt.Println("NowResources->", node.NowResources)
    }

    fmt.Println("schedulePcInfo->")
    fmt.Printf("%+v", s.nodeSet.nodes)
    return nil
}

func PcInfoTranslate(pcInfo []string, bias string) (map[string] NodePcInfo, error) {
    nodePcInfoArr :=  make(map[string] NodePcInfo)
    for _, infoStr := range pcInfo {
        nodePcInfo := NodePcInfo{}
        defaultWeight := getpcinfo.DefaultResourcesWeight()

        infoArr := strings.Split(infoStr, ":")
        if cpuUsageRate, err := strconv.ParseFloat(infoArr[0], 64); err == nil {
            nodePcInfo.CpuUsageRate = cpuUsageRate
        } else {
            return nil, err
        }

        if memFree, err := strconv.ParseInt(infoArr[1], 10, 64); err == nil {
            nodePcInfo.MemFree = memFree
        } else {
            return nil, err
        }

        if memAll, err := strconv.ParseInt(infoArr[2], 10, 64); err == nil {
            nodePcInfo.MemAll = memAll
        } else {
            return nil, err
        }

        if netUsagerate, err := strconv.ParseFloat(infoArr[3], 64); err == nil {
            nodePcInfo.NETUsageRate = netUsagerate
        } else {
            return nil, err
        }
        nodePcInfo.ID = infoArr[4]
        memUsageRete := float64(nodePcInfo.MemAll-nodePcInfo.MemFree)/float64(nodePcInfo.MemAll)
        nodePcInfo.MemUsageRete = math.Trunc(memUsageRete*1e4+0.5)*1e-4

        cpuWeight := defaultWeight.Cpu * nodePcInfo.CpuUsageRate
        memWeight := defaultWeight.Mem * nodePcInfo.MemUsageRete

        switch bias {
        case "c1": cpuWeight = cpuWeight * defaultWeight.C1
        case "c2": cpuWeight = cpuWeight * defaultWeight.C2
        case "m1": cpuWeight = memWeight * defaultWeight.M1
        case "m2": cpuWeight = memWeight * defaultWeight.M2
        default:

        }

        weight := cpuWeight + memWeight
        nodePcInfo.Weight = math.Trunc(weight*1e4+0.5)*1e-4

        nodePcInfoArr[nodePcInfo.ID] = nodePcInfo
    }

    return nodePcInfoArr, nil
}

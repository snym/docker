package getpcinfo

import (

)

// send ip
type NodeIp struct {
    ID		string
    Addr	[]byte
    //State	swarm.NodeState
}

// weight set
type ResourcesWeight struct {
    Cpu     float64
    Mem     float64
    C1      float64
    C2      float64
    M1      float64
    M2      float64
    Wup     float64
}

func DefaultResourcesWeight() *ResourcesWeight {
    return &ResourcesWeight{
        Cpu:    40.0,
        Mem:    60.0,
        C1:     1.5,
        C2:     3.0,
        M1:     1.5,
        M2:     3.0,
        Wup:    1.5,
    }
}
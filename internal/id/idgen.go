package id

import (
    "github.com/sony/sonyflake"
    "time"
)

type Generator struct {
    sf *sonyflake.Sonyflake
}

func NewGenerator(machineID uint16) *Generator {
    var st sonyflake.Settings
    st.StartTime = time.Now().Add(-1 * time.Minute)
    st.MachineID = func() (uint16, error) { return machineID, nil }
    sf := sonyflake.NewSonyflake(st)
    return &Generator{sf: sf}
}

func (g *Generator) Next() (uint64, error) {
    return g.sf.NextID()
}

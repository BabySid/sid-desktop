package base

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

func GetCPUUsage() float64 {
	percents, _ := cpu.Percent(time.Second*1, false)
	return percents[0]
}

func GetMEMUsage() float64 {
	mem, _ := mem.VirtualMemory()
	return mem.UsedPercent
}

func GetUpTime() uint64 {
	uptime, _ := host.Uptime()
	return uptime
}

func GetDiskPartitions() []string {
	ps, _ := disk.Partitions(false)
	rs := make([]string, len(ps), len(ps))
	for i := range ps {
		rs[i] = ps[i].Device
	}

	return rs
}

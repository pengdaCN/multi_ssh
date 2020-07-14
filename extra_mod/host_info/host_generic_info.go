package host_info

type HostGenericInfo struct {
	User *UserInfo
	Net  []*NetInterInfo
	Disk *DiskInfo
	Cpu  *CpuInfo
	Mem  *MemoryInfo
}

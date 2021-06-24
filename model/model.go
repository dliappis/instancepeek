package model

type InstanceInfo struct {
	InstanceType string
	CPU          CPUInfo
	Disk         LocalDiskInfo
	Memory       MemoryInfo
	Network      NetworkInfo
	Hardware     VMInfo
}

type CPUInfo struct {
	VCPUCount string
}

type LocalDiskInfo struct {
	Typ     string
	Count   string
	SizeGiB string
}

type MemoryInfo struct {
	SizeMiB string
}

type NetworkInfo struct {
	Performance string
}

type VMInfo struct {
	Hypervisor string
	Baremetal  string
}

type Catalog struct {
	Header       []string
	InstanceData []map[string]string
}

package model

type InstanceInfo struct {
	InstanceType Data
	CPU          CPUInfo
	Disk         LocalDiskInfo
	Memory       MemoryInfo
	Network      NetworkInfo
	Hardware     VMInfo
	Meta         map[string]string
}

type CPUInfo struct {
	VCPUCount Data
}

type LocalDiskInfo struct {
	Typ     Data
	Count   Data
	SizeGiB Data
}

type MemoryInfo struct {
	SizeMiB Data
}

type NetworkInfo struct {
	Performance Data
}

type VMInfo struct {
	Hypervisor Data
	Baremetal  Data
}

type Data struct {
	Label, Value string
}

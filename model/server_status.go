package model

// ServerStatus 服务器运行状态信息
type ServerStatus struct {
	CPU         CPUInfo     `json:"cpu"`
	Memory      MemoryInfo  `json:"memory"`
	Disk        DiskInfo    `json:"disk"`
	System      SystemInfo  `json:"system"`
	Uptime      UptimeInfo  `json:"uptime"`
	Network     NetworkInfo `json:"network"`
	ProcessInfo ProcessInfo `json:"process"`
	Timestamp   int64       `json:"timestamp"`
}

// CPUInfo CPU 信息
type CPUInfo struct {
	Cores       int32   `json:"cores"`       // CPU 核心数
	Usage       float64 `json:"usage"`       // CPU 使用率 (%)
	Temperature float64 `json:"temperature"` // CPU 温度 (°C)
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total     uint64  `json:"total"`      // 总内存 (MB)
	Used      uint64  `json:"used"`       // 已用内存 (MB)
	Free      uint64  `json:"free"`       // 空闲内存 (MB)
	UsageRate float64 `json:"usage_rate"` // 使用率 (%)
	Available uint64  `json:"available"`  // 可用内存 (MB)
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Total     uint64  `json:"total"`      // 总容量 (MB)
	Used      uint64  `json:"used"`       // 已用容量 (MB)
	Free      uint64  `json:"free"`       // 空闲容量 (MB)
	UsageRate float64 `json:"usage_rate"` // 使用率 (%)
	Path      string  `json:"path"`       // 挂载点
}

// SystemInfo 系统信息
type SystemInfo struct {
	OS           string `json:"os"`             // 操作系统 (linux/darwin/windows)
	OSVersion    string `json:"os_version"`     // 操作系统版本
	Distro       string `json:"distro"`         // Linux 发行版 (centos/ubuntu/debian等)
	Platform     string `json:"platform"`       // 平台
	HostName     string `json:"host_name"`      // 主机名
	KernelVer    string `json:"kernel_version"` // 内核版本
	Architecture string `json:"architecture"`   // 系统架构
}

// UptimeInfo 系统运行时间
type UptimeInfo struct {
	TotalSeconds uint64 `json:"total_seconds"` // 运行秒数
	Days         uint64 `json:"days"`          // 天数
	Hours        uint64 `json:"hours"`         // 小时
	Minutes      uint64 `json:"minutes"`       // 分钟
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	BytesSent   uint64 `json:"bytes_sent"`   // 发送字节数
	BytesRecv   uint64 `json:"bytes_recv"`   // 接收字节数
	PacketsSent uint64 `json:"packets_sent"` // 发送包数
	PacketsRecv uint64 `json:"packets_recv"` // 接收包数
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	RunningProcesses uint64 `json:"running_processes"` // 运行中的进程数
	AllProcesses     uint64 `json:"all_processes"`     // 总进程数
}

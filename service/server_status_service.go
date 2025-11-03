package service

import (
	"gin/model"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// GetServerStatus 获取服务器运行状态信息
func GetServerStatus() (*model.ServerStatus, error) {
	status := &model.ServerStatus{
		Timestamp: time.Now().Unix(),
	}

	// 获取 CPU 信息
	cpuInfo, err := getCPUInfo()
	if err == nil {
		status.CPU = cpuInfo
	}

	// 获取内存信息
	memInfo, err := getMemoryInfo()
	if err == nil {
		status.Memory = memInfo
	}

	// 获取磁盘信息
	diskInfo, err := getDiskInfo()
	if err == nil {
		status.Disk = diskInfo
	}

	// 获取系统信息
	sysInfo, err := getSystemInfo()
	if err == nil {
		status.System = sysInfo
	}

	// 获取运行时间
	uptimeInfo, err := getUptimeInfo()
	if err == nil {
		status.Uptime = uptimeInfo
	}

	// 获取网络信息
	netInfo, err := getNetworkInfo()
	if err == nil {
		status.Network = netInfo
	}

	// 获取进程信息
	procInfo, err := getProcessInfo()
	if err == nil {
		status.ProcessInfo = procInfo
	}

	return status, nil
}

// getCPUInfo 获取 CPU 信息
func getCPUInfo() (model.CPUInfo, error) {
	info := model.CPUInfo{}

	// CPU 核心数
	cores, err := cpu.Counts(true)
	if err == nil {
		info.Cores = int32(cores)
	}

	// CPU 使用率
	percent, err := cpu.Percent(time.Second, false)
	if err == nil && len(percent) > 0 {
		info.Usage = percent[0]
	}

	return info, nil
}

// getMemoryInfo 获取内存信息
func getMemoryInfo() (model.MemoryInfo, error) {
	info := model.MemoryInfo{}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return info, err
	}

	info.Total = vmStat.Total / 1024 / 1024         // 转换为 MB
	info.Used = vmStat.Used / 1024 / 1024           // 转换为 MB
	info.Free = vmStat.Free / 1024 / 1024           // 转换为 MB
	info.Available = vmStat.Available / 1024 / 1024 // 转换为 MB
	info.UsageRate = vmStat.UsedPercent

	return info, nil
}

// getDiskInfo 获取磁盘信息
func getDiskInfo() (model.DiskInfo, error) {
	info := model.DiskInfo{}

	var diskPath string
	osType := runtime.GOOS

	// 根据操作系统选择磁盘路径
	switch osType {
	case "windows":
		diskPath = "C:\\"
	case "darwin":
		diskPath = "/"
	case "linux":
		diskPath = "/"
	default:
		diskPath = "/"
	}

	// 获取分区信息
	diskStat, err := disk.Usage(diskPath)
	if err != nil {
		return info, err
	}

	info.Total = diskStat.Total / 1024 / 1024 // 转换为 MB
	info.Used = diskStat.Used / 1024 / 1024   // 转换为 MB
	info.Free = diskStat.Free / 1024 / 1024   // 转换为 MB
	info.UsageRate = diskStat.UsedPercent
	info.Path = diskStat.Path

	return info, nil
}

// getSystemInfo 获取系统信息
func getSystemInfo() (model.SystemInfo, error) {
	info := model.SystemInfo{}

	// 基础系统信息
	info.OS = runtime.GOOS
	info.Platform = runtime.GOARCH

	// 主机信息
	hostStat, err := host.Info()
	if err == nil {
		info.HostName = hostStat.Hostname
		info.KernelVer = hostStat.KernelVersion
		info.OSVersion = hostStat.PlatformVersion
	}

	info.Architecture = runtime.GOARCH

	// 获取 Linux 发行版信息
	if runtime.GOOS == "linux" {
		distro := getLinuxDistro()
		info.Distro = distro
	}

	return info, nil
}

// getLinuxDistro 获取 Linux 发行版
func getLinuxDistro() string {
	// 尝试读取 /etc/os-release (Ubuntu/Debian/CentOS 7+)
	distro := getFromOSRelease()
	if distro != "" {
		return distro
	}

	// 备用方案：检查特定文件
	if fileExists("/etc/redhat-release") {
		return "centos"
	}
	if fileExists("/etc/lsb-release") {
		return "ubuntu"
	}
	if fileExists("/etc/debian_version") {
		return "debian"
	}

	return "linux"
}

// getFromOSRelease 从 /etc/os-release 读取发行版信息
func getFromOSRelease() string {
	cmd := exec.Command("cat", "/etc/os-release")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	content := string(output)

	// 提取 ID 字段
	re := regexp.MustCompile(`ID=["']?([^"\n']+)["']?`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		distro := strings.TrimSpace(matches[1])
		return strings.ToLower(distro)
	}

	return ""
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	cmd := exec.Command("test", "-f", path)
	return cmd.Run() == nil
}

// getUptimeInfo 获取系统运行时间
func getUptimeInfo() (model.UptimeInfo, error) {
	info := model.UptimeInfo{}

	hostStat, err := host.Uptime()
	if err != nil {
		return info, err
	}

	info.TotalSeconds = hostStat
	info.Days = hostStat / (24 * 3600)
	info.Hours = (hostStat % (24 * 3600)) / 3600
	info.Minutes = (hostStat % 3600) / 60

	return info, nil
}

// getNetworkInfo 获取网络统计信息
func getNetworkInfo() (model.NetworkInfo, error) {
	info := model.NetworkInfo{}

	netStat, err := net.IOCounters(false)
	if err != nil {
		return info, err
	}

	if len(netStat) > 0 {
		info.BytesSent = netStat[0].BytesSent
		info.BytesRecv = netStat[0].BytesRecv
		info.PacketsSent = netStat[0].PacketsSent
		info.PacketsRecv = netStat[0].PacketsRecv
	}

	return info, nil
}

// getProcessInfo 获取进程信息
func getProcessInfo() (model.ProcessInfo, error) {
	info := model.ProcessInfo{}

	processes, err := process.Processes()
	if err != nil {
		return info, err
	}

	info.AllProcesses = uint64(len(processes))

	// 统计运行中的进程
	runningCount := 0
	for _, proc := range processes {
		status, err := proc.Status()
		if err == nil {
			if status[0] == "R" || status[0] == "S" {
				runningCount++
			}
		}
	}
	info.RunningProcesses = uint64(runningCount)

	return info, nil
}

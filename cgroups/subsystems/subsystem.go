package subsystems

// Subsystem 指系统的某个资源，cpu、memory
// 使用 resource config 对进程进行设置 cgroups
type Subsystem interface {
	// cpu、memory
	Name() string

	// 根据 config 设置 cgroup
	Set(path string, res *ResourceConfig) error

	// 应用 cgroups
	Apply(path string, pid int, res *ResourceConfig) error

	// 移除某个 cgroups
	Remove(path string) error
}

// 解析传入的参数：mem cpu cpuset 得到
type ResourceConfig struct {
	MemoryLimit string
	CpuCfsQuota int
	CpuShare    string
	CpuSet      string
}

var SubsystemsInstances = []Subsystem{
	&MemorySubsystem{},
	&CpuSubSystem{},
	&CpusetSubSystem{},
}

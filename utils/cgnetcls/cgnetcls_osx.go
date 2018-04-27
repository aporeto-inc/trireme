// +build !linux,darwin windows

//Package cgnetcls implements functionality to manage classid for processes belonging to different cgroups
package cgnetcls

type netCls struct {
}

//Creategroup creates a cgroup/net_cls structure and writes the allocated classid to the file.
//To add a new process to this cgroup we need to write to the cgroup file
func (s *netCls) Creategroup(cgroupname string) error {
	return nil
}

//AssignMark writes the mark value to net_cls.classid file.
func (s *netCls) AssignMark(cgroupname string, mark uint64) error {
	return nil
}

//AddProcess adds the process to the net_cls group
func (s *netCls) AddProcess(cgroupname string, pid int) error {
	return nil
}

//RemoveProcess removes the process from the cgroup by writing the pid to the
//top of net_cls cgroup cgroup.procs
func (s *netCls) RemoveProcess(cgroupname string, pid int) error {
	return nil
}

// DeleteCgroup removes the cgroup
func (s *netCls) DeleteCgroup(cgroupname string) error {
	return nil
}

func (s *netCls) Deletebasepath(contextID string) bool {
	return true
}

func (s *netCls) GetCgroupList() {

}

// ListCgroupProcesses lists the processes of the cgroup
func (s *netCls) ListCgroupProcesses(cgroupname string) ([]string, error) {
	return []string{}, nil
}

// ListAllCgroups returns a list of the cgroups that are managed in the Trireme path
func (s *netCls) ListAllCgroups(path string) []string {
	return []string{}
}

//NewCgroupNetController returns a handle to call functions on the cgroup net_cls controller
func NewCgroupNetController(triremepath string, releasePath string) Cgroupnetcls {
	return &netCls{}
}

//NewDockerCgroupNetController returns a handle to call functions on the cgroup net_cls controller
func NewDockerCgroupNetController() Cgroupnetcls {
	return &netCls{}
}

// MarkVal returns a new Mark
func MarkVal() uint64 {
	return 0
}

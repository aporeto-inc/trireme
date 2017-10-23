// +build linux,!darwin,!windows

//Package cgnetcls implements functionality to manage classid for processes belonging to different cgroups
package cgnetcls

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/kardianos/osext"

	"go.uber.org/zap"
)

const (
	// TriremeBasePath is the base path of the Trireme tree in cgroups
	TriremeBasePath = "trireme"
	// CgroupNameTag is the tag for the cgroup name
	CgroupNameTag = "@cgroup_name"
	// CgroupMarkTag is the tag for the cgroup mark
	CgroupMarkTag = "@cgroup_mark"
	// PortTag is the tag for a port
	PortTag = "@usr:port"

	markFile             = "/net_cls.classid"
	procs                = "/cgroup.procs"
	releaseAgentConfFile = "/release_agent"
	notifyOnReleaseFile  = "/notify_on_release"
	//Initialmarkval is the start of mark values we assign to cgroup
	Initialmarkval = 100
)

var basePath = "/sys/fs/cgroup/net_cls"
var markval uint64 = Initialmarkval

//Empty receiver struct
type netCls struct {
	markchan         chan uint64
	ReleaseAgentPath string
}

//Initialize only ince
func init() {
	mountCgroupController()
}

// Creategroup creates a cgroup/net_cls structure and writes the allocated classid to the file.
// To add a new process to this cgroup we need to write to the cgroup file
func (s *netCls) Creategroup(cgroupname string) error {

	//Create the directory structure
	_, err := os.Stat(basePath + procs)
	if os.IsNotExist(err) {
		syscall.Mount("cgroup", basePath, "cgroup", 0, "net_cls,net_prio")
	}

	os.MkdirAll(filepath.Join(basePath, TriremeBasePath, cgroupname), 0700)

	//Write to the notify on release file and release agent files

	if s.ReleaseAgentPath != "" {
		err = ioutil.WriteFile(filepath.Join(basePath, releaseAgentConfFile), []byte(s.ReleaseAgentPath), 0644)
		if err != nil {
			return fmt.Errorf("Failed to register a release agent error %s", err.Error())
		}

		err = ioutil.WriteFile(filepath.Join(basePath, notifyOnReleaseFile), []byte("1"), 0644)
		if err != nil {
			return fmt.Errorf("Failed to write to the notify file %s", err.Error())
		}

		err = ioutil.WriteFile(filepath.Join(basePath, TriremeBasePath, notifyOnReleaseFile), []byte("1"), 0644)
		if err != nil {
			return fmt.Errorf("Failed to write to the notify file %s", err.Error())
		}

		err = ioutil.WriteFile(filepath.Join(basePath, TriremeBasePath, cgroupname, notifyOnReleaseFile), []byte("1"), 0644)
		if err != nil {
			return fmt.Errorf("Failed to write to the notify file %s", err.Error())
		}
	}

	return nil

}

//AssignMark writes the mark value to net_cls.classid file.
func (s *netCls) AssignMark(cgroupname string, mark uint64) error {

	_, err := os.Stat(filepath.Join(basePath, TriremeBasePath, cgroupname))
	if os.IsNotExist(err) {
		return errors.New("Cgroup does not exist")
	}

	//16 is the base since the mark file expects hexadecimal values
	markval := "0x" + (strconv.FormatUint(mark, 16))

	if err := ioutil.WriteFile(filepath.Join(basePath, TriremeBasePath, cgroupname, markFile), []byte(markval), 0644); err != nil {
		return errors.New("Failed to  write to net_cls.classid file for new cgroup")
	}

	return nil
}

// AddProcess adds the process to the net_cls group
func (s *netCls) AddProcess(cgroupname string, pid int) error {

	_, err := os.Stat(filepath.Join(basePath, TriremeBasePath, cgroupname))
	if os.IsNotExist(err) {
		return errors.New("Cannot add process. Cgroup does not exist")
	}

	PID := []byte(strconv.Itoa(pid))
	if err := syscall.Kill(pid, 0); err != nil {
		return nil
	}

	if err := ioutil.WriteFile(filepath.Join(basePath, TriremeBasePath, cgroupname, procs), PID, 0644); err != nil {
		return errors.New("Cannot add process. Failed to add process to cgroup")
	}

	return nil
}

//RemoveProcess removes the process from the cgroup by writing the pid to the
//top of net_cls cgroup cgroup.procs
func (s *netCls) RemoveProcess(cgroupname string, pid int) error {

	_, err := os.Stat(filepath.Join(basePath, TriremeBasePath, cgroupname))
	if os.IsNotExist(err) {
		return errors.New("Cannot clean up process. Cgroup does not exist")
	}

	data, err := ioutil.ReadFile(filepath.Join(basePath, procs))
	if err != nil || !strings.Contains(string(data), strconv.Itoa(pid)) {
		return errors.New("Cannot cleanup process. Process is not a part of this cgroup")
	}

	if err := ioutil.WriteFile(filepath.Join(basePath, procs), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return errors.New("Cannot clean up process. Failed to remove process to cgroup")
	}

	return nil
}

// DeleteCgroup assumes the cgroup is already empty and destroys the directory structure.
// It will return an error if the group is not empty. Use RempoveProcess to remove all processes
// Before we try deletion
func (s *netCls) DeleteCgroup(cgroupname string) error {

	_, err := os.Stat(filepath.Join(basePath, TriremeBasePath, cgroupname))
	if os.IsNotExist(err) {
		zap.L().Debug("Group already deleted", zap.Error(err))
		return nil
	}

	err = os.Remove(filepath.Join(basePath, TriremeBasePath, cgroupname))
	if err != nil {
		return fmt.Errorf("Failed to delete cgroup %s error returned %s", cgroupname, err.Error())
	}

	return nil
}

// GetAssignedMarkVal -- Gets the mark val assigned to the group
func GetAssignedMarkVal(cgroupName string) string {
	mark, _ := ioutil.ReadFile(filepath.Join(basePath, TriremeBasePath, cgroupName, markFile))
	return string(mark[:len(mark)-1])
}

//Deletebasepath removes the base aporeto directory which comes as a separate event when we are not managing any processes
func (s *netCls) Deletebasepath(cgroupName string) bool {

	if cgroupName == TriremeBasePath {
		os.Remove(filepath.Join(basePath, cgroupName))
		return true
	}

	return false
}

func mountCgroupController() {
	mounts, _ := ioutil.ReadFile("/proc/mounts")
	sc := bufio.NewScanner(strings.NewReader(string(mounts)))
	var netCls = false
	var cgroupMount string
	for sc.Scan() {
		if strings.HasPrefix(sc.Text(), "cgroup") {
			cgroupMount = strings.Split(sc.Text(), " ")[1]
			cgroupMount = cgroupMount[:strings.LastIndex(cgroupMount, "/")]
			if strings.Contains(sc.Text(), "net_cls") {
				basePath = strings.Split(sc.Text(), " ")[1]
				netCls = true
				return
			}
		}

	}

	if len(cgroupMount) == 0 {
		zap.L().Error("Cgroups are not enabled or net_cls is not mounted")
		return
	}
	if !netCls {
		basePath = cgroupMount + "/net_cls"
		os.MkdirAll(basePath, 0700)
		syscall.Mount("cgroup", basePath, "cgroup", 0, "net_cls,net_prio")
		return

	}

}
func CgroupMemberCount(cgroupName string) int {
	_, err := os.Stat(filepath.Join(basePath, TriremeBasePath, cgroupName))
	if os.IsNotExist(err) {
		return 0
	}
	data, err := ioutil.ReadFile(filepath.Join(basePath, TriremeBasePath, cgroupName, "cgroup.procs"))
	if err != nil {
		return 0
	}
	return len(data)
}

// NewDockerCgroupNetController returns a handle to call functions on the cgroup net_cls controller
func NewDockerCgroupNetController() Cgroupnetcls {

	controller := &netCls{
		markchan:         make(chan uint64),
		ReleaseAgentPath: "",
	}

	return controller
}

//NewCgroupNetController returns a handle to call functions on the cgroup net_cls controller
func NewCgroupNetController(releasePath string) Cgroupnetcls {
	binpath, _ := osext.Executable()
	controller := &netCls{
		markchan:         make(chan uint64),
		ReleaseAgentPath: binpath,
	}

	if releasePath != "" {
		controller.ReleaseAgentPath = releasePath
	}

	return controller
}

// MarkVal returns a new Mark Value
func MarkVal() uint64 {
	return atomic.AddUint64(&markval, 1)
}

// ListCgroupProcesses lists the processes of the cgroup
func ListCgroupProcesses(cgroupname string) ([]string, error) {

	_, err := os.Stat(filepath.Join(basePath, TriremeBasePath, cgroupname))

	if os.IsNotExist(err) {
		return []string{}, errors.New("Cgroup does not exist")
	}

	data, err := ioutil.ReadFile(filepath.Join(basePath, TriremeBasePath, cgroupname, "cgroup.procs"))
	if err != nil {
		return []string{}, errors.New("Cannot read procs file")
	}

	procs := []string{}

	for _, line := range strings.Split(string(data), "\n") {
		if len(line) > 0 {
			procs = append(procs, string(line))
		}
	}

	return procs, nil
}

// GetCgroupList geta list of all cgroup names
func GetCgroupList() []string {
	cgroupList := []string{}
	filelist, err := ioutil.ReadDir(filepath.Join(basePath, TriremeBasePath))
	if err != nil {
		return []string{}
	}
	for _, file := range filelist {
		if file.IsDir() {
			cgroupList = append(cgroupList, file.Name())
		}
	}
	return cgroupList
}

package gorapl

//this file represents my attemps at algos for uncovering the toplogy of the system we're running on, with respect to the physical CPU packages
//We have a problem: if we're running on a box with more than one physical CPU, we need to figure out what /dev/cpu/$cpu/msr device to access
//So we need to map logical CPUs to physical sockets

//I'm not really sure how portable this algo is
//it is, however, the simplest way to do this. The intel power gaget iterates through each CPU using affinity masks, and runs `cpuid` in a loop to
//figure things out
//This uses  /sys/devices/system/cpu/cpu*/topology/physical_package_id, which is what lscpu does. I *think* geopm does something similar to this.
// func topoPkgCPUMap() ([]int, error) {

// 	sysdir := "/sys/devices/system/cpu/"
// 	var cpuMap map[int]int
// 	re := regexp.MustCompile("cpu([0-9]*)")

// 	files, err := ioutil.ReadDir(sysdir)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, file := range files {
// 		if file.IsDir() && strings.Contains(file.Name(), "cpu") {
// 			fullPkg := fmt.Sprintf("%s%s/topology/physical_package_id")
// 			match := re.FindStringSubmatch("git commit -m '${abc}'")
// 			thread, err := strconv.ParseInt(match[1], 10, 64)
// 			if err != nil {
// 				return nil, err
// 			}
// 			dat, err := ioutil.ReadFile(fullPkg)
// 			if err != nil {
// 				return nil, err
// 			}
// 			phys, err := strconv.ParseInt(string(dat), 10, 64)
// 			if err != nil {
// 				return nil, err
// 			}
// 			cpuMap[int(phys)] = int(thread)
// 		}
// 	}
// }

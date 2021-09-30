package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fearful-symmetry/gorapl/rapl"
	"github.com/pkg/errors"
)

/*
TODO:
Still unsure about the logic for accessing various MSR devices.
In theory, and based on what I've observed from intel's code,
the RAPL MSRs are package-wide, i.e it doesn't matter what MSR you access on
the package in question. Still, I'm not 100% sure.
*/

func main() {
	err := DumpRAPL()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

// DumpRAPL Prints detailed RAPL data
func DumpRAPL() error {
	top, err := topoPkgCPUMap()
	if err != nil {
		return errors.Wrap(err, "error fetching CPU topology")
	}

	fmt.Printf("Topology:\n")
	for pkg, cores := range top {
		fmt.Printf("\tHardware Package: %d, Cores: %v\n", pkg, cores)
	}

	for pkg, cores := range top {
		fmt.Printf("CPU Package %d\n", pkg)
		handler, err := rapl.CreateNewHandler(cores[0], "")
		if err != nil {
			return errors.Wrap(err, "error creating handler")
		}
		units, err := handler.ReadPowerUnit()
		if err != nil {
			return errors.Wrap(err, "error reading power units")
		}
		fmt.Printf("\tUnits: %f Watts; %f Joules; %f Seconds\n", units.PowerUnits, units.EnergyStatusUnits, units.TimeUnits)

		domains := handler.GetDomains()
		fmt.Printf("\tRAPL Domains:\n")
		for _, domain := range domains {
			fmt.Printf("\t\t%s:\n", domain.Name)

			limit, err := handler.ReadPowerLimit(domain)
			if err != nil {
				return errors.Wrapf(err, "error reading Power Limit on domain %s", domain.Name)
			}
			fmt.Printf("\t\t\tLimit 1: %f Watts; Limit 2: %f Watts; locked: %v\n", limit.Limit1.PowerLimit, limit.Limit2.PowerLimit, limit.Lock)

			status, err := handler.ReadEnergyStatus(domain)
			if err != nil {
				return errors.Wrapf(err, "error reading energy status on domain %s", domain.Name)
			}
			fmt.Printf("\t\t\tCumulative Power usage: %f Joules\n", status)

			policy, err := handler.ReadPolicy(domain)
			if err != nil {
				fmt.Printf("\t\t\tRAPL Policy not available on Domain %s\n", domain.Name)
			}
			if err == nil {
				fmt.Printf("\t\t\tRAPL Policy: %d\n", policy)
			}

			pwrInfo, err := handler.ReadPowerInfo(domain)
			if err != nil {
				fmt.Printf("\t\t\tRAPL Power Info not available on Domain %s\n", domain.Name)
			}
			if err == nil {
				fmt.Printf("\t\t\tPower Info: Thermal Spec: %f; Min: %f; Max: %f; Time Window: %f\n", pwrInfo.ThermalSpecPower, pwrInfo.MinPower, pwrInfo.MaxPower, pwrInfo.MaxTimeWindow)
			}
		}

	}

	return nil
}

//this file represents my attempts at algos for uncovering the toplogy of the system we're running on, with respect to the physical CPU packages
//We have a problem: if we're running on a box with more than one physical CPU, we need to figure out what /dev/cpu/$cpu/msr device to access
//So we need to map logical CPUs to physical sockets

//I'm not really sure how portable this algo is
//it is, however, the simplest way to do this. The intel power gaget iterates through each CPU using affinity masks, and runs `cpuid` in a loop to
//figure things out
//This uses /sys/devices/system/cpu/cpu*/topology/physical_package_id, which is what lscpu does. I *think* geopm does something similar to this.
func topoPkgCPUMap() (map[int][]int, error) {

	sysdir := "/sys/devices/system/cpu/"
	cpuMap := make(map[int][]int)

	files, err := ioutil.ReadDir(sysdir)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("cpu[0-9]+")

	for _, file := range files {
		if file.IsDir() && re.MatchString(file.Name()) {

			fullPkg := filepath.Join(sysdir, file.Name(), "/topology/physical_package_id")
			dat, err := ioutil.ReadFile(fullPkg)
			if err != nil {
				return nil, errors.Wrapf(err, "error reading file %s", fullPkg)
			}
			phys, err := strconv.ParseInt(strings.TrimSpace(string(dat)), 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing value from %s", fullPkg)
			}
			var cpuCore int
			_, err = fmt.Sscanf(file.Name(), "cpu%d", &cpuCore)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching CPU core value from string %s", file.Name())
			}
			pkgList, ok := cpuMap[int(phys)]
			if !ok {
				cpuMap[int(phys)] = []int{cpuCore}
			} else {
				pkgList = append(pkgList, cpuCore)
				cpuMap[int(phys)] = pkgList
			}

		}
	}

	return cpuMap, nil
}

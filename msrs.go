package gorapl

type domainMSRs struct {
	PowerLimit   int64
	EnergyStatus int64
	Policy       int64
	PerfStatus   int64
	PowerInfo    int64
}

// The various RAPL domains

//RAPLDomain is a string type that covers the various RAPL domains
type RAPLDomain struct {
	mask uint
	name string
	msrs domainMSRs
}

//Package is the RAPL domain for the CPU package
var Package = RAPLDomain{0x1, "Package", domainMSRs{0x610, 0x611, 0x0, 0x613, 0x614}}

//DRAM is the RAPL domain for the DRAM
var DRAM = RAPLDomain{0x2, "DRAM", domainMSRs{0x618, 0x619, 0x0, 0x61b, 0x61c}}

//PP0 is the RAPL domain for the processor core
var PP0 = RAPLDomain{0x4, "PP0", domainMSRs{0x638, 0x639, 0x63a, 0x63a, 0x0}}

//PP1 is platform-dependant, although it usually referrs to some uncore power plane
var PP1 = RAPLDomain{0x8, "PP1", domainMSRs{0x640, 0x641, 0x642, 0x0, 0x0}}

// struct defs

//PowerLimitSetting specifies a power limit for a given time window
type PowerLimitSetting struct {
	PowerLimit      uint64
	EnableLimit     bool
	ClampingLimit   bool
	TimeWindowLimit uint64
}

//RAPLPowerLimit contains the data in the MSR_[DOMAIN]_POWER_LIMIT MSR
//This MSR containers two power limits. From the SDM:
//"Two power lmits can be specified, corresponding to time windows of different sizes"
//"Each power limit provides independent clamping control that would permit the processor cores to go below OS-requested state to meet the power limits."
type RAPLPowerLimit struct {
	Limit1 PowerLimitSetting
	Limit2 PowerLimitSetting
	Lock   bool
}

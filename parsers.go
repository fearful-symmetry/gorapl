package gorapl

//parsers for turning the raw MSR uint into a struct

//Handle the MSR_[DOMAIN]_POWER_LIMIT MSR
func parsePowerLimit(msr uint64) RAPLPowerLimit {

	var powerLimit RAPLPowerLimit

	powerLimit.Limit1.PowerLimit = msr & 0x7fff
	powerLimit.Limit1.EnableLimit = ((msr >> 15) & 1) == 1
	powerLimit.Limit1.ClampingLimit = ((msr >> 16) & 1) == 1
	powerLimit.Limit1.TimeWindowLimit = (msr >> 17) & 0x7f

	powerLimit.Limit2.PowerLimit = (msr >> 32) & 0x7fff
	powerLimit.Limit2.EnableLimit = ((msr >> 47) & 1) == 1
	powerLimit.Limit2.ClampingLimit = ((msr >> 48) & 1) == 1
	powerLimit.Limit2.TimeWindowLimit = (msr >> 49) & 0x7f

	powerLimit.Lock = ((msr >> 63) & 1) == 1

	return powerLimit
}

//handle the MSR_RAPL_POWER_UNIT MSR
func parsePowerUnit(msr uint64) RAPLPowerUnit {

	var powerUnit RAPLPowerUnit
	powerUnit.PowerUnits = msr & 0xf
	powerUnit.EnergyStatusUnits = (msr >> 8) & 0x1f
	powerUnit.TimeUnits = (msr >> 16) & 0xf

	return powerUnit
}

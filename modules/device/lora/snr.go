package lora

import "math"

// FSPL in dB (dKm = distance in kilometers, fMHz = frequency in MHz)
func FSPLdB(dKm, fMHz float64) float64 {
	return 32.44 + 20*math.Log10(dKm) + 20*math.Log10(fMHz)
}

// Prx in dBm
func Prx_dBm(Pt_dBm, Gt_dBi, Gr_dBi, dKm, fMHz, Lsys_dB float64) float64 {
	return Pt_dBm + Gt_dBi + Gr_dBi - FSPLdB(dKm, fMHz) - Lsys_dB
}

// Sensitivity (dBm) for given BW (Hz), NF (dB), and required SNR (dB)
func Sensitivity_dBm(BWHz, NFdB, requiredSNRdB float64) float64 {
	noiseFloor := -174.0 + 10*math.Log10(BWHz) // dBm
	return noiseFloor + NFdB + requiredSNRdB
}

// Max free-space distance (km) for given minimum receivable power (PrxMin_dBm)
func MaxDistanceKmForPrx(Pt_dBm, Gt_dBi, Gr_dBi, Lsys_dB, PrxMin_dBm, fMHz float64) float64 {
	// allowed FSPL
	FSPLallowed := Pt_dBm + Gt_dBi + Gr_dBi - Lsys_dB - PrxMin_dBm
	exp := (FSPLallowed - 32.44 - 20*math.Log10(fMHz)) / 20.0
	return math.Pow(10, exp)
}

// Radio horizon approx (km) given antenna heights in meters
// d_km ≈ 3.57 * (sqrt(h1) + sqrt(h2))
func RadioHorizonKm(h1m, h2m float64) float64 {
	return 3.57 * (math.Sqrt(h1m) + math.Sqrt(h2m))
}

// Convenience check: canReach returns true if received power is >= sensitivity + margin
func CanReach(Pt_dBm, Gt_dBi, Gr_dBi, dKm, fMHz, Lsys_dB, BWHz, NFdB, requiredSNRdB, fadeMargin_dB float64) (can bool, prx, sens float64) {
	prx = Prx_dBm(Pt_dBm, Gt_dBi, Gr_dBi, dKm, fMHz, Lsys_dB)
	sens = Sensitivity_dBm(BWHz, NFdB, requiredSNRdB)
	can = prx >= (sens + fadeMargin_dB)
	return
}

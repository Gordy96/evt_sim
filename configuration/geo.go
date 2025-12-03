package configuration

type position struct {
	Type string  `hcl:"type,label"`
	Lat  float64 `hcl:"lat"`
	Lon  float64 `hcl:"lon"`
	Elev float64 `hcl:"elev"`
}

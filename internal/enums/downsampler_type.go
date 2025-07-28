package enums

type DownsamplerType string

const (
	DownsamplerRDP           DownsamplerType = "rdp"
	DownsamplerTimeThreshold DownsamplerType = "time_threshold"
)

func (d DownsamplerType) IsValid() bool {
	switch d {
	case DownsamplerRDP, DownsamplerTimeThreshold:
		return true
	default:
		return false
	}
}

func (d DownsamplerType) String() string {
	return string(d)
}

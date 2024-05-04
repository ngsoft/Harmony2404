package util

type Status int

const (
	STATUS_INITIAL Status = iota
	STATUS_ON
	STATUS_OFF
)

func (v Status) Is(s Status) bool {
	return v == s
}
func (v Status) HasStatus() bool {
	return v > STATUS_INITIAL
}
func (v Status) On() bool {
	return v == STATUS_ON
}
func (v Status) Off() bool {
	return v == STATUS_OFF
}

// func (v Status) String() string {
// 	switch v {
// 	case STATUS_ON:
// 		return "STATUS_ON"
// 	case STATUS_OFF:
// 		return "STATUS_OFF"
// 	default:
// 		return "STATUS_INITIAL"
// 	}
// }

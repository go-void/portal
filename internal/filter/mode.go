package filter

type FilterMode int

const (
	NxDomainMode FilterMode = iota
	LocalIPMode
	NoDataMode
	NullMode
)

var modes = []string{"NXDOMAIN", "LOCALIP", "NODATA", "NULL"}

func (m FilterMode) String() string {
	return modes[m]
}

func MethodFromString(m string) (FilterMode, error) {
	for i, method := range modes {
		if m == method {
			return FilterMode(i), nil
		}
	}

	return -1, ErrInvalidFilterMethod
}

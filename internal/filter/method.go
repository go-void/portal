package filter

type FilterMethod int

const (
	NxDomainMethod FilterMethod = iota
	LocalIPMethod
	NoDataMethod
	NullMethod
)

var methods = []string{"NXDOMAIN", "LOCALIP", "NODATA", "NULL"}

func (m FilterMethod) String() string {
	return methods[m]
}

func MethodFromString(m string) (FilterMethod, error) {
	for i, method := range methods {
		if m == method {
			return FilterMethod(i), nil
		}
	}

	return -1, ErrInvalidFilterMethod
}

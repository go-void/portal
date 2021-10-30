package cache

type Status int

const (
	Hit Status = iota
	Miss
	Expired
)

func (s Status) String() string {
	return []string{"HIT", "MISS", "EXPIRE"}[s]
}

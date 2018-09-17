package flow

type Stat struct {
	Request *Request
	Error   error
	Len     int64
}

func (stat Stat) Successful() bool {
	return stat.Error == nil
}

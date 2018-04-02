package eventz

var (
	curID = 0
	ids   = make(chan int)
)

func init() {
	go inc()
}

// ID returns a new ID that is unique
func ID() int {
	return <-ids
}

func inc() int {
	for {
		curID++
		ids <- curID
	}
}

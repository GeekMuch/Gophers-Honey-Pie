package filewatcher

// LogChannel struct used by the filewatcher and logparser.
type LogChannel struct {
	Name string
	Logs chan string
}

// NewLogChannel initializes a new log channel with the given name.
func NewLogChannel(name string) *LogChannel {
	return &LogChannel{
		Name: name,
		Logs: make(chan string),
	}
}

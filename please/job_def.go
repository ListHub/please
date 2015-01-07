package please

// JobDef defines information needed to launch and run a scheduled container
type JobDef struct {
	Name        string
	Schedule    string
	Image       string
	Command     string
	Ports       []string
	Environment map[string]string
}

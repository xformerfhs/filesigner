package cmdline

// CommandLiner is the interface for all command line processors.
type CommandLiner interface {
	Parse([]string) (error, bool)
	PrintUsage()
	ExtractCommandData() error
}

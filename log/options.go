package log

// Options declares logger options for builder
type Options struct {
	Debug     bool
	LogLevel  string
	AppName   string
	AppID     string
	Version   string
	Revision  string
	SentryDSN string
}

// -----------------------------------------------------------------------------

// DefaultOptions defines default logger options
var DefaultOptions = &Options{
	Debug:     false,
	LogLevel:  "info",
	AppName:   "changeme",
	AppID:     "changeme",
	Version:   "0.0.1",
	Revision:  "123456789",
	SentryDSN: "",
}

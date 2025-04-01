package config

import "time"

var (
	Addr            string        // listen address for http server
	DB              string        // database connection string
	LogConfig       string        // log config file
	LogLevel        string        // sets the log level (zap log level values)
	LogFile         string        // log file to write to
	WaitForServices time.Duration // duration to wait for other services to be ready
	TenantID        int           // default tenant id to use
)

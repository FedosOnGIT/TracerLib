package tracer

import (
	"flag"
	"os"
)

type Configuration struct {
	tracerFile string
	flags      Flags
}

type Flags struct {
	OsVersion   string
	Host        string
	DC          string
	VersionName string
	Environment string
	DeviceID    string
}

func InitConfiguration() *Configuration {
	configuration := &Configuration{}
	flag.StringVar(&configuration.tracerFile, "tracer", "", "-tracer allows to pass path for tracer logs")

	flags := Flags{}
	flags.OsVersion = os.Getenv("cloud_minion_kernel")
	flags.Host = os.Getenv("cloud_hostname")
	flags.DC = os.Getenv("cloud_name")
	flags.VersionName = os.Getenv("cloud_version")
	flags.Environment = os.Getenv("tracer_environment")
	flags.DeviceID = os.Getenv("runid")
	configuration.flags = flags

	return configuration
}

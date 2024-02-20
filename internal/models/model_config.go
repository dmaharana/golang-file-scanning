package models

type AppConfig struct {
	Database Database `ini:"database"`
	Data     Data     `ini:"data"`
	Server   Server   `ini:"server"`
}

type Database struct {
	DatabaseServer string `ini:"dbServer"`
	Port           int    `ini:"dbPort"`
	Username       string `ini:"dbUser"`
	Password       string `ini:"dbPassword"`
	DatabaseName   string `ini:"dbName"`
}

type Data struct {
	DataDir    string `ini:"dataDir"`
	SuccessDir string `ini:"processedSuccessDir"`
	FailDir    string `ini:"processedFailDir"`
}

type Server struct {
	MaxCurrJobs int    `ini:"maxConcurrentJobs"`
	LogFile     string `ini:"logfile"`
}

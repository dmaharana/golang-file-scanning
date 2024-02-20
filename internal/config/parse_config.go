package config

import (
	"file/handling/internal/models"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

func LoadConfigFromIni(iniFile string) (config models.AppConfig, err error) {
	cfg, err := ini.Load(iniFile)
	if err != nil {
		// log.Fatal(err)
		WriteDefaultConfig(iniFile)
		log.Printf("Config file not found, wrote default one: %s\n", iniFile)
		log.Println("Please fill the  fields and restart. Exiting...")
		os.Exit(1)
	}

	// Get the value of a specific section and key
	err = cfg.MapTo(&config)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", iniFile, err)
	}

	DecryptCredentials(&config, "")

	// Print the contents of the Config struct
	// log.Printf("Database server: %s\n", config.DatabaseServer)
	// log.Printf("Database port: %d\n", config.Port)
	// log.Printf("Database username: %s\n", config.Username)
	// log.Printf("Database password: %s\n", config.Password)
	// log.Printf("Database name: %s\n", config.DatabaseName)
	// log.Printf("Data directory: %s\n", config.DataDir)
	// log.Printf("Success directory: %s\n", config.SuccessDir)
	// log.Printf("Fail directory: %s\n", config.FailDir)
	// log.Printf("Max concurrent jobs: %d\n", config.MaxCurrJobs)
	// log.Printf("Log file: %s\n", config.LogFile)

	return
}

func CreateDataDirStruct(config models.AppConfig, inifile, inputdir, soutdir, faildir string) (dataDirs models.DataDir) {
	dataDirs.IniFile = inifile
	dataDirs.InputDir = config.Data.DataDir
	if dataDirs.InputDir == "" {
		dataDirs.InputDir = inputdir
	}
	dataDirs.SuccessDir = config.Data.SuccessDir
	if dataDirs.SuccessDir == "" {
		dataDirs.SuccessDir = soutdir
	}
	dataDirs.FailDir = config.Data.FailDir
	if dataDirs.FailDir == "" {
		dataDirs.FailDir = faildir
	}

	// create dirs if not present
	os.MkdirAll(dataDirs.InputDir, os.ModePerm)
	os.MkdirAll(dataDirs.SuccessDir, os.ModePerm)
	os.MkdirAll(dataDirs.FailDir, os.ModePerm)

	return
}

func WriteDefaultConfig(configfile string) (err error) {
	// create a new empty configuration struct
	var cfg models.AppConfig

	// fill in data
	cfg.Database.DatabaseServer = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.Username = "postgres"
	cfg.Database.Password = "password"
	cfg.Database.DatabaseName = ""
	cfg.Data.DataDir = "./data"
	cfg.Data.SuccessDir = "./success"
	cfg.Data.FailDir = "./failure"
	cfg.Server.MaxCurrJobs = 10
	cfg.Server.LogFile = "./logs/app.log"

	return WriteCfg(&cfg, configfile)
}

func WriteCfg(cfg *models.AppConfig, configfile string) (err error) {
	// create a new INI file
	iniFile := ini.Empty()

	// map the config struct to INI file
	err = iniFile.ReflectFrom(&cfg)
	if err != nil {
		log.Fatalf("Failed to reflect config to INI: %v", err)
	}

	// save cfg to configfile
	err = iniFile.SaveTo(configfile)
	if err != nil {
		log.Fatalf("Unable to write default configuration: %v", err)
		return
	}

	return
}

// Decrypt credentials from config
func DecryptCredentials(cfg *models.AppConfig, password string) {
	decryptedDBPass, ok := Decrypt(cfg.Database.Password, password)
	if ok {
		cfg.Database.Password = decryptedDBPass
	}
}

// Encrypt credentials from config
func EncryptCredentials(cfg *models.AppConfig, password string) {
	enCryptedDBPass, ok := Encrypt(cfg.Database.Password, password)
	if ok {
		cfg.Database.Password = enCryptedDBPass
	}
}

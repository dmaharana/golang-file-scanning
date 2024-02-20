package main

import (
	"file/handling/internal/config"
	"file/handling/internal/models"
	"file/handling/internal/services"
	"log"
)

const IniFile = "config.ini"
const InputDir = "./data"
const SuccessOutDir = "./success"
const FailOutDir = "./failed"

func main() {
	log.SetFlags(log.LstdFlags | log.Ldate | log.Lshortfile)

	iniFile, inputDir := services.ParseCmdFlags(IniFile, InputDir)
	//  Load configuration from INI file
	cfg, err := config.LoadConfigFromIni(iniFile)
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	printCfgValues(&cfg)

	if inputDir == "" {
		inputDir = InputDir
	}

	// Populate data dirs
	dataDirs := config.CreateDataDirStruct(cfg, iniFile, inputDir, SuccessOutDir, FailOutDir)

	log.Println("data dirs: ", dataDirs)

	// process all existing files
	services.ProcessAllFilesInInputDir(cfg, dataDirs)

	services.ProcessInputFile(cfg, dataDirs) // Process all files in the input directory and move them to appropriate

	// encrypt creds and save config
	config.EncryptCredentials(&cfg, "")
	config.WriteCfg(&cfg, dataDirs.IniFile)
}

// printCfgValues prints all available fields of cfg and their respective values
func printCfgValues(cfg *models.AppConfig) {
	log.Println("\n--- Configuration Values ---")
	log.Printf("ServerAddr: %v\n", cfg.Database.DatabaseServer)
	log.Printf("DBUser: %v\n", cfg.Database.Username)
	log.Printf("Port:       %d\n", cfg.Database.Port)
}

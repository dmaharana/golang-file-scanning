package services

import (
	"flag"
	"log"
)

func ParseCmdFlags(iniFile string, inputDir string) (string, string) {
	log.Println("Parsing command line flags...")
	nIniFile := flag.String("c", iniFile, "Path to the .ini file.")
	nInputDir := flag.String("i", inputDir, "Data input directory.")
	flag.Parse()

	if *nIniFile == "" || *nInputDir == "" {
		log.Fatalf("Error: Both -c and -i parameters are required.\n\t-c specifies the path of the configuration file.\n\t-i specified the input directory")
	} else {
		return *nIniFile, *nInputDir
	}

	return iniFile, inputDir
}

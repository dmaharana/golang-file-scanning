package services

import (
	"file/handling/internal/models"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

func ProcessInputFile(cfg models.AppConfig, dataDirs models.DataDir) {
	log.Println("Processing file in  data directory: ", dataDirs.InputDir)

	// handle system interrupts gracefully
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)

	// create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// add the directory to the watcher
	err = watcher.Add(dataDirs.InputDir)
	if err != nil {
		log.Printf("Failed to watch %s: %v", dataDirs.InputDir, err)
		return
	}

	// start monitoring for events
	for {
		select {
		case event, ok := <-watcher.Events: // a file has been added or removed
			if !ok { // we've been notified that the watch
				return
			}
			processEvent(&event, dataDirs)
		case err, ok := <-watcher.Errors: // an error occurred while watching
			if !ok {
				return
			} // we've been notified that the watcher stopped
			log.Printf("Error: %v\n", err)
		// check if there is user interrupt
		case sig := <-interruptChan:
			switch sig {
			case os.Interrupt:
				log.Println("\rReceived SIGINT - Stopping...")
				// os.Exit(0)
				return
			default:
				log.Printf("Unknown signal received: %+v\n", sig)
			}
		}
	}
}

func processEvent(event *fsnotify.Event, dataDirs models.DataDir) {
	log.Printf("Event: %v", event.Op)

	// Check if a new file is created
	if event.Op&fsnotify.Create == fsnotify.Create {
		log.Println("New file created:", event.Name)
		// You can perform any action you want here
		// move the file to success directory
		moveFileToDir(event.Name, dataDirs.SuccessDir)
	}

}

func moveFileToDir(filePath string, targetDir string) {
	// add current date and time to the target filename
	t := time.Now()
	ts := fmt.Sprintf("%d%02d%02d_%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	// get file extension
	ext := path.Ext(filePath)
	fname := path.Base(filePath)
	newFileName := fname[:len(fname)-len(ext)] + "_" + ts + ext

	targetPath := path.Join(targetDir, newFileName)

	err := os.Rename(filePath, targetPath)
	if err != nil {
		log.Printf("Could not move file to Success directory: %s\n", err)
		return
	} else {
		log.Printf("Moved '%s' to '%s'.\n", filePath, targetPath)
	}
}

func GetAllFilesInInputDir(dirPath string) ([]string, error) {
	// iterate through all the files in dirPath and store in a  slice of strings
	fileList := []string{}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fileList, err
	}

	for _, f := range files {
		fileList = append(fileList, path.Join(dirPath, f.Name())) // join with dirpath for full path
	}

	return fileList, nil
}

func ProcessAllFilesInInputDir(cfg models.AppConfig, dataDirs models.DataDir) {
	// process each file in inputDir
	// if successful - move to successDir
	// if unsuccessful - move to failedDir
	files, err := GetAllFilesInInputDir(dataDirs.InputDir)
	if err != nil {
		log.Printf("Failed to retrieve list of files from Input Directory: %w", err)
		return
	}

	for _, file := range files {
		ProcessFile(file, dataDirs)
	}
}

func ProcessFile(filename string, dataDirs models.DataDir) {
	// open the file
	moveFileToDir(filename, dataDirs.SuccessDir)
}

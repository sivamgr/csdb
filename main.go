package main

import (
	"log"
	"os"
	"time"
)

var appDataPath string = "/opt/appdata/csdb/"
var dataPath string = "/opt/appdata/csdb/data"
var inboxPath string = "/opt/appdata/csdb/inbox"
var tmpPath string = "/tmp"

type cs struct {
	UTC uint32
	O   uint32
	H   uint32
	L   uint32
	C   uint32
	V   uint32
	OI  uint32
}

func main() {
	log.Printf("CSDB\n")
	os.MkdirAll(appDataPath, os.ModePerm)
	if !exists(appDataPath) {
		appDataPath, _ = os.Getwd()
		appDataPath += "/appdata/csdb/"
		os.MkdirAll(appDataPath, os.ModePerm)
		log.Printf("Failed to create required appdata directories in /opt/appdata. Switching to %s\n", appDataPath)
	}

	tmpPath = os.TempDir()
	dataPath = appDataPath + "/data/"
	inboxPath = appDataPath + "/inbox/"
	os.MkdirAll(dataPath, os.ModePerm)
	os.MkdirAll(inboxPath, os.ModePerm)

	log.Println("Data Folder is ", appDataPath)
	log.Println("Temp Folder is ", tmpPath)

	time.LoadLocation("Asia/Kolkata")

	processIncomingFiles()
	startWebServer()
}

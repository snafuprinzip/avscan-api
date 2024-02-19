package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type logMessage struct {
	TimeStamp  string `json:"timestamp"`
	LogLevel   string `json:"logLevel"`
	Message    string `json:"message"`
	Method     string `json:"method"`
	HTTPStatus int    `json:"http_status"`
	Logger     string `json:"logger"`
	System     System `json:"system"`
}

type System struct {
	Process     string `json:"process"`
	Application string `json:"application"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func logf(loglevel, method string, httpstatus int, format string, a ...any) {
	// print message to stdout
	message := fmt.Sprintf(format, a...)
	log.Printf("[%d] [%s] %s\n", os.Getpid(), strings.ToUpper(loglevel), message)

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to read hostname, setting to localhost")
		hostname = "localhost"
	}

	// logfile with pid as in the original python script
	//logfile := path.Join(Config.Server.LogDir, hostname+"."+strconv.Itoa(os.Getpid())+".json")

	// logfile with application and hostname
	logfile := path.Join(Config.Server.LogDir, strings.ToLower(Config.Global.Name+"-"+hostname+".json"))

	// create logdir if it doesn't exist already
	err = os.MkdirAll(Config.Server.LogDir, 0750)
	if err != nil {
		log.Printf("Unable to create logdir %s: %s", Config.Server.LogDir, err)
	}

	// Check Logsize
	fi, err := os.Stat(logfile)
	if err != nil {
		log.Printf("Unable to open logfile %s: %s\n", logfile, err)
	} else {
		// Check if the log file exceeds max file size and rotate if necessary
		if int64(Config.Server.LogsizeMB*1024*1024) < fi.Size() {
			log.Printf("Log file size exceeded the limit. Creating a new log file.\n")
			os.Rename(logfile, logfile+".old")
			_, err := os.Create(logfile)
			if err != nil {
				log.Printf("Unable to create new log file: %s\n", err)
			}
		}
	}

	// Open logfile in append mode
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Unable to open logfile %s: %s\n", logfile, err)
		return
	}
	defer file.Close()

	// Generating Logmessage
	msg := logMessage{
		TimeStamp:  time.Now().Format(time.DateTime),
		LogLevel:   strings.ToUpper(loglevel),
		Message:    message,
		Method:     method,
		HTTPStatus: httpstatus,
		Logger:     "clamav",
		System: System{
			Process:     strconv.Itoa(os.Getpid()),
			Application: Config.Global.Name,
			Version:     Config.Global.Version,
			Environment: Config.Global.Environment,
		},
	}

	// marshall into json format
	b, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Unable to marshal log message: %s\n", err)
		return
	}

	// convert []byte to string to add a newline char at the end
	s := string(b) + "\n"

	// and write to logfile for filebeat to pick up
	if _, err := file.Write([]byte(s)); err != nil {
		file.Close() // ignore error; Write error takes precedence
		log.Printf("Unable to write message to log file: %s\n", err)
	}
}

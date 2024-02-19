package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const alphanum = `^[a-zA-Z0-9]*$`

type ScanResult struct {
	AppID         string `json:"app-id"`
	CorrelationID string `json:"correlation-id"`
	Filename      string `json:"fname"`
	HTTPStatus    int    `json:"http_status"`
	Infected      int    `json:"infected files"`
	Message       string `json:"message"`
	RemoteIP      string `json:"remote_address"`
	Result        string `json:"result_string"`
	Size          int64  `json:"size"`
	StartDate     string `json:"start date"`
	EndDate       string `json:"end date"`
	TimeNeeded    string `json:"time"`
	UUID          string `json:"uuid"`
	ExitCode      int    `json:"result"`
}

func (res *ScanResult) JSON() string {
	resJSON, err := json.Marshal(res)
	if err != nil {
		logf("error", "scan", http.StatusOK, "Error marshaling scan result: %s", err)
		return ""
	}
	return string(resJSON)
}

func (res *ScanResult) JSONIndent() string {
	resJSON, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		logf("error", "scan", http.StatusOK, "Error marshaling scan result: %s", err)
		return ""
	}
	return string(resJSON)
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	if !isAccessGrantedByIP(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	result := ScanResult{
		ExitCode:   -2,
		HTTPStatus: 200,
		RemoteIP:   getHeaderIP(r),
	}

	appid := r.FormValue("app-id")
	corid := r.FormValue("correlation-id")

	if !regexp.MustCompile(alphanum).MatchString(appid) {
		result.Message = "app-id has invalid format: "
		result.HTTPStatus = 400
		logf("warn", r.RequestURI, http.StatusBadRequest, result.Message+result.JSON())
		return
	}
	result.AppID = appid

	if !regexp.MustCompile(alphanum).MatchString(corid) {
		result.Message = "correlation-id has invalid format: "
		result.HTTPStatus = 400
		logf("warn", r.RequestURI, http.StatusBadRequest, result.Message+result.JSON())
		return
	}
	result.CorrelationID = corid

	// Maximum upload of MaxMB files
	r.ParseMultipartForm(Config.Scanner.MaxMB << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		logf("crit", "scan", http.StatusInternalServerError, "Error retrieving file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	result.Filename = handler.Filename
	result.Size = handler.Size
	result.UUID = uuid.New().String()

	// Create the uploads folder
	temppath := path.Join(Config.Scanner.UploadDir, result.UUID)
	err = os.MkdirAll(temppath, os.ModePerm)
	if err != nil {
		logf("crit", "POST", http.StatusInternalServerError, "Unable to create upload directory %s: %s", temppath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create new file in the uploads directory
	destpath := path.Join(temppath, handler.Filename)
	dst, err := os.Create(destpath)
	if err != nil {
		logf("crit", "POST", http.StatusInternalServerError, "Unable to create upload file %s: %s", destpath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		logf("crit", "POST", http.StatusInternalServerError, "Can't write to destination file %s: %s", destpath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// scan the uploaded file for viruses
	var cmd *exec.Cmd

	if Config.Scanner.RemoteScan {
		cmd = exec.Command("clamdscan", "-c", Config.Scanner.Configpath, "--stream", destpath)
	} else {
		cmd = exec.Command("clamscan", destpath)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		if cmd.ProcessState.ExitCode() > 1 { // ignore virus found message and show only errors
			logf("crit", "POST", http.StatusInternalServerError, "Error executing scan: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// sort command output into result fields
	outputlines := strings.Split(string(out), "\n")
	//	fmt.Println(len(outputlines), ": ", outputlines)

	for _, line := range outputlines {
		switch {
		case strings.Contains(line, "Infected files:"):
			_, infected, _ := strings.Cut(line, ":")
			result.Infected, _ = strconv.Atoi(strings.Trim(infected, " "))
		case strings.Contains(line, "Time:"):
			_, timeneeded, _ := strings.Cut(line, ":")
			result.TimeNeeded = strings.Trim(timeneeded, " \t")
		case strings.Contains(line, "Start Date:"):
			_, sdate, _ := strings.Cut(line, ":")
			result.StartDate = strings.Trim(sdate, " \t")
		case strings.Contains(line, "End Date:"):
			_, edate, _ := strings.Cut(line, ":")
			result.EndDate = strings.Trim(edate, " \t")
		}
	}

	result.ExitCode = cmd.ProcessState.ExitCode()
	switch result.ExitCode {
	case 0:
		_, msg, _ := strings.Cut(outputlines[0], ":")
		result.Result = strings.Trim(msg, " \t")
		result.Message = result.Result
	case 1:
		_, msg, _ := strings.Cut(outputlines[0], ":")
		result.Result = "VIRUS " + strings.Trim(msg, " \t")
		result.Message = result.Result + ": " + string(out)
		logf("warn", r.RequestURI, http.StatusBadGateway, result.Message+result.JSON())
		result.HTTPStatus = 502
	case 2:
		_, msg, _ := strings.Cut(outputlines[0], ":")
		result.Result = "scanning error: " + strings.Trim(msg, " \t")
		result.Message = result.Result
		result.HTTPStatus = 504
	case -1:
		result.Message = "passthrough"
	default:
		result.Message = fmt.Sprintf("Ein Fehler ist aufgetreten. Kontaktieren Sie bitte %s <%s>!", Config.Global.Maintainer, Config.Global.Email)
		result.HTTPStatus = 500
	}

	// delete temporary upload directory
	os.RemoveAll(temppath)

	w.WriteHeader(result.HTTPStatus)
	w.Write([]byte(result.JSONIndent()))
}

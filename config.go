package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

type ConfigStruct struct {
	Global struct {
		Name        string `yaml:"name" json:"name"`
		Title       string `yaml:"title" json:"title"`
		Environment string `yaml:"environment" json:"environment"`
		URL         string `yaml:"url" json:"url"`
		Version     string `yaml:"version" json:"version"`
		Author      string `yaml:"author" json:"author"`
		Team        string `yaml:"team" json:"team"`
		Maintainer  string `yaml:"maintainer" json:"maintainer"`
		Email       string `yaml:"email" json:"email"`
	} `yaml:"global" json:"global"`
	Server struct {
		Listen    string   `yaml:"listen" json:"listen"`
		Port      int      `yaml:"port" json:"port"`
		LogDir    string   `yaml:"logdir" json:"logdir"`
		LogFile   string   `yaml:"logfile" json:"logfile"`
		LogsizeMB int      `yaml:"logsizemb" json:"logsizemb"`
		Access    []string `yaml:"access" json:"access"`
	} `yaml:"server" json:"server"`
	Passthrough struct {
		Appid []string `yaml:"appid" json:"appid"`
	} `yaml:"passthrough" json:"passthrough"`
	Scanner struct {
		Name       string `yaml:"name" json:"name"`
		RemoteScan bool   `yaml:"remote_scan" json:"remoteScan"`
		MaxMB      int64  `yaml:"maxmb" json:"maxmb"`
		UploadDir  string `yaml:"uploaddir" json:"uploaddir"`
		Configpath string `yaml:"configpath" json:"configpath"`
	} `yaml:"scanner" json:"scanner"`
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	if isAccessGrantedByIP(r) {
		s, _ := json.MarshalIndent(Config, "", "\t")
		w.Write([]byte(s))
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
	}
}

var (
	Config *ConfigStruct
)

func DefaultConfig() {
	Config = &ConfigStruct{}
	Config.Global.Name = "AVScan API"
	Config.Server.LogDir = "./log"
	Config.Server.LogsizeMB = 2
}

func ReadConfig(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("%s\n", err)
		log.Printf("Unable to open configuration file for reading.\nUsing default configuration\n")
		Config = &ConfigStruct{}
		Config.Save(filename)
		return err
	}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Config.Global.Author = "Michael Leimenmeier"
	return nil
}

func (c *ConfigStruct) Yaml() []byte {
	y, err := yaml.Marshal(c)
	if err != nil {
		log.Printf("%s\n", err)
		return nil
	}
	return y
}

func (c *ConfigStruct) Save(path string) {
	err := os.WriteFile(path, c.Yaml(), 0600)
	if err != nil {
		log.Printf("Unable to save configuration to file %s.\n%s\n", path, err)
	}
}

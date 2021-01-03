package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

//
// getConfig returns a new decoded Config struct
func getConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// parseFlag will create and parse CLI flags
// and returns the path to be userd elsewhere
func parseFlags() (string, bool, bool, string) {
	var (
		configPath  string
		logPath     string
		showVersion bool
		showHelp    bool
	)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&configPath, "config", fmt.Sprintf("%s/zabbix-trueconf.yml", dir), "path to config file")
	flag.StringVar(&logPath, "logfile", fmt.Sprintf("%s/zabbix-tconf.log", dir), "path to log file")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.Parse()

	return configPath, showVersion, showHelp, logPath
}

func help() {
	_, pname := filepath.Split(os.Args[0])
	fmt.Printf("Usage: \n %s [OPTIONS]\n", pname)
	fmt.Printf(" %s --version\n", pname)
	fmt.Printf(" %s --help\n", pname)
	fmt.Printf("OPTIONS:\n --config=path to config file\n --log=path to log file \n")
	os.Exit(0)
}

func main() {
	// Program start time
	t0 := time.Now()

	// Generate our config based on the config supplied by the user in the flags
	cfgPath, showVersion, showHelp, logPath := parseFlags()

	if showHelp {
		help()
	}

	if showVersion {
		ShowVersion()
	}

	cfg, err := getConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.LogFile == "" {
		cfg.LogFile = logPath
	}

	logfile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Open Log File ERROR!!!", err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	client, err := NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	client.GetTrueConfInfo()

	// Program finish time
	t1 := time.Now()
	if cfg.LogDebug {
		log.Println("Program timing:", t1.Sub(t0))
	}

	os.Exit(0)

}

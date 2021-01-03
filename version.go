package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// PROGNAME ...
	//progName = "tconf-cli"
	// VERSION ...
	progVersion = "1.0"
)

// ShowVersion ...
func ShowVersion() {
	_, pname := filepath.Split(os.Args[0])
	fmt.Printf("%s version %s \n", pname, progVersion)
	fmt.Printf("(c)gkv@gmpro.ru\n")
	os.Exit(0)
}

package main

import (
	"flag"
	"fmt"
	"path/filepath"
	//	"io/ioutil"
	"crypto/rand"
	"errors"
	"io"
	"os"
)

const Usage = "Usage: kiwf filename\n"

func main() {
	overwrites := *flag.Int("overwrites", 3, "How many times to overwrite the file with random data before deleting")
	flag.Parse()

	pattern, err := parseArgs(os.Args)
	if err != nil {
		fatal(err)
	}

	// Ensure we're starting from the correct dir

	fs, err := filepath.Glob(pattern)
	if err != nil {
		fatal(err)
	}

	for _, f := range fs {
		openFileAndKillRecursive(f, overwrites)
	}
}

func openFileAndKillRecursive(fName string, overwrites int) error {
	f, err := os.OpenFile(fName, os.O_WRONLY, 0000)

	if err != nil {
		// TODO: Find a cleaner way to do this
		err = handleDir(fName, overwrites)
		if err != nil {
			return err
		} else {
			return nil
		}
	}

	defer f.Close()

	killFile(f, overwrites)

	fmt.Printf("Successfully removed file %s", f.Name())
	return nil
}

func handleDir(dName string, overwrites int) error {
	d, err := os.Open(dName)

	if err != nil {
		return err
	}

	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		err = openFileAndKillRecursive(dName+"."+string(filepath.Separator)+file, overwrites)
		if err != nil {
			return err
		}
	}

	err = os.Remove(dName)
	if err != nil {
		return err
	}

	return nil
}

func fatal(err error) {
	fmt.Printf("%s", err.Error())
	os.Exit(1)
}

func killFile(f *os.File, overwrites int) error {
	for i := 0; i < overwrites; i++ {
		err := overwriteFile(f)
		if err != nil {
			fatal(err)
		}
	}

	return os.Remove(f.Name())
}

func overwriteFile(f *os.File) error {
	fileInfo, err := f.Stat()
	if err != nil {
		return errors.New("Failed to retrieve file info")
	}

	if fileInfo.IsDir() {
		fatal(errors.New("Trying to remove directory - not handled yet"))
	}

	s := fileInfo.Size()

	if _, err := io.CopyN(f, rand.Reader, s); err != nil {
		return errors.New("Could not write to file")
	}

	return nil
}

func parseArgs(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("No file supplied to remove\n\t" + Usage)
	} else {
		return args[1], nil
	}
}

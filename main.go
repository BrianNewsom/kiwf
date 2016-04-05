package main

import (
	"flag"
	"fmt"
	//	"io/ioutil"
	"crypto/rand"
	"errors"
	"os"
)

const Usage = "Usage: kiwf filename\n"

func main() {
	overwrites := *flag.Int("overwrites", 3, "How many times to overwrite the file with random data before deleting")
	flag.Parse()

	fName, err := parseArgs(os.Args)
	if err != nil {
		fatal(err)
	}

	f, err := os.OpenFile(fName, os.O_WRONLY, 0000)
	if err != nil {
		fatal(err)
	}

	for i := 0; i < overwrites; i++ {
		err = overwriteFile(f)
		if err != nil {
			fatal(err)
		}
	}

	removeFile(f)

	fmt.Printf("Successfully removed file %s", f.Name())
}

func fatal(err error) {
	fmt.Printf("%s", err.Error())
	os.Exit(1)
}

func removeFile(f *os.File) error {
	return os.Remove(f.Name())
}

func overwriteFile(f *os.File) error {
	fileInfo, err := f.Stat()
	if err != nil {
		return errors.New("Failed to retrieve file info")
	}

	s := fileInfo.Size()

	b, err := RandomBytes(s)
	if err != nil {
		return errors.New("Could not retrieve random bytes")
	}

	_, err = f.Write(b)
	if err != nil {
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

func RandomBytes(n int64) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

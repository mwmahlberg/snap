package main

import (
	"fmt"
	"os"

	"github.com/mwmahlberg/snap/pkg"
	dbg "github.com/visionmedia/go-debug"
)

func setup(c config) ([]pkg.Option, error) {
	var debug = dbg.Debug("setup:setup")

	var opts = make([]pkg.Option, 0)

	debug("Setting up files")
	files, err := setUpFiles(c)
	if err != nil {
		debug("error setting up files: %s", err)
		return nil, fmt.Errorf("setting up files: %s", err)
	}
	opts = append(opts, files...)

	if c.Unsnap {
		debug("Setting mode to unsnap")
		opts = append(opts, pkg.Mode(pkg.UNSNAP))
	}

	return opts, nil
}

func setUpFiles(c config) ([]pkg.Option, error) {
	var debug = dbg.Debug("setup:files")

	opts := make([]pkg.Option, 0)

	if isStdin() {
		debug("Reading from StdIn")
		opts = append(opts, fromStdIn...)
	} else if c.InFile != nil {
		debug("Reading from infile %s", c.InFile.Name())

		if c.Stdout {
			debug("writing to StdOut")
			return append(opts, pkg.InFile(c.InFile), pkg.OutFile(os.Stdout)), nil
		} else if c.Unsnap {
			f, err := createUnsnapFromInfile(c.InFile)
			if err != nil {
				return nil, fmt.Errorf("creating outfile: %s", err)
			}
			debug("writing uncompressed data to %s", f.Name())
			return append(opts, pkg.InFile(c.InFile), pkg.OutFile(f)), nil
		}

		o, err := createFile(c.InFile.Name() + c.Suffix)
		if err != nil {
			return nil, fmt.Errorf("creating outfile: %s", err)
		}

		debug("writing compressed data to %s", o.Name())

		return append(opts, pkg.InFile(c.InFile), pkg.OutFile(o)), nil
	} else {
		panic("No Infile")
	}

	return opts, nil
}

func createUnsnapFromInfile(f *os.File) (*os.File, error) {
	in := []rune(f.Name())
	return createFile(string(in[:len(in)-3]))
}

func createFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
}

func isStdin() bool {
	s, _ := os.Stdin.Stat()
	return (s.Mode() & os.ModeCharDevice) == 0

}

// Copyright ©2016-2022 Markus W Mahlberg <markus@mahlberg.io>
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"os"

	"github.com/alecthomas/kong"
	"github.com/mwmahlberg/snap/pkg"
	dbg "github.com/visionmedia/go-debug"
)

var (

	// Commit upon which snap was built, set by build flag
	Commit = "unknown"

	// Version is the semantic version number of the current binary, should be set by build flag
	//
	//     go build -ldflags "-X main.Version=$(git describe --tags $(git rev-list --tags --max-count=1)"
	Version = "0.0.0"

	cfg struct {
		Unsnap  bool     `help:"decompress file instead of compressing it" short:"d"`
		Keep    bool     `help:"keep original file" short:"k"`
		Stdout  bool     `help:"write to stdout" short:"c"`
		Suffix  string   `help:"set the suffix" default:".sz" short:"S"`
		InFile  *os.File `help:"file to (de)compress" arg:"" optional:""`
		Version kong.VersionFlag
	}
)

func main() {

	var debug = dbg.Debug("MAIN")

	ctx := kong.Parse(
		&cfg,
		kong.Name(os.Args[0]),
		kong.Description("(de-)compress files using snappy algorithm"),
		kong.Vars{
			"version": Version + " " + Commit,
		},
	)

	if os.Args[0] == "unsnap" {
		debug("Called as 'unsnap'. Decompressing source file")
		cfg.Unsnap = true
	} else if os.Args[0] == "scat" || os.Args[0] == "szcat" {
		debug("Called as s(z)cat. Decompressing source file to stdout")
		cfg.Keep = true
		cfg.Unsnap = true
		cfg.Stdout = true
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		debug("Reading from STDIN")
		cfg.InFile = os.Stdin
		cfg.Keep = true
		cfg.Stdout = true
	} else if cfg.InFile == nil {
		ctx.Printf("No input file given.")
		os.Exit(2)
	}

	fi, _ := (*cfg.InFile).Stat()
	debug("Infile %s", fi.Name())

	var outFile *os.File
	var outErr error
	var inFileName = cfg.InFile.Name()

	if cfg.Stdout {
		outFile = os.Stdout
		cfg.Keep = true
	} else if cfg.Unsnap {
		in := []rune(inFileName)
		outFile, outErr = os.OpenFile(string(in[:len(in)-3]), os.O_CREATE|os.O_EXCL|os.O_WRONLY, fi.Mode())

	} else {
		outFile, outErr = os.OpenFile(inFileName+cfg.Suffix, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fi.Mode())
	}
	defer func() {
		debug("Closing outfile")
		outFile.Close()
	}()
	ctx.FatalIfErrorf(outErr, "unable to open '%s'", (*cfg.InFile).Name())

	debug("Outfile: %s", outFile.Name())

	inbuf := bufio.NewReader(cfg.InFile)
	outbuf := bufio.NewWriter(outFile)

	defer func() {
		debug("Flushing buffer")
		outbuf.Flush()
	}()

	s := pkg.NewSnapper(inbuf, outbuf)

	if cfg.Unsnap {
		debug("Unsnapping")
		err := s.Unsnap()
		ctx.FatalIfErrorf(err, "Error decompressing file: %s", err)
	} else {
		debug("Snapping")
		err := s.Snap()
		ctx.FatalIfErrorf(err, "error compressing file: %s", err)
	}

	defer func() {
		debug("Closing inFile")
		cfg.InFile.Close()
	}()

	if !cfg.Keep {
		debug("Removing source file after completion")
		err := os.Remove(inFileName)
		ctx.FatalIfErrorf(err, "error while removing '%s': %v", inFileName, err)
	}
}

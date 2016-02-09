// Copyright ©2016 Markus W Mahlberg <markus@mahlberg.io>
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

	"github.com/mwmahlberg/snap/internal"
	    "github.com/andrew-d/go-termutil"


	dbg "github.com/tj/go-debug"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	unsnap = kingpin.Flag("unsnap", "uncompress file").Short('u').Bool()
	keep   = kingpin.Flag("keep", "keep original file").Short('k').Default("false").Bool()

	stdout = kingpin.Flag("stdout", "write to stdout").Short('c').Default("false").Bool()
	suffix = kingpin.Flag("suffix", "changes the default suffix from '.sz' to the given value").Short('S').Default(".sz").String()

	inFile = kingpin.Arg("file", "file to (de)compress").File()

	outFile *os.File
)

func init() {

	debug := dbg.Debug("INIT")
	kingpin.UsageTemplate(kingpin.DefaultUsageTemplate).Version("0.1").Author("Markus W Mahlberg")
	kingpin.CommandLine.Author("Markus W Mahlberg")
	kingpin.CommandLine.Help = "tool to (de-)compress files using snappy algorithm"
	kingpin.CommandLine.HelpFlag.Short('h')

	kingpin.Parse()

	if os.Args[0] == "unsnap" {
		debug("Called as 'unsnap'. Decompressing source file")
		*unsnap = true
	} else if os.Args[0] == "scat" || os.Args[0] == "szcat" {
		debug("Called as s(z)cat. Decompressing source file to stdout")
		*keep = true
		*unsnap = true
		*stdout = true
	}

}

func main() {

	var debug = dbg.Debug("MAIN")

	if !termutil.Isatty(os.Stdin.Fd()){
		debug("Reading from STDIN")
		*inFile = os.Stdin
		*keep = true
		*stdout = true
	}

	fi, err := (*inFile).Stat()
	kingpin.FatalIfError(err, "unable to access '%s'", (*inFile).Name())

	var outFile *os.File
	var outErr error

	if *stdout {
		outFile = os.Stdout
		*keep = true
	} else if *unsnap {
		in := []rune((*inFile).Name())
		outFile, outErr = os.OpenFile(string(in[:len(in)-3]), os.O_CREATE|os.O_EXCL|os.O_WRONLY, fi.Mode())

	} else {
		outFile, outErr = os.OpenFile((*inFile).Name()+*suffix, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fi.Mode())
	}
	defer outFile.Close()

	kingpin.FatalIfError(outErr, "unable to open '%s'", (*inFile).Name())

	debug("Outfile: %s", outFile.Name())

	inbuf := bufio.NewReader(*inFile)
	outbuf := bufio.NewWriter(outFile)
	defer outbuf.Flush()

	s := internal.NewSnapper(inbuf, outbuf)

	if *unsnap {
		err := s.Unsnap()
		kingpin.FatalIfError(err, "error during decompression")
	} else {
		err = s.Snap()
		kingpin.FatalIfError(err, "error during compression")
	}

	(*inFile).Close()

	if !(*keep) {
		debug("Removing source file after completion")
		err := os.Remove((*inFile).Name())
		kingpin.FatalIfError(err, "error while removing '%s': %v", (*inFile).Name(), err)
	}
}

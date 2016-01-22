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
	"os"
	"bufio"

	"github.com/mwmahlberg/snap/internal"

	dbg "github.com/tj/go-debug"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	unsnap  = kingpin.Flag("unsnap", "uncompress file").Short('u').Bool()
	inFile  = kingpin.Arg("file", "file to (de)compress").Required().File()
//	keep    = kingpin.Flag("keep", "keep original file").Short('k').Default("false").Bool()
	outFile *os.File
)

func init() {

	debug := dbg.Debug("INIT")
	kingpin.Version("0.1")
	kingpin.CommandLine.HelpFlag.Short('h')

	kingpin.Parse()

	if os.Args[0] == "unsnap" {
		debug("Called as 'unsnap'. Decompressing source file")
		*unsnap = true
	}

}

func main() {
	var debug = dbg.Debug("MAIN")
	defer (*inFile).Close()
	
	var outFileName string

	if *unsnap {
		in := []rune((*inFile).Name())
		outFileName = string(in[:len(in)-3])
		debug("Name: %s", outFileName)
	} else {
		outFileName = (*inFile).Name() + `.sz`
	}

	fi, err := (*inFile).Stat()
	kingpin.FatalIfError(err,"unable to access '%s'", (*inFile).Name())

	outFile, err := os.OpenFile(outFileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fi.Mode())
	kingpin.FatalIfError(err,"unable to open '%s'",(*inFile).Name())
	defer outFile.Close()

	inbuf := bufio.NewReader(*inFile)
	outbuf := bufio.NewWriter(outFile)
	defer outbuf.Flush()

	s := internal.NewSnapper(inbuf, outbuf)

	if *unsnap {
		err := s.Unsnap()
		kingpin.FatalIfError(err,"error during decompression")
		return
	}

	err = s.Snap()
	kingpin.FatalIfError(err,"error during compression")

}

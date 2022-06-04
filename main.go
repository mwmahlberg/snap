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
	"os"

	"github.com/alecthomas/kong"
	"github.com/mwmahlberg/snap/pkg"
	dbg "github.com/visionmedia/go-debug"
)

type config struct {
	Unsnap  bool     `help:"decompress file instead of compressing it" short:"d"`
	Keep    bool     `help:"keep original file" short:"k"`
	Stdout  bool     `help:"write to stdout" short:"c"`
	Suffix  string   `help:"set the suffix" default:".sz" short:"S"`
	InFile  *os.File `help:"file to (de)compress" arg:"" optional:""`
	outFile *os.File
	Version kong.VersionFlag
}

var (

	// Commit upon which snap was built, set by build flag
	Commit = "unknown"

	// Version is the semantic version number of the current binary, should be set by build flag
	Version = "0.0.0"

	fromStdIn = []pkg.Option{
		pkg.InFile(os.Stdin),
	}

	szcat = []pkg.Option{
		pkg.OutFile(os.Stdout),
		pkg.Mode(pkg.UNSNAP),
	}

	toStdout = []pkg.Option{
		pkg.OutFile(os.Stdout),
	}

	cfg config
)

func main() {

	var debug = dbg.Debug("main")

	var opts = make([]pkg.Option, 0)

	ctx := kong.Parse(
		&cfg,
		kong.Name(os.Args[0]),
		kong.Description("(de-)compress files using snappy algorithm"),
		kong.Vars{
			"version": Version + "-" + Commit,
		},
	)

	if o, err := setup(cfg); err != nil {
		panic(err)
	} else {
		opts = append(opts, o...)
	}

	if os.Args[0] == "unsnap" {
		debug("Called as 'unsnap'. Decompressing source file")
		opts = append(opts, pkg.Mode(pkg.UNSNAP))
	} else if os.Args[0] == "scat" || os.Args[0] == "szcat" || (cfg.Stdout && cfg.Unsnap) {
		opts = append(opts, szcat...)
	}

	if !isStdin() {
		debug("Not reading from stding, defering closing of source file")

		debug("Keep is not set, source file will be removed")
		if !cfg.Keep {
			defer os.Remove(cfg.InFile.Name())
		}

		defer cfg.InFile.Close()

	}
	if !cfg.Stdout {
		debug("Not writing to stdout, defering close of destination file")
		defer cfg.outFile.Close()
	}

	s := pkg.NewSnapper(opts...)
	err := s.Do()

	ctx.FatalIfErrorf(err, "processing data: %s", err)
}

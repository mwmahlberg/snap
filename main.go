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
	"flag"
	"fmt"
	"os"

	"github.com/mwmahlberg/snap/internal"
	dbg "github.com/tj/go-debug"
)

var (
	unsnap      bool
	inFileName  string
	outFileName string
	file        *os.File
)

func init() {
	debug := dbg.Debug("INIT")
	flag.BoolVar(&unsnap, "u", false, "unsnap source file")
	flag.Parse()

	if inFileName = flag.Arg(0); inFileName == "" {
		fmt.Println("no input file given")
	}

	if unsnap {
		in := []rune(inFileName)
		outFileName = string(in[:len(in)-3])
		debug("Name: %s", outFileName)
	} else {
		outFileName = inFileName + `.sz`
	}
}

func main() {

	s, err := internal.NewSnapper(inFileName, outFileName)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	if unsnap {
		if err := s.Unsnap(); err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}

		return
	}

	if err := s.Snap(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

}

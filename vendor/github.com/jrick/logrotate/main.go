// Copyright (c) 2017, moshee
// Copyright (c) 2017, Josh Rickmar <jrick@devio.us>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Command logrotate writes and rotates logs read from stdin.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jrick/logrotate/rotator"
)

var (
	flagT = flag.Bool("t", false, "Behave like tee(1)")
	flagC = flag.Int("c", 5000, "Max (uncompressed) logfile size in kB")
	flagR = flag.Int("r", 0, "Max number of roll files to keep, 0 is unlimited")
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: <process that outputs to stdout> | logrotate [-t] [-c <N>] <filename>")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	r, err := rotator.New(flag.Arg(0), int64(*flagC), *flagT, *flagR)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	if err := r.Run(os.Stdin); err != nil {
		log.Print(err)
		return
	}
}

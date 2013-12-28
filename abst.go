/*
 GoABM - Agent Based Modeling

Copyright 2013 by Remo Hertig <remo.hertig@bluewin.ch>

GoABM is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

// abst.go is the interface to the abst webinterface
package goabm

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

var JournaledSimulation bool
var LogToFile bool
var OutputDir string
var RunID string

// add some flags
func Init() {
	flag.BoolVar(&JournaledSimulation, "abst.journal", false, "log all simulation states (agent moves)")
	flag.BoolVar(&LogToFile, "abst.logtofile", false, "log aggregated states to file in abst.out")
	flag.StringVar(&OutputDir, "abst.out", "out", "output dir")
	flag.StringVar(&RunID, "abst.runid", "", "id of the run, random if not provided")
}

func GetAbstPath() {

}

type Abst struct {
	Log     *os.File
	Journal *os.File
}

func (a *Abst) Init() {
	_, err := os.Stat(OutputDir)
	if err != nil {
		// create directory
		err = os.Mkdir(OutputDir, 0700)
		if err != nil {
			panic(err)
		}
	}
	if RunID == "" {
		RunID = fmt.Sprintf("%d", rand.Int())
	}
	// create output streams
	if LogToFile {
		f, err := os.Create(OutputDir + "/goabm-" + RunID + "-log")
		if err != nil {
			panic(err)
		}
		a.Log = f

	} else {
		// just use stdout
		a.Log = os.Stdout

	}
	if JournaledSimulation {
		// create journal file
		f, err := os.Create(OutputDir + "/goabm-" + RunID + "-journal")
		if err != nil {
			panic(err)
		}
		a.Journal = f
	}
}

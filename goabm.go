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

// GoABM is a Agent Based Modeling library
// currently there is only a static 2d grid landscape, see examples/
package goabm

import ("math/rand"
 "reflect"
 "fmt"
"os"
)

type AgentID int

type Agenter interface {
	Act()
}

type Modeler interface {
	LandscapeAction()
	Init(interface{})//Landscaper)
	CreateAgent(interface{}) Agenter
}

type Landscaper interface {
	Init(Modeler)
	//Action()
	GetAgents() []Agenter
	Dump() []byte //TODO: cleanup dump and use streams
	GetAgentById(AgentID) Agenter
}



type Statistics struct {
	Events int
	Steps int
}

type Simulation struct {
	Features  int
	Traits    int
	Landscape Landscaper
	Stats     Statistics
	Model     Modeler
	Log Logger
	AbstInterface Abst
}

func (s *Simulation) Init() {

	s.Landscape.Init(s.Model)
	s.Model.Init(s.Landscape)
	s.Log.Model = s.Model


	s.AbstInterface.Init()
		s.Log.Out = s.AbstInterface.Log
	s.Log.Init()


}
func (s *Simulation) Step() {
	s.Model.LandscapeAction()
	//r := rand.New(rand.NewSource(90))
	order := rand.Perm(len(s.Landscape.GetAgents()))
	//fmt.Printf("order: %v\n", order)
	for _, i := range order {
		s.Landscape.GetAgents()[i].Act()
		//fmt.Printf("running Agent #%d\n",i)
		s.Stats.Events = s.Stats.Events + 1
	}
	s.Stats.Steps = s.Stats.Steps + 1
	s.Log.Step(s.Stats)
	
	if(JournaledSimulation) { // dump landscape
	
	//fmt.Println(s.Landscape.Dump())
	s.AbstInterface.Journal.Write(s.Landscape.Dump())
	s.AbstInterface.Journal.Write([]byte("\n"))
	}

}

type Logger struct {
	StdOut bool
	Model Modeler
	FirstOut bool
	Out *os.File

}

func (l *Logger) Init() {
l.FirstOut = true
	if(l.StdOut) {
	fmt.Fprintf(l.Out,"Step,\tEvents,\t")
	}

}

func (l *Logger) Step(stats Statistics) {
	if l.StdOut {
		// get the fields of the model through reflection
		s := reflect.ValueOf(l.Model).Elem()
		//s := reflect.Indirect(in).Elem()
		typeOfT := s.Type()
		if !l.FirstOut {
			fmt.Fprintf(l.Out,"%d,\t%d,\t", stats.Steps, stats.Events)
		}
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			if f.Type().Kind() != reflect.Interface {
				if typeOfT.Field(i).Tag.Get("goabm") != "hide" {
					if l.FirstOut {
						fmt.Fprintf(l.Out,"%s,\t", typeOfT.Field(i).Name)
					} else {
						fmt.Fprintf(l.Out,"%v,\t",f.Interface())
					}
				}

			}
		}
		fmt.Fprintf(l.Out,"\n")
		if l.FirstOut {
			l.FirstOut = false
		}
	}

}



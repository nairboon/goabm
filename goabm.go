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

import "math/rand"


type Agenter interface {
	Act()
}

type Modeler interface {
	LandscapeAction()
	Init(Landscaper)
	CreateAgent(interface{}) Agenter
}

type Landscaper interface {
	Init(Modeler)
	//Action()
	GetAgents() []Agenter
}



type Statistics struct {
	Events int
}

type Simulation struct {
	Features  int
	Traits    int
	Landscape Landscaper
	Stats     Statistics
	Model     Modeler
}

func (s *Simulation) Init() {

	s.Landscape.Init(s.Model)
	s.Model.Init(s.Landscape)

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

}

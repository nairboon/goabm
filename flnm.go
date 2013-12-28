/*
This file is part of GoABM
Copyright 2013 by Remo Hertig <remo.hertig@bluewin.ch>
*/
package goabm

import "fmt"
import "math/rand"

// 2d landscape with no movement
type FixedLandscapeNoMovement struct {
	Agents     []FLNMAgent // library agent object, implements neighbor selection etc.
	UserAgents []Agenter   // agents from the user
	Size       int
	width      int
	height     int
}

type FLNMAgenter interface {
	GetRandomNeighbor() Agenter
}

type FLNMAgent struct {
	seqnr int
	x     int
	y     int
	ls    *FixedLandscapeNoMovement
	//exe Agenter
}

func (l *FixedLandscapeNoMovement) GetAgents() []Agenter {

	return l.UserAgents
}

func (l *FixedLandscapeNoMovement) GetAgent(x, y int) Agenter {
	//fmt.Printf("accessing %d/%d\n", x,y)
	// outerbounds enter on the opposite side
	if x >= l.width || x < 0 {
		x = 0
	}
	if y >= l.height || y < 0 {
		y = 0
	}

	return l.UserAgents[l.width*x+y]
}

func (a *FLNMAgent) GetRandomNeighbor() Agenter {
	switch choice := rand.Int31n(3); choice {
	case 0: // top
		return a.ls.GetAgent(a.x, a.y+1)
	case 1: // right
		return a.ls.GetAgent(a.x+1, a.y)
	case 2: // down
		return a.ls.GetAgent(a.x, a.y-1)
	case 3: // left
		return a.ls.GetAgent(a.x-1, a.y)
	default:
		panic(">3")

	}
}

func (l *FixedLandscapeNoMovement) Init(model Modeler) {
	numAgents := l.Size * l.Size
	fmt.Printf("Init landscape with %d agents\n", numAgents)

	l.width = l.Size
	l.height = l.Size

	l.Agents = make([]FLNMAgent, numAgents)
	l.UserAgents = make([]Agenter, numAgents)
	y := 0
	x := 0
	for i := range l.Agents {
		//for i:=0;i<numAgents;i++ {
		l.UserAgents[i] = model.CreateAgent(&l.Agents[i])

		l.Agents[i].seqnr = i

		l.Agents[i].x = x
		l.Agents[i].y = y
		x += 1
		if x >= l.width {
			// new row
			x = 0
			y += 1

		}
		l.Agents[i].ls = l

	}
}
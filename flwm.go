/*
This file is part of GoABM
Copyright 2013 by Remo Hertig <remo.hertig@bluewin.ch>
*/
package goabm

import "fmt"
import "math/rand"
import qt "github.com/larspensjo/quadtree" 
import vector "github.com/proxypoke/vector" 

// 2d landscape with movement
type FixedLandscapeWithMovement struct {
	Agents     []FLWMAgent // library agent object, implements neighbor selection etc.
	UserAgents []Agenter   // agents from the user
	Size       int
	Sight	   float64
	NAgents    int
	width      int
	height     int
	tree	   *qt.Quadtree
}

type FLWMAgenter interface {
	GetRandomNeighbor() Agenter
	MoveRandomly(float64)
}

type FLWMAgent struct {
	seqnr int
	x     float64
	y     float64
	ls    *FixedLandscapeWithMovement
	qt.Handle
	//exe Agenter
}

func (l *FixedLandscapeWithMovement) Dump() string {

	return ""
}

func (l *FixedLandscapeWithMovement) GetAgents() []Agenter {

	return l.UserAgents
}


func (a *FLWMAgent) MoveRandomly(steplength float64) {
	//random direction
	v:= vector.NewFrom([]float64{rand.Float64(),rand.Float64()})
	v.Normalize()
	v.Scale(steplength)
	x,_ := v.Get(0)
	y,_ := v.Get(1)
	a.ls.tree.Move(a,qt.Twof{x,y})
}

func (a *FLWMAgent) GetRandomNeighbor() Agenter {
	possibleNeighbors := a.ls.tree.FindNearObjects(qt.Twof{a.x, a.y},a.ls.Sight)
	if len(possibleNeighbors) < 1 {
	  return nil
}
	choice := rand.Int31n(int32(len(possibleNeighbors)))
	i:= possibleNeighbors[choice].(*FLWMAgent).seqnr
	return a.ls.UserAgents[i]
}

func (l *FixedLandscapeWithMovement) Init(model Modeler) {
	numAgents := l.Size * l.Size
	fmt.Printf("Init landscape with %d agents\n", numAgents)

	l.width = l.Size
	l.height = l.Size
	
	l.tree = qt.MakeQuadtree(qt.Twof{0, 0}, qt.Twof{float64(l.Size), float64(l.Size)})
	
	l.Agents = make([]FLWMAgent, numAgents)
	l.UserAgents = make([]Agenter, numAgents)
	y := 0
	x := 0
	for i := range l.Agents {
		//for i:=0;i<numAgents;i++ {
		l.UserAgents[i] = model.CreateAgent(&l.Agents[i])

		l.Agents[i].seqnr = i
		l.Agents[i].x = float64(x)
		l.Agents[i].y = float64(y)

		l.tree.Add(&l.Agents[i],qt.Twof{float64(x), float64(y)})

		x += 1
		if x >= l.width {
			// new row
			x = 0
			y += 1

		}
		l.Agents[i].ls = l

	}
}

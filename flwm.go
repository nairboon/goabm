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
	Sight      float64
	NAgents    int
	tree       *qt.Quadtree
}

type FLWMAgenter interface {
	GetRandomNeighbor() Agenter
	MoveRandomly(float64)
	Id() AgentID
}

type FLWMAgent struct {
	Seqnr AgentID `json:"index"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	ls    *FixedLandscapeWithMovement `json:"-"`
	qt.Handle `json:"-"`
	//exe Agenter
}

func (a *FLWMAgent) Id() AgentID {
	return a.Seqnr
}

func (l *FixedLandscapeWithMovement) Dump() NetworkDump {
	// dump as a network
	nodes := l.UserAgents
	var links []Link

	for _, a := range l.Agents {

		tmp := a.ls.tree.FindNearObjects(qt.Twof{a.X, a.Y}, a.ls.Sight)

		for _, v := range tmp {
			if v.(*FLWMAgent).Seqnr != a.Seqnr {
				//panic("self link")
				link := Link{Source: a.Seqnr, Target: v.(*FLWMAgent).Seqnr}
				links = append(links, link)

			}
		}

	}

	
	return NetworkDump{Nodes:nodes,Links:links}
}

func (l *FixedLandscapeWithMovement) GetAgents() *[]Agenter {

	return &l.UserAgents
}

func (l *FixedLandscapeWithMovement) GetAgentById(id AgentID) Agenter {
	for i, a := range l.Agents {
		if a.Seqnr == id {
			return l.UserAgents[i]
		}
	}
	return nil
}
func random(min, max float64) float64 {
  return rand.Float64() * (max - min) + min
}

func (a *FLWMAgent) MoveRandomly(steplength float64) {
	//random direction
	bsize := float64(a.ls.Size)
	v := vector.NewFrom([]float64{random(-bsize,bsize), random(-bsize,bsize)})
	v.Normalize()
	v.Scale(steplength)
	x, _ := v.Get(0) 
	y, _ := v.Get(1)
	x = x + a.X
	y = y + a.Y
	
	// check bounds
	if x < 0 {
	 x = bsize + x // reenter world on the other side
	}
	if y < 0 {
	 y = bsize + y // reenter world on the other side
	}
	if x > bsize {
	x = x - bsize
	}
	if y > bsize {
	y = y - bsize
	}
	if x >= float64(a.ls.Size) || x < 0 || y >= float64(a.ls.Size) || y < 0  {
        // just try again
        fmt.Printf("out of bounds %f/%f of %f\n",x,y, bsize)
         a.MoveRandomly(steplength)
         return
	}
	//fmt.Printf("move from %f/%f to %f/%f\n",a.X,a.Y,x,y)
	a.ls.tree.Move(a, qt.Twof{x, y})

	a.X = x
	a.Y = y
}

func (a *FLWMAgent) GetRandomNeighbor() Agenter {
	tmp := a.ls.tree.FindNearObjects(qt.Twof{a.X, a.Y}, a.ls.Sight)
	var possibleNeighbors []qt.Object
	for _, v := range tmp {
		if v.(*FLWMAgent).Seqnr != a.Seqnr {
			possibleNeighbors = append(possibleNeighbors, v)
		}
	}
	if len(possibleNeighbors) < 1 {
		return nil
	}

	choice := rand.Int31n(int32(len(possibleNeighbors)))
	i := possibleNeighbors[choice].(*FLWMAgent).Seqnr
	if i == a.Seqnr {
		panic("same agent")
	}
	return a.ls.UserAgents[i]
}

func (l *FixedLandscapeWithMovement) Init(model Modeler) {
	numAgents := l.NAgents
	//fmt.Printf("Init landscape with %d agents\n", numAgents)



	l.tree = qt.MakeQuadtree(qt.Twof{0, 0}, qt.Twof{float64(l.Size), float64(l.Size)})

	l.Agents = make([]FLWMAgent, numAgents)
	l.UserAgents = make([]Agenter, numAgents)
	y := 0
	x := 0
	for i := range l.Agents {
		//for i:=0;i<numAgents;i++ {
		l.UserAgents[i] = model.CreateAgent(&l.Agents[i])

		l.Agents[i].Seqnr = AgentID(i)
		l.Agents[i].X = float64(x)
		l.Agents[i].Y = float64(y)

		l.tree.Add(&l.Agents[i], qt.Twof{float64(x), float64(y)})

		x += 1
		if x >= l.Size {
			// new row
			x = 0
			y += 1

		}
		l.Agents[i].ls = l

	}
}

/*
This file is part of GoABM
Copyright 2013 by Remo Hertig <remo.hertig@bluewin.ch>
*/

// An go implementation of Robert Axelrods ABM model of disseminating culture [1]
// [1]: Axelrod, Robert. "The dissemination of culture a model with local convergence and global polarization." Journal of conflict resolution 41, no. 2 (1997): 203-226.
package main

import "fmt"
import "math/rand"
import "goabm"
import "flag"

// Implementation of the Agent, cultural Traits are stored in features
type AxelrodAgent struct {
	features Feature
	Agent    goabm.FLNMAgenter
}

// returns the culture as a string
func (a *AxelrodAgent) Culture() string {
	return fmt.Sprintf("%v", a.features)
}

// required for the simulation interface, called everytime when the agent is activated
func (a *AxelrodAgent) Act() {

	//fmt.Printf("Agent culture: %v\n",a.features)
	other := a.Agent.GetRandomNeighbor().(*AxelrodAgent)
	sim := a.Similarity(other)
	if sim >= 0.99 {
		// agents are already equal
		return
	}
	probabilityToInteract := rand.Float32()
	//fmt.Printf("interacting %f <= %f\n",probabilityToInteract,sim)
	//interact with sim% chance
	if probabilityToInteract <= sim {

		for i := range a.features {
			if a.features[i] != other.features[i] {
				//fmt.Printf("%d influenced %d\n", other.seqnr, a.seqnr)
				a.features[i] = other.features[i]
				return
			}

		}
	}

}

// helper function to determine the similarity between to agents
func (a *AxelrodAgent) Similarity(other *AxelrodAgent) float32 {
	c := float32(0.0)
	// count equal traits, final score = shared traits/total traits
	for i := range a.features {
		if a.features[i] == other.features[i] {
			c = c + 1
		}
	}
	//fmt.Printf("sim: %f/%d\n",c,len(a.features))
	return c / float32(len(a.features))
}

type Feature []int
// implementation of the model
type Axelrod struct {
	Cultures  int
	Landscape goabm.Landscaper
	Traits    int
	Features  int
}

func (a *Axelrod) Init(l goabm.Landscaper) {
	a.Landscape = l
}

func (a *Axelrod) CreateAgent(agenter interface{}) goabm.Agenter {

	agent := &AxelrodAgent{Agent: agenter.(goabm.FLNMAgenter)}

	f := make(Feature, a.Features)
	for i := range f {
		f[i] = rand.Intn(a.Traits)
	}
	agent.features = f
	return agent
}

func (a *Axelrod) LandscapeAction() {
	a.Cultures = a.CountCultures()

}

func (a *Axelrod) CountCultures() int {
	cultures := make(map[string]int)
	for _, b := range a.Landscape.GetAgents() {
		a := b.(*AxelrodAgent)
		cul := a.Culture()
		if _, ok := cultures[cul]; ok {
			cultures[cul] = 1
		} else {
			cultures[cul] = cultures[cul] + 1
		}
	}
	return len(cultures)
}

func main() {
	fmt.Println("ABM simulation")

	var traits = flag.Int("traits", 5, "help message for flagname")
	var features = flag.Int("features", 5, "help message for flagname")
	var size = flag.Int("size", 10, "help message for flagname")

	var runs = flag.Int("runs", 200, "help message for flagname")
	flag.Parse()

	model := &Axelrod{Traits: *traits, Features: *features}
	sim := &goabm.Simulation{Landscape: &goabm.FixedLandscapeNoMovement{Size: *size}, Model: model}
	sim.Init()
	for i := 0; i < *runs; i++ {
		fmt.Printf("Step #%d, Events:%d, Cultures:%d\n", i, sim.Stats.Events, model.Cultures)
		if model.Cultures == 1 {
			return
		}
		sim.Step()

	}
	//fmt.Printf("%v\n",sim.Landscape.GetAgents())

}

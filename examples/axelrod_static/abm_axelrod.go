/*
This file is part of GoABM
Copyright 2013 by Remo Hertig <remo.hertig@bluewin.ch>
*/

// An go implementation of Robert Axelrods ABM model of disseminating culture [1].
// [1]: Axelrod, Robert. "The dissemination of culture a model with local convergence and global polarization." Journal of conflict resolution 41, no. 2 (1997): 203-226.
package main

import "fmt"
import "math/rand"
import "goabm"
import "flag"

// Implementation of the Agent, cultural Traits are stored in features
type AxelrodAgent struct {
	Features Feature
	*goabm.FLNMAgent
	Model *Axelrod
}

// returns the culture as a string
func (a *AxelrodAgent) Culture() string {
	return fmt.Sprintf("%v", a.Features)
}
/*
func (a *AxelrodAgent) ID() string {
}*/

// required for the simulation interface, called everytime when the agent is activated
func (a *AxelrodAgent) Act() {

	// (ii) (a) selects a neighbor for cultural interaction
	other := a.Model.GetRandomNeighbor(a)
	sim := a.Similarity(other)
	if sim >= 0.99 {
		// agents are already equal
		return
	}
	dice := rand.Float32()
	//interact with sim% chance
	if dice <= sim {

		for i := range a.Features {
			if a.Features[i] != other.Features[i] {
				//fmt.Printf("%d influenced %d\n", other.seqnr, a.seqnr)
				a.Features[i] = other.Features[i]
				return
			}

		}
	}

}

// helper function to determine the similarity between to agents
func (a *AxelrodAgent) Similarity(other *AxelrodAgent) float32 {
	c := float32(0.0)
	// count equal traits, final score = shared traits/total traits
	for i := range a.Features {
		if a.Features[i] == other.Features[i] {
			c = c + 1
		}
	}
	//fmt.Printf("sim: %f/%d\n",c,len(a.features))
	return c / float32(len(a.Features))
}

type Feature []int
// implementation of the model
type Axelrod struct {
	Cultures  int
	Landscape goabm.Landscaper
	Traits    int `goabm:"hide"` // don't show these in the stats'
	Features  int `goabm:"hide"`
}

func (a *Axelrod) GetRandomNeighbor(origin *AxelrodAgent) *AxelrodAgent{
       id, err := origin.GetRandomLink()
       if(err != nil) {
       panic(err)
       }
       return a.Landscape.GetAgentById(id).(*AxelrodAgent)
}

func (a *Axelrod) Init(l interface{}) {
	a.Landscape = l.(goabm.Landscaper)
}

func (a *Axelrod) CreateAgent(agenter interface{}) goabm.Agenter {

	agent := &AxelrodAgent{FLNMAgent: agenter.(*goabm.FLNMAgent)}

	f := make(Feature, a.Features)
	for i := range f {
		f[i] = rand.Intn(a.Traits)
	}
	agent.Features = f
	agent.Model = a
	return agent
}

func (a *Axelrod) LandscapeAction() {
	a.Cultures = a.CountCultures()

}

func (a *Axelrod) CountCultures() int {
	cultures := make(map[string]int)
	for _, b := range *a.Landscape.GetAgents() {
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
        //initialize the goabm library (logs & flags)
	goabm.Init()

	var traits = flag.Int("traits", 15, "number of cultural traits per feature")
	var features = flag.Int("features", 15, "number of cultural features")
	var size = flag.Int("size", 10, "size (width/height) of the landscape")

	var runs = flag.Int("runs", 200, "number of simulation runs")
	flag.Parse()
		fmt.Println("ABM simulation")
        // create your model which has to satisfy the Modeler interface
	model := &Axelrod{Traits: *traits, Features: *features}
	// create the simulation with a Landscape, your model and a logger
	sim := &goabm.Simulation{Landscape: &goabm.FixedLandscapeNoMovement{Size: *size}, Model: model, Log: goabm.Logger{StdOut: true}}
	sim.Init()
	for i := 0; i < *runs; i++ {
		//fmt.Printf("Step #%d, Events:%d, Cultures:%d\n", i, sim.Stats.Events, model.Cultures)
		if model.Cultures == 1 {
			sim.Stop()
	fmt.Printf("Stimulation prematurely done\n")
			return
		}
		sim.Step()

	}
	sim.Stop()
	fmt.Printf("Stimulation done\n")

}

/*
This file is part of GoABM
Copyright 2013 by Remo Hertig <remo.hertig@bluewin.ch>
*/

// An extension to Robert Axelrods ABM model of disseminating culture [1]
// [1]: Axelrod, Robert. "The dissemination of culture a model with local convergence and global polarization." Journal of conflict resolution 41, no. 2 (1997): 203-226.
// based on the work of Arezky Hernandez Rodriguez
// (Overview , Design Concepts , and Details ( ODD ) for the Axelrod ’ s model for Cultural Dissemination)
//
// In this model agents will move around the landscape
package main

import "fmt"
import "math/rand"
import "goabm"
import "flag"

// Implementation of the Agent, cultural Traits are stored in Features
type AxelrodAgent struct {
	Features Feature
	Agent    goabm.FLWMAgenter
	Steplength float64  `goabm:"hide"`
	ProbVeloc float64  `goabm:"hide"`
}

// returns the culture as a string
func (a *AxelrodAgent) Culture() string {
	return fmt.Sprintf("%v", a.Features)
}

// required for the simulation interface, called everytime when the agent is activated
func (a *AxelrodAgent) Act() {
	//fmt.Printf("Agent culture: %v\n",a.Features)


	dice := rand.Float64()
	// (i) agent decides to move according to the probability veloc
	if dice <= a.ProbVeloc {
		a.Agent.MoveRandomly(a.Steplength)
	}

	// (ii) (a) selects a neighbor for cultural interaction
	res := a.Agent.GetRandomNeighbor()
	// there is no agent in sight
	if res == nil {
	return
	}
	other :=res.(*AxelrodAgent)
	sim := a.Similarity(other)
	if sim >= 0.99 {
		// agents are already equal
		return
	}
	dice2 := rand.Float32()
	//fmt.Printf("interacting %f <= %f\n",probabilityToInteract,sim)
	//interact with sim% chance
	if dice2 <= sim {

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
	//fmt.Printf("sim: %f/%d\n",c,len(a.Features))
	return c / float32(len(a.Features))
}

type Feature []int
// implementation of the model
type Axelrod struct {
	Cultures  int
	Landscape goabm.Landscaper
	Traits    int  `goabm:"hide"`
	Features  int  `goabm:"hide"`
	Steplength float64  `goabm:"hide"`
	ProbVeloc float64  `goabm:"hide"`
}

func (a *Axelrod) Init(l interface{}) {
	a.Landscape = l.(goabm.Landscaper)
}

func (a *Axelrod) CreateAgent(agenter interface{}) goabm.Agenter {

	agent := &AxelrodAgent{Agent: agenter.(goabm.FLWMAgenter)}

	f := make(Feature, a.Features)
	for i := range f {
		f[i] = rand.Intn(a.Traits)
	}
	agent.Features = f
	agent.ProbVeloc = a.ProbVeloc
	agent.Steplength = a.Steplength
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
       //initialize the goabm library (logs & flags)
	goabm.Init()
	        
	fmt.Println("ABM simulation")
	// general model parameters
	var traits = flag.Int("traits", 5, "number of cultural traits per feature")
	var Features = flag.Int("Features", 5, "number of cultural Features")
	var size = flag.Int("size", 10, "size (width/height) of the landscape")

	// parameters for the moving model
	var probveloc = flag.Float64("pveloc", 0.05, "probability that an agent moves")
	var steplength = flag.Float64("steplength", 0.1, "maximal distance a agent can travel per step")
	var sight = flag.Float64("sight", 1, "radius in which agent can interact")
	var numAgents = flag.Int("agents", 100, "number of agents to simulate")

	var runs = flag.Int("runs", 200, "number of simulation runs")
	flag.Parse()

	model := &Axelrod{Traits: *traits, Features: *Features, ProbVeloc: *probveloc, Steplength: *steplength}
sim := &goabm.Simulation{Landscape: &goabm.FixedLandscapeWithMovement{Size: *size, NAgents: *numAgents,Sight:*sight},
 Model: model , Log: goabm.Logger{StdOut: true}}
	sim.Init()
	for i := 0; i < *runs; i++ {
		//fmt.Printf("Step #%d, Events:%d, Cultures:%d\n", i, sim.Stats.Events, model.Cultures)
		if model.Cultures == 1 {
			return
		}
		sim.Step()

	}
	//fmt.Printf("%v\n",sim.Landscape.GetAgents())

}

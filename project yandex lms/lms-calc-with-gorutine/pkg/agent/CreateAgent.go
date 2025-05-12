package agent

import (
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"sync"
)

type Agent struct {
	Computing_power int
	Wg              sync.WaitGroup
	chanTasks       chan entites.Task
	mutex           sync.Mutex
	chanResults     chan float64
	chanErrors      chan error
}

func NewAgent(computing_power int) *Agent {
	return &Agent{
		Computing_power: computing_power,
		Wg:              sync.WaitGroup{},
		mutex:           sync.Mutex{},
		chanTasks:       make(chan entites.Task),
		chanResults:     make(chan float64),
		chanErrors:      make(chan error),
	}
}

package machine

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-pg/pg"
	"github.com/m49n3t0/boretto/models"
	"github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////

// to define the dispatcher parameters
type DispatcherParams struct {
	// robot configuration
    Function string
	// dispatcher configuration
	MaxWorker int64
    MaxQueue int64
    // machine interface object
    Machine *Machine
}

// dispatcher object
type Dispatcher struct {

	// function on what the robot work
	function string

	// store datas from database
	definitions map[string]*models.Definition
	endpoints   map[string]*models.Endpoint

	// manage the distribution of workflow
	workerPool chan chan string
	queue      chan string

	// manage the quit process
	signal chan os.Signal
	quit   chan bool

	// database handler
	db *pg.DB

    // logger handler
    log *logrus.Entry
}

// dispatcher creation handler
func NewDispatcher(params *DispatcherParams) (*Dispatcher, error) {

    machine := params.Machine
	definitions := make(map[string]*models.Definition)
	endpoints := make(map[string]*models.Endpoint)
	workerPool := make(chan chan string, params.MaxWorker)
	queue := make(chan string, params.MaxQueue)
	signal := make(chan os.Signal, 2)
	quit := make(chan bool)

    // get the logger from interface
    logger, err := machine.GetLogger()
    if err != nil {
        return nil, err
    }
    if logger == nil {
        return nil, errors.New("the logger initialization return an empty object")
    }

    // get the database handler from interface
    database, err := machine.GetDatabase(logger)
    if err != nil {
        return nil, err
    }

    // create the object
	dispatcher := &Dispatcher{
		function:    params.Function,
		definitions: definitions,
		endpoints:   endpoints,
		workerPool:  workerPool,
		queue:       queue,
		signal:      signal,
		quit:        quit,
        db: database,
        log: logger,
	}

	return dispatcher, nil
}

// stop signals programmatically
func (dispatcher *Dispatcher) Stop() {
	go func() {
		dispatcher.log.Println("STOP")
		dispatcher.quit <- true
	}()
}

// stop signals from system
func (dispatcher *Dispatcher) Signal() {
	go func() {

        // link system signal to the dispatcher signal
        signal.Notify(dispatcher.signal, os.Interrupt, syscall.SIGTERM)

        // when receive syscall signal
		<-dispatcher.signal

        // do a stopper
		dispatcher.log.Println("SIGNAL")
		dispatcher.Stop()
	}()
}

///////////////////////////////////////////////////////////////////////////////















































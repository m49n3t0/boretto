package machine

import (
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-pg/pg"
	"github.com/m49n3t0/boretto/models"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////

// interface to implements to correctly works with this executor
type Machine interface {
	// permit to retrieve logger handler
	GetLogger() (*logrus.Entry, error)
	// permit to retrieve database handler
	GetDatabase(*logrus.Entry) (*pg.DB, error)
}

///////////////////////////////////////////////////////////////////////////////

// default machine definition
type DefaultMachine struct{}

// permit to retrieve logger handler
func (machine *DefaultMachine) GetLogger() (*logrus.Entry, error) {

	// create a default logger
	var logger = logrus.New()

	// define the logger default level
	log.Level = logrus.DebugLevel

	// define the logger default output
	log.Out = os.Stdout

	return logger, nil
}

// permit to retrieve database handler
func (machine *DefaultMachine) GetDatabase(logger *logrus.Entry) (*pg.DB, error) {

	// database configuration
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	// build database host address
	var address = host
	if port != "" {
		address = host + ":" + port
	}

	// pg database connector
	db := pg.Connect(&pg.Options{
		Addr:       address,
		User:       user,
		Password:   password,
		Database:   database,
		MaxRetries: 2,
	})

	// check connection
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		logger.Println("Problem while check database connection")
		return nil, err
	}

	return db, nil
}

///////////////////////////////////////////////////////////////////////////////

// to launch the machine executor with a default struct
func RunDefault(function string) {
	return Run(DefaultMachine{}, function)
}

// to launch the machine executor
func Run(machine *Machine, function string) {

	// check mandatory parameter
	if machine != nil {
		logrus.Panic("Missing mandatory machine parameter")
	}

	// create a dispatcher
	dispatcher, err := NewDispatcher(&DispatcherParams{
		Machine:   machine,
		Function:  function,
		MaxWorker: 20,
		MaxQueue:  5,
	})

	// deferred the database connection closes
	defer dispatcher.DB.Close()

	// dispatcher creation catch error
	if err != nil {
		logrus.Panic(err)
	}

	// dispatcher logger
	log := dispatcher.Logger

	// get robot configuration
	err = dispatcher.GetRobotConfiguration()
	if err != nil {
		log.Panic(err)
	}

	// listen the channel
	go dispatcher.Signal()

	// launch a first task ID listing
	go dispatcher.GetTasks()

	//	// launch the database NOTIFY listener
	//	go dispatcher.Listen()

	// starting n number of workers
	for i := int64(0); i < dispatcher.MaxWorker; i++ {

		// create a new worker
		worker := NewWorker(&i, dispatcher)

		// start it
		worker.Start()
	}

	// launch the dispatch
	log.Info("Worker dispatch started...")

	for {
		select {
		case ID := <-dispatcher.Queue:

			log.Info("Dispatch to taskChannel with ID: %s", ID)

			// try to obtain a worker task channel that is available.
			// this will block until a worker is idle
			taskChannel := <-dispatcher.WorkerPool

			// dispatch the task to the worker task channel
			taskChannel <- ID

		case <-dispatcher.Quit:

			// we have received a signal to stop
			log.Info("Dispatch is stopping")

			// XXX : how to stop workers correctly

			return
		}
	}
}

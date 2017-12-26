package bot

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-pg/pg"

	"github.com/m49n3t0/boretto/models"
)

///////////////////////////////////////////////////////////////////////////////

// dispatcher object
type Dispatcher struct {

	// function on what the robot work
	function        string

	// store datas from database
	definitions     map[int64]*models.Definition
	endpoints       map[int64]*models.Endpoint

	// manage the distribution of workflow
	workerPool      chan chan int64
	queue           chan int64

	// manage the quit process
	signal          chan os.Signal
	quit            chan bool

	// database handler
	db *pg.DB
}

// dispatcher creation handler
func New() (*Dispatcher, error) {

	definitions := make(map[int64]*models.Definition)
	endpoints   := make(map[int64]*models.Endpoint)
	workerPool  := make(chan chan int64, MAX_WORKER)
	queue       := make(chan int64, MAX_QUEUE)
	signal      := make(chan os.Signal, 2)
	quit        := make(chan bool)

	dispatcher := &Dispatcher{
		function:       FUNCTION,
        definitions:    definitions,
		endpoints:      endpoints,
		workerPool:     workerPool,
		queue:          queue,
		signal:         signal,
		quit:           quit,
	}

	return dispatcher, nil
}

// do the dispatching processes
func (dispatcher *Dispatcher) Run() {

	// link system signal to the dispatcher signal
	signal.Notify(dispatcher.signal, os.Interrupt, syscall.SIGTERM)

	// listen the channel
	go dispatcher.Signal()

	// get database connection
	err := dispatcher.dbConnect()
	if err != nil {
		panic(err)
	}

	// defer the disconnection
	defer dispatcher.dbDisconnect()

	// get robot configuration
	err = dispatcher.getConfiguration()
	if err != nil {
		panic(err)
	}

	// starting n number of workers
	for i := 0; i < ENV_MAX_WORKER; i++ {
		worker := NewWorker(dispatcher)

		log.Println(worker)

		worker.Start()
	}

	// launch a first task ID listing
	go dispatcher.getTaskIDs()

	// launch the database NOTIFY listener
	go dispatcher.listen()

	// launch the dispatch
	dispatcher.launch()
}

// launch task in free workers
func (dispatcher *Dispatcher) launch() {

	log.Println("LAUNCHED")

	log.Println("Worker dispatch started...")

	for {
		select {
		case taskId := <-dispatcher.queue:

			log.Printf("Dispatch to taskChannel with ID : " + strconv.Itoa(int(taskId)))

			// try to obtain a worker task channel that is available.
			// this will block until a worker is idle
			taskChannel := <-dispatcher.workerPool

			// dispatch the task to the worker task channel
			taskChannel <- taskId

		case <-dispatcher.quit:

			// we have received a signal to stop
			log.Println("RECEIVE QUIT")

			// XXX : how to stop workers correctly

			os.Exit(1)
		}
	}
}

// stop signals programmatically
func (dispatcher *Dispatcher) Stop() {
	go func() {
		log.Println("STOP")
		dispatcher.quit <- true
	}()
}

// stop signals from system
func (dispatcher *Dispatcher) Signal() {
	go func() {
		<-dispatcher.signal
		log.Println("SIGNAL")
		dispatcher.Stop()
	}()
}

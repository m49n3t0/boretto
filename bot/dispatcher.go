package bot

import (
	"github.com/go-pg/pg"
	"github.com/m49n3t0/boretto/models"
	"log"
	"strconv"
    "syscall"
    "os"
    "os/signal"
)

///////////////////////////////////////////////////////////////////////////////

// Dispatcher object
type Dispatcher struct {
	//    //	Function string
	//    //	Version  int64
	//

	// function on what the robot work
	function string

	// store datas from database
	robots    map[int64]*models.Definition
	endpoints map[int64]interface{}

	// manage the distribution of workflow
	workerPool chan chan int64
	queue      chan int64

    // manage the quit process
    signal chan os.Signal
	quit          chan bool

	// database handler
	db *pg.DB

	//
	//    //	definition    models.Definition
	//    //	endpoint_http map[int64]models.EndpointHttp
	//    //	quit          chan bool
}

// dispatcher create handler
func New() (*Dispatcher, error) {
	//    // get the
	//    version, err := strconv.ParseInt(ENV_VERSION, 10, 64)
	//    if err != nil {
	//        return nil, err
	//    }

	robots := make(map[int64]*models.Definition)
	endpoints := make(map[int64]interface{})
	workerPool := make(chan chan int64, ENV_MAX_WORKER)
	queue := make(chan int64, ENV_MAX_QUEUE)
    signal := make(chan os.Signal, 2)
    quit := make(chan bool)

	dispatcher := &Dispatcher{
		function: ENV_FUNCTION,
		//		Function:   ENV_FUNCTION,
		//		Version:    version,
		robots:     robots,
		endpoints:  endpoints,
		workerPool: workerPool,
		queue:      queue,
        signal: signal,
        quit: quit,
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

	//	// launch a first read on database data task
	//	go dispatcher.readTaskIds()
	//
	//	// launch the listener for database events
	//	go dispatcher.initializeListenerAndListen()

	// launch the dispatch
	dispatcher.launch()
}

// Launch task in free workers
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



// Stop signals programmatically
func (dispatcher *Dispatcher) Stop() {
	go func() {
        log.Println("STOP")
		dispatcher.quit <- true
	}()
}

// Stop signals from system
func (dispatcher *Dispatcher) Signal() {
	go func() {
        <-dispatcher.signal
        log.Println("SIGNAL")
        dispatcher.Stop()
	}()
}






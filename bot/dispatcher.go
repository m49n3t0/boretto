package bot

import (
	"github.com/go-pg/pg"
	"github.com/m49n3t0/boretto/models"
	"log"
	//	"strconv"
)

///////////////////////////////////////////////////////////////////////////////

// Dispatcher object
type Dispatcher struct {
	//    //	Function string
	//    //	Version  int64
	//

	// function on what the robot work
	function string

	// manage the distribution of workflow
	workerPool chan chan int64
	queue      chan int64

	// store datas from database
	robots    map[int64]*models.Definition
	endpoints map[int64]interface{}

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

	workerPool := make(chan chan int64, ENV_MAX_WORKER)
	queue := make(chan int64, ENV_MAX_QUEUE)
	robots := make(map[int64]*models.Definition)
	endpoints := make(map[int64]interface{})

	dispatcher := &Dispatcher{
		function: ENV_FUNCTION,
		//		Function:   ENV_FUNCTION,
		//		Version:    version,
		workerPool: workerPool,
		queue:      queue,
		robots:     robots,
		endpoints:  endpoints,
	}

	return dispatcher, nil
}

// do the dispatching processes
func (dispatcher *Dispatcher) Run() {

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

	//	// starting n number of workers
	//	for i := 0; i < dispatcher.Configuration.MaxWorkers; i++ {
	//		worker := NewWorker(dispatcher)
	//		worker.Start()
	//	}
	//
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

	//	log.Println("Worker dispatch started...")
	//
	//	for {
	//		select {
	//		case taskId := <-dispatcher.queue:
	//
	//			log.Printf("Dispatch to taskChannel with ID : " + strconv.Itoa(int(taskId)))
	//
	//			// try to obtain a worker task channel that is available.
	//			// this will block until a worker is idle
	//			taskChannel := <-dispatcher.workerPool
	//
	//			// dispatch the task to the worker task channel
	//			taskChannel <- taskId
	//
	//		case <-dispatcher.quit:
	//			// we have received a signal to stop
	//
	//			// XXX : how to stop workers correctly
	//		}
	//	}
}

// XXX : how to improve this part ?
// Stop signals
func (dispatcher *Dispatcher) Stop() {
	//	go func() {
	//		dispatcher.quit <- true
	//	}()
}

package bot

import (
	"github.com/m49n3t0/boretto/models"
	"log"
	"strconv"
	"errors"
)

///////////////////////////////////////////////////////////////////////////////

// the worker executes the task process
type Worker struct {
	workerPool  chan chan int64
	taskChannel chan int64
	dispatcher  *Dispatcher
	quit        chan bool
}

// function to create a new worker
func NewWorker(dispatcher *Dispatcher) Worker {
	return Worker{
		workerPool:  dispatcher.workerPool,
		taskChannel: make(chan int64),
		dispatcher:  dispatcher,
		quit:        make(chan bool)}
}

// start method starts the run loop for the worker
// listening for a quit channel in case we need to stop it
func (worker *Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue
			worker.workerPool <- worker.taskChannel

			// read from the channel
			select {

                case taskId := <-worker.taskChannel:

                    log.Printf("Entry to taskChannel with ID : " + strconv.Itoa(int(taskId)) + "\n")

                    // we have received a work request
                    if err := worker.DoAction(taskId); err != nil {
                        log.Printf("Error while working on task: %s", err.Error())
                    }

                    log.Println(".", "ENDJOB")

                case <-worker.quit:

                    // we have received a signal to stop
                    return
			}
		}
	}()
}

// stop signals the worker to stop listening for work requests.
func (worker *Worker) Stop() {
	go func() {
		worker.quit <- true
	}()
}

///////////////////////////////////////////////////////////////////////////////

//// Response http from each API response
//type HttpResponse struct {
//	Buffer   *JsonB
//	Interval *int
//	Step     *string
//	Comment  *string
//}

func (worker *Worker) DoAction(id int64) error {

	// retrieve a task
	task, err := worker.dispatcher.getTask( id )

	if err != nil {
		log.Fatalln("Error while fetching one task", err)
	}

	// logger
	log.Println("Working on task " + strconv.Itoa(int(task.Id)) + "|" + task.Function + " on step: " + task.Step)

	// --------------------------------------------------------------------- //

	log.Println("42.0.1----------------------------------------------")

	// get robot for this task version
	robot, ok := worker.dispatcher.robots[ task.Version ]

	if !ok {
		log.Println("Robot definition for this version doesn't exists")
		return errors.New("Robot not found")
	}

	log.Println("42.0.2----------------------------------------------")

	log.Println( robot.Sequence )

	log.Println("42.0.3----------------------------------------------")

	// vars on step
	var step *models.Step

	// get the actual step in the sequence
	for _, s := range robot.Sequence {

        log.Println("42.0.4----------------------------------------------")
        log.Printf("task step : %s / step name : %s , %s , %s", task.Step, s.Name, s.EndpointType, s.EndpointID)

		// check with the local task
		if task.Step == s.Name {
            log.Println("42.0.5----FOUND")
			step = &s
			break
		}
	}

	log.Println("42.1.1----------------------------------------------")

	log.Println( step )

	log.Println("42.1.2----------------------------------------------")

	// check step found
	if step == nil {
		log.Println("Step not found for this version")
		return errors.New("Step not found")
	}

	log.Println("42.2----------------------------------------------")

	// get the associated endpoint
	endpoint, ok := worker.dispatcher.endpoints[ step.EndpointID ]

	if !ok {
		log.Println("Associated endpoint to this step doesn't exists")
		return errors.New("Endpoint not found")
	}

	log.Println("42.3----------------------------------------------")

	// --------------------------------------------------------------------- //

	// change the task data
	task.Status = "DOING"

	// update in database the task
	err = worker.dispatcher.updateTask( task )

	if err != nil {
		log.Println("Error while updating the task for status lock : ", err)
		return err
	}

	log.Println("42.4----------------------------------------------")

	// --------------------------------------------------------------------- //


	log.Println(endpoint)





	// --------------------------------------------------------------------- //

	log.Println("42.9----------------------------------------------")

	log.Println("Task updation")

	// update in database the task
	err = worker.dispatcher.updateTask( task )

	if err != nil {
		log.Println("Error while updating the task result : ", err)
		return err
	}

	log.Println("42.10----------------------------------------------")

	log.Println("Task updated")

	return nil
}




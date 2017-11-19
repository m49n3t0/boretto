package bot

import (
	//	"bytes"
	//	"encoding/json"
	//	_ "github.com/lib/pq"
	//	"gopkg.in/gorp.v2"
	//	"io/ioutil"
	//	"net/http"
	//	"time"
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


//	// vars on step data
//	var actualStep Step
//	var foundActualStep = false
//
//	// which is the actual step data
//	for _, s := range worker.Definition.Sequence {
//		if s.Name == task.Step {
//			actualStep = s
//			foundActualStep = true
//			break
//		}
//	}
//
//	// vars
//	var httpResponse HttpResponse
//	var statusCode int
//
//	// actual step process
//	if foundActualStep {
//
//		// do the http call to retrieve API data/informations
//		httpResponse, statusCode, err = worker.CallHttp(task, actualStep)
//
//		if err != nil {
//
//			// error while doing the http call on api's
//			var comment = "Error while doing the http call on api's"
//
//			log.Fatalln(comment)
//
//			// forge an error status code
//			statusCode = 500
//
//			// forge a fake http response
//			httpResponse = HttpResponse{
//				Buffer:  &task.Buffer, // repush the same buffer
//				Comment: &comment}     // forge a comment
//
//		}
//
//	} else {
//
//		// no step associated for this step name
//		var comment = "Error while check the actual step informations"
//
//		log.Fatalln(comment)
//
//		// forge an error status code
//		statusCode = 600
//
//		// forge a fake http response
//		httpResponse = HttpResponse{
//			Buffer:  &task.Buffer, // repush the same buffer
//			Comment: &comment}     // forge a comment
//
//	}
//
//	// default value
//	var statusName = "todo"
//
//	// switch on each status code
//	switch statusCode {
//
//        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//        //  200  = next                                                         //
//        //______________________________________________________________________//
//        case 200:
//
//            // not an ending step
//            if !actualStep.EndStep {
//
//                var found = false
//
//                // retrieve the next step data
//                for _, s := range worker.Definition.Sequence {
//                    if found {
//                        // set the next step
//                        task.Step = s.Name
//                        break
//                    }
//                    if s.Name == task.Step {
//                        found = true
//                    }
//                }
//
//                // no next step found, error
//                if !found {
//
//                    // error status due to no next step founded
//                    statusName = "error"
//
//                    // forge error comment
//                    var comment = "Impossible to found the next step, maybe a problem in the step sequence"
//
//                    httpResponse.Comment = &comment
//                }
//
//            } else {
//
//                // terminate the task
//                statusName = "done"
//
//                // update the done date
//                var timeNow = time.Now()
//                task.DoneDate = &timeNow
//            }
//
//        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//        //  301  = next to '...' step or/and next in '...' interval of seconds  //
//        //______________________________________________________________________//
//        case 301:
//
//            // an interval is setup to the next execution laster
//            if httpResponse.Interval != nil {
//
//                // interval send in seconds
//                var interval = *httpResponse.Interval
//
//                // if defined an interval
//                if interval > 0 {
//
//                    // compute the todo date with the interval
//                    var todoDate = *task.TodoDate
//                    var newTodoDate = todoDate.Add(time.Duration(interval) * time.Second)
//
//                    task.TodoDate = &newTodoDate
//
//                    // logger
//                    log.Println("Change the todoDate to '" + newTodoDate.String() + "'")
//                }
//            }
//
//            // a next step definition
//            if httpResponse.Step != nil {
//
//                // new step name
//                var stepName = *httpResponse.Step
//
//                // flag founded
//                var found = false
//
//                // new step exists in the sequence
//                for _, s := range worker.Definition.Sequence {
//                    // found the asked overwritted step
//                    if s.Name == stepName {
//                        found = true
//                        break
//                    }
//                }
//
//                // step founded in the sequence
//                if found {
//
//                    // overwrite the step
//                    task.Step = stepName
//
//                    // logger
//                    log.Println("Change the next step to '" + stepName + "'")
//
//                } else {
//
//                    // error while the overwriting of the next step
//                    statusName = "error"
//
//                    // forge error comment
//                    var comment = "Impossible to found the next step, maybe the overwrite step name doesn't exists"
//
//                    // logger
//                    log.Println(comment)
//
//                    httpResponse.Comment = &comment
//                }
//            }
//
//        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//        //  420  = cancelled                                                    //
//        //______________________________________________________________________//
//        case 420:
//
//            // Setup the status
//            statusName = "cancelled"
//
//        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//        //  520  = problem                                                      //
//        //______________________________________________________________________//
//        case 520:
//
//            // Setup the status
//            statusName = "problem"
//
//        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//        // other = error ( 5XX : auto-retry )                                   //
//        //______________________________________________________________________//
//        default:
//
//            // Error process
//            statusName = "error"
//
//            // Exception for 5XX status code, auto retry if possible
//            if ((statusCode / 100) == 5) && (task.Retry > 1) {
//                statusName = "todo"
//            }
//
//	}
//
//	// buffer key update
//	task.Buffer = *httpResponse.Buffer
//
//	// status key update
//	task.Status = statusName
//
//	// last update key update
//	timeNow = time.Now()
//	task.LastUpdate = &timeNow
//
//	// retry key update
//	task.Retry = task.Retry - 1
//
//	// comment key update
//	if httpResponse.Comment != nil {
//
//		// retrieve the comment string
//		var comment = *httpResponse.Comment // string
//
//		// change it only if not empty
//		if comment != "" {
//			task.Comment = comment
//		}
//	}
//
//
//
//



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









































///////////////////////////////////////////////////////////////////////////////

//type HttpOut struct {
//	Name      string
//	Arguments JsonB
//	Buffer    JsonB
//}
//
//func (w Worker) CallHttp(task Task, step Step) (httpResponse HttpResponse, statusCode int, err error) {
//
//	// initialize the http client
//	httpclient := http.Client{}
//
//	// http call data
//	var dataOut = HttpOut{
//		Name:      task.Name,
//		Arguments: task.Arguments,
//		Buffer:    task.Buffer}
//
//	// encode the http call data
//	jsonValue, err := json.Marshal(dataOut)
//
//	if err != nil {
//		log.Fatalln("Error while encode the http call data", err)
//
//		return httpResponse, statusCode, err
//	}
//
//	// create the http request
//	req, err := http.NewRequest("POST", step.Url, bytes.NewBuffer(jsonValue))
//
//	if err != nil {
//		log.Fatalln("Error while create the http resquest", err)
//
//		return httpResponse, statusCode, err
//	}
//
//	// set some headers
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("X-Custom-Header", "my-custom-header")
//
//	// do the http call
//	resp, err := httpclient.Do(req)
//
//	if err != nil {
//		log.Fatalln("Error while do the http call", err)
//
//		return httpResponse, statusCode, err
//	}
//
//	defer resp.Body.Close()
//
//	// read the response body
//	body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		log.Fatalln("Error while read the http body data", err)
//
//		return httpResponse, statusCode, err
//	}
//
//	log.Println("===============================")
//	log.Printf("Post data request was '%s'\n", string(jsonValue))
//	log.Println("Response StatusCode:", resp.StatusCode)
//	log.Println("Response Headers:", resp.Header)
//	log.Println("Response Body:", string(body))
//	log.Println("-------------------------------")
//
//	// decoding the returned body data
//	err = json.Unmarshal(body, &httpResponse)
//
//	if err != nil {
//		log.Fatalln("Error while decoding the http response body", err)
//
//		return httpResponse, statusCode, err
//	}
//
//	// retrieve the statusCode data
//	statusCode = resp.StatusCode
//
//	log.Println("Http call work fine")
//
//	return httpResponse, statusCode, nil
//}

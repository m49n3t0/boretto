package bot

import (
	"errors"
	"github.com/m49n3t0/boretto/models"
	"log"
	"strconv"

	"bytes"
	//	"encoding/json"
	_ "github.com/lib/pq"
	//	"gopkg.in/gorp.v2"
	"io/ioutil"
	"net/http"
	//	"strconv"
	//	"time"
)

///////////////////////////////////////////////////////////////////////////////

// the worker executes the task process
type Worker struct {
	workerPool  chan chan int64
	taskChannel chan int64
	dispatcher  *Dispatcher
	quit        chan bool
}

// worker creation handler
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

				log.Printf("Entry to taskChannel with ID : '%d'\n", strconv.Itoa(int(taskId)))

				// we have received a work request
				if err := worker.DoAction(taskId); err != nil {
					log.Printf("Error while working on task: %s", err.Error())
				}

				log.Println(".", "ENDJOB")

			case <-worker.quit:

				// we have received a signal to stop
				// exit this function
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

//
//  type EndpointResponse struct {
//      Action  EndpointResponseAction  `json:"action,notnull"`
//      Buffer  *JsonB                  `json:"buffer"`
//      Data    *EndpointResponseData   `json:"data"`
//  }
//
//  buffer: interface{},
//  action: ENUM("GOTO","NEXT","GOTO_LATER","NEXT_LATER","RETRY","RETRY_NOW","ERROR","PROBLEM","CANCELED"),
//  data: {
//      step: string                                --> optional : only for GOTO/GOTO_LATER action
//      interval: int64 in seconds ( default: 60 )  --> optional : only for GOTO_LATER/NEXT_LATER/RETRY action
//      comment: string                             --> optional : only for ERROR/PROBLEM/CANCELED action
//      detail: map[string]string{}                 --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger
//      no_decrement: bool                          --> optional : only for RETRY action
//  },
//

// do the action for this task with the good action
func (worker *Worker) DoAction(id int64) error {

	// retrieve a task
	task, err := worker.dispatcher.getTask(id)
	if err != nil {
		log.Fatalln("Error while fetching one task", err)
	}

	// logger
	log.Printf("Working on task '%d-%s' on step: '%s'\n", strconv.Itoa(int(task.Id)), task.Function, task.Step)

	// --------------------------------------------------------------------- //

	log.Println("42.0.1----------------------------------------------")

	// get robot for this task version
	definition, ok := worker.dispatcher.definitions[task.Version]
	if !ok {
		log.Println("Robot definition for this version doesn't exists")
		return errors.New("Robot definition for this version doesn't exists")
	}

	log.Println("42.0.2----------------------------------------------")

	log.Println(definition.Sequence)

	log.Println("42.0.3----------------------------------------------")

	// vars on step
	var step *models.Step

	// get the actual step in the sequence
	for _, s := range definition.Sequence {

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

	log.Println(step)

	log.Println("42.1.2----------------------------------------------")

	// check step found
	if step == nil {
		log.Println("Step not found for this version")
		return errors.New("Step not found for this version")
	}

	log.Println("42.2----------------------------------------------")

	// get the associated endpoint
	endpoint, ok := worker.dispatcher.endpoints[step.EndpointID]
	if !ok {
		log.Println("Associated endpoint to this step doesn't exists")
		return errors.New("Associated endpoint to this step doesn't exists")
	}

	log.Println("42.3----------------------------------------------")

	// --------------------------------------------------------------------- //

	// change the task data
	task.Status = "DOING"

	// last update key
	task.LastUpdate = time.Now()

	// fetch the object
	err := dispatcher.db.Update(task)

	if err != nil {
		log.Println("Error while updating the task for status lock : ", err)
		return err
	}

	log.Println("42.4----------------------------------------------")

	// --------------------------------------------------------------------- //

	log.Println("-~-~-~/endpoint\\~-~-~-")

	log.Println(endpoint)

	err = worker.DoHttpAction(endpoint)
	if err != nil {
		log.Println("Error while do the http action : ", err)
		return err
	}

    //
    //  type EndpointResponse struct {
    //      Action  EndpointResponseAction  `json:"action,notnull"`
    //      Buffer  *JsonB                  `json:"buffer"`
    //      Data    *EndpointResponseData   `json:"data"`
    //  }
    //
    //  buffer: interface{},
    //  action: ENUM("GOTO","NEXT","GOTO_LATER","NEXT_LATER","RETRY","RETRY_NOW","ERROR","PROBLEM","CANCELED"),
    //  data: {
    //      step: string                                --> optional : only for GOTO/GOTO_LATER action
    //      interval: int64 in seconds ( default: 60 )  --> optional : only for GOTO_LATER/NEXT_LATER/RETRY action
    //      comment: string                             --> optional : only for ERROR/PROBLEM/CANCELED action
    //      detail: map[string]string{}                 --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger
    //      no_decrement: bool                          --> optional : only for RETRY action
    //  },
    //

	// vars
	var httpResponse HttpResponse
	var statusCode int

	// actual step process
	if foundActualStep {

		// do the http call to retrieve API data/informations
		httpResponse, statusCode, err = worker.CallHttp(task, actualStep)

		if err != nil {

			// error while doing the http call on api's
			var comment = "Error while doing the http call on api's"

			log.Fatalln(comment)

			// forge an error status code
			statusCode = 500

			// forge a fake http response
			httpResponse = HttpResponse{
				Buffer:  &task.Buffer, // repush the same buffer
				Comment: &comment}     // forge a comment

		}

	} else {

		// no step associated for this step name
		var comment = "Error while check the actual step informations"

		log.Fatalln(comment)

		// forge an error status code
		statusCode = 600

		// forge a fake http response
		httpResponse = HttpResponse{
			Buffer:  &task.Buffer, // repush the same buffer
			Comment: &comment}     // forge a comment

	}









    //
    //  type EndpointResponse struct {
    //      Action  EndpointResponseAction  `json:"action,notnull"`
    //      Buffer  *JsonB                  `json:"buffer"`
    //      Data    *EndpointResponseData   `json:"data"`
    //  }
    //
    //  buffer: interface{},
    //  action: ENUM("GOTO","NEXT","GOTO_LATER","NEXT_LATER","RETRY","RETRY_NOW","ERROR","PROBLEM","CANCELED"),
    //  data: {
    //      step: string                                --> optional : only for GOTO/GOTO_LATER action
    //      interval: int64 in seconds ( default: 60 )  --> optional : only for GOTO_LATER/NEXT_LATER/RETRY action
    //      comment: string                             --> optional : only for ERROR/PROBLEM/CANCELED action
    //      detail: map[string]string{}                 --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger
    //      no_decrement: bool                          --> optional : only for RETRY action
    //  },
    //

    if response.Buffer defined ... {

        // update the local buffer

    }

    //        // buffer key update
    //        task.Buffer = *httpResponse.Buffer
    //
    //        // status key update
    //        task.Status = statusName
    //
    //        // last update key update
    //        timeNow = time.Now()
    //        task.LastUpdate = &timeNow
    //
    //        // retry key update
    //        task.Retry = task.Retry - 1
    //
    //        // comment key update
    //        if httpResponse.Comment != nil {
    //
    //            // retrieve the comment string
    //            var comment = *httpResponse.Comment // string
    //
    //            // change it only if not empty
    //            if comment != "" {
    //                task.Comment = comment
    //            }
    //        }
    //
    //        log.Println("Task updation")
    //
    //        log.Println("-~-~-~/endpoint\\~-~-~-")
    //
    //        // --------------------------------------------------------------------- //
    //
    //        log.Println("42.9----------------------------------------------")
    //
    //        log.Println("Task updation")
    //
    //        // change the task data
    //        task.Status = "TODO"
    //
    //        // last update key
    //        task.LastUpdate = time.Now()
    //
    //        // fetch the object
    //        err := dispatcher.db.Update(task)
    //
    //        if err != nil {
    //            log.Println("Error while updating the task result : ", err)
    //            return err
    //        }
    //
    //        log.Println("42.10----------------------------------------------")
    //
    //        log.Println("Task updated")
    //
    //        return nil


    // vars for process function response
    var ( task, mErr ) ( *models.Task, error )

    // do the action correctly
    switch response.Action {

        // GOTO action
        case models.EndpointResponseAction_GOTO:
            task, mErr = worker.processActionGoto( task, definition, response)

        // GOTO_LATER action
        case models.EndpointResponseAction_GOTO_LATER :
            task, mErr = worker.processActionGotoLater( task, definition, response)

        // NEXT action
        case models.EndpointResponseAction_NEXT :
            task, mErr = worker.processActionNext( task, definition, response)

        // NEXT_LATER action
        case models.EndpointResponseAction_NEXT_LATER :
            task, mErr = worker.processActionNextLater( task, definition, response)

        // RETRY_NOW action
        case models.EndpointResponseAction_RETRY_NOW :
            task, mErr = worker.processActionRetryNow( task, definition, response)

        // RETRY action
        case models.EndpointResponseAction_RETRY :
            task, mErr = worker.processActionRetry( task, definition, response)

        // CANCELED action
        case models.EndpointResponseAction_CANCELED :
            task, mErr = worker.processActionCanceled( task, definition, response)

        // PROBLEM action
        case models.EndpointResponseAction_PROBLEM:
            task, mErr = worker.processActionProblem( task, definition, response)

        // DEFAULT / ERROR action
        default :
            task, mErr = worker.processActionError( task, definition, response)
    }

    return worker.updateTask( task, mErr )
}


// Function to process the GOTO_LATER action
func (worker *Worker) processActionGotoLater(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

    // interval settings
    //

    //default
    var interval = 60

    // interval is correctly defined ?
    if response.Data.Interval != nil && *response.Data.Interval > 60 {
        interval = *response.Data.Interval
    }

    // update the task TodoDate key
    task.TodoDate = task.TodoDate.Add( time.Duration( interval ) * time.Second )

    // logger
    log.Println("TodoDate updated to '%s'", task.TodoDate) // task.TodoDate.String()

    // later == retry
    task.Retry = task.Retry - 1

    return worker.processGoto(task, definition, response)
}


// Function to process the GOTO action
func (worker *Worker) processActionGoto(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

    // step settings
    //

    // step is defined ?
    if response.Data.Step == nil || *response.Data.Step == "" {
        return task, errors.New("Missing step parameter from API response for GOTO actions")
    }

    // flag to know if found or not
    var found = false

    // asked step exists in the sequence
    for _, s := range definition.Sequence {
        // found the asked step
        if s.Name == *response.Data.Step {
            found = true
            break
        }
    }

    // not found, error
    if !found {
        return task, errors.New("Impossible to found the asked step from API response")
    }

    // setup the new step
    task.Step = *response.Data.Step

    // logger
    log.Println("Goto step updated to '%s'", task.Step)

    return task, nil
}


// Function to process the NEXT_LATER action
func (worker *Worker) processActionNextLater(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

    // interval settings default
    var interval = 60

    // interval is correctly defined ?
    if response.Data.Interval != nil && *response.Data.Interval > 60 {
        interval = *response.Data.Interval
    }

    // update the task TodoDate key
    task.TodoDate = task.TodoDate.Add( time.Duration( interval ) * time.Second )

    // logger
    log.Println("TodoDate updated to '%s'", task.TodoDate) // task.TodoDate.String()

    // later == retry
    task.Retry = task.Retry - 1

    return worker.processNext(task, definition, response)
}


// Function to process the NEXT action
func (worker *Worker) processActionNext(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

    // next step settings
    //

    // flag to know if actual step found or not
    var found = false

    // retrieve the next step data
    for _, s := range definition.Sequence {
        // actual founded, this one is the classic next step
        if found {
            // setup the new step
            task.Step = *response.Data.Step
            task.Status = models.TaskStatus_TODO

            // logger
            log.Println("Next step updated to '%s'", task.Step)

            return task, nil
        }
        // this actual step was here, founded
        if s.Name == task.Step {
            found = true
        }
    }

    // no next step found, error
    return nil, errors.New("Impossible to found the next step")
}


// Function to process the RETRY_NOW action
func (worker *Worker) processActionRetryNow(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

    // step setup
    task.Status = models.TaskStatus_TODO
    task.Retry = task.Retry - 1

    return task, nil
}


// Function to process the RETRY action
func (worker *Worker) processActionRetry(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

    // interval settings default
    var interval = 60

    // interval is correctly defined ?
    if response.Data.Interval != nil && *response.Data.Interval > 60 {
        interval = *response.Data.Interval
    }

    // update the task TodoDate key
    task.TodoDate = task.TodoDate.Add( time.Duration( interval ) * time.Second )

    // logger
    log.Println("TodoDate updated to '%s'", task.TodoDate)
    //log.Println("TodoDate updated to '%s'", task.TodoDate.String())

    // no_decrement settings
    //

    // not exists/defined no_decrement flag
    if response.Data.NoDecrement != nil && response.Data.NoDecrement == false {
        // later == retry
        task.Retry = task.Retry - 1
    }

    return task, nil
}














// Function to process the RETRY action
func (worker *Worker) processActionRetry(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

        // CANCELED actions
        case models.EndpointResponseAction_CANCELED :

            //      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
            //      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

// Function to process the RETRY action
func (worker *Worker) processActionRetry(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

        // ERROR & PROBLEM actions
        default :

            //      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
            //      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

    }


// Function to process the RETRY action
func (worker *Worker) processActionRetry(task *models.Task, definition *models.Definition, response *models.EndpointResponse) (*models.Task, error) {

        // ERROR & PROBLEM actions
        default :

            //      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
            //      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

    }












func (worker *Worker) updateTask(task *models.Task, mErr error) error {

    // if callback return an error
    if mErr != nil {

        // set correct status/comment
        task.Status = models.TaskStatus_ERROR
        task.Comment = mErr.Error()
        task.Retry = task.Retry - 1

        log.Println(mErr)
    }

    // last update date key
	task.LastUpdate = time.Now()

	// update the database object
    err := dispatcher.db.Update(task)
	if err != nil {
		log.Printf("Error while updating the task result : %s", err)
		return err
	}

	log.Println("42.10----------------------------------------------")

	log.Println("Task updated")

	return nil
}















///////////////////////////////////////////////////////////////////////////////




///////////////////////////////////////////////////////////////////////////////

//type HttpOut struct {
//	Name      string
//	Arguments JsonB
//	Buffer    JsonB
//}

//func (w Worker) CallHttp(task Task, step Step) (httpResponse HttpResponse, statusCode int, err error) {
func (worker *Worker) CallHttp(task *models.Task, step *models.Step) (error, error) {

	// initialize the http client
	httpclient := http.Client{}

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

	body := &bytes.Buffer{}

	//	if e.Body != "" {
	//		body = bytes.NewBuffer([]byte(e.Body))
	//	}
	//	if len(e.BasicAuthUser) > 0 || len(e.BasicAuthPassword) > 0 {
	//	    req.SetBasicAuth(e.BasicAuthUser, e.BasicAuthPassword)
	//	}
	//	return req, err
	//}

	var endpoint *models.HttpEndpoint

	// create the http request
	req, err := http.NewRequest(endpoint.Method, endpoint.Url, body) //  bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Fatalln("Error while create the http resquest", err)
		//	return httpResponse, statusCode, err
		return nil, nil
	}

	//	// set some headers
	//	req.Header.Set("Content-Type", "application/json")
	//	req.Header.Set("X-Custom-Header", "my-custom-header")

	//
	//	for k, v := range t.Headers {
	//		req.Header.Set(k, v)
	//	}
	//
	//    tr := &http.Transport{
	//        TLSClientConfig: &tls.Config{InsecureSkipVerify: t.IgnoreVerifySSL},
	//    }
	//    client := &http.Client{Transport: tr}
	//
	//	start := time.Now()
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		return nil, err
	//	}

	// do the http call
	resp, err := httpclient.Do(req)

	if err != nil {
		log.Println("Error while do the http call", err)
		//return httpResponse, statusCode, err
		return nil, nil
	}

	//	var bb []byte
	//	if resp.Body != nil {
	//		defer resp.Body.Close()
	//		var errr error
	//		bb, errr = ioutil.ReadAll(resp.Body)
	//		if errr != nil {
	//			return nil, errr
	//		}
	//		r.Body = string(bb)
	//
	//		bodyJSONArray := []interface{}{}
	//		if err := json.Unmarshal(bb, &bodyJSONArray); err != nil {
	//			bodyJSONMap := map[string]interface{}{}
	//			if err2 := json.Unmarshal(bb, &bodyJSONMap); err2 == nil {
	//				r.BodyJSON = bodyJSONMap
	//			}
	//		} else {
	//			r.BodyJSON = bodyJSONArray
	//		}
	//	}
	//
	//	r.StatusCode = resp.StatusCode

	defer resp.Body.Close()

	// read the response body
	iobody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error while read the http body data", err)
		//return httpResponse, statusCode, err
		return nil, nil
	}

	log.Println("===============================")
	//	log.Printf("Post data request was '%s'\n", string(jsonValue))
	log.Println("Response StatusCode:", resp.StatusCode)
	log.Println("Response Headers:", resp.Header)
	log.Println("Response Body:", string(iobody))
	log.Println("-------------------------------")

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

	return nil, nil
}

































































































































































































































































































































































































//// do the action for this task with the good action
//func (worker *Worker) DoAction(id int64) error {
//
//	// retrieve a task
//	task, err := worker.dispatcher.getTask(id)
//	if err != nil {
//		log.Fatalln("Error while fetching one task", err)
//	}
//
//	// logger
//    log.Printf("Working on task '%d-%s' on step: '%s'\n", strconv.Itoa(int(task.Id)), task.Function, task.Step)
//
//	// --------------------------------------------------------------------- //
//
//	log.Println("42.0.1----------------------------------------------")
//
//	// get robot for this task version
//	definition, ok := worker.dispatcher.definitions[task.Version]
//	if !ok {
//		log.Println("Robot definition for this version doesn't exists")
//		return errors.New("Robot definition for this version doesn't exists")
//	}
//
//	log.Println("42.0.2----------------------------------------------")
//
//	log.Println(definition.Sequence)
//
//	log.Println("42.0.3----------------------------------------------")
//
//	// vars on step
//	var step *models.Step
//
//	// get the actual step in the sequence
//	for _, s := range definition.Sequence {
//
//		log.Println("42.0.4----------------------------------------------")
//		log.Printf("task step : %s / step name : %s , %s , %s", task.Step, s.Name, s.EndpointType, s.EndpointID)
//
//		// check with the local task
//		if task.Step == s.Name {
//			log.Println("42.0.5----FOUND")
//			step = &s
//			break
//		}
//	}
//
//	log.Println("42.1.1----------------------------------------------")
//
//	log.Println(step)
//
//	log.Println("42.1.2----------------------------------------------")
//
//	// check step found
//	if step == nil {
//		log.Println("Step not found for this version")
//		return errors.New("Step not found for this version")
//	}
//
//	log.Println("42.2----------------------------------------------")
//
//	// get the associated endpoint
//	endpoint, ok := worker.dispatcher.endpoints[step.EndpointID]
//	if !ok {
//		log.Println("Associated endpoint to this step doesn't exists")
//		return errors.New("Associated endpoint to this step doesn't exists")
//	}
//
//	log.Println("42.3----------------------------------------------")
//
//	// --------------------------------------------------------------------- //
//
//	// change the task data
//	task.Status = "DOING"
//
//	// update in database the task
//	err = worker.dispatcher.updateTask(task)
//
//	if err != nil {
//		log.Println("Error while updating the task for status lock : ", err)
//		return err
//	}
//
//	log.Println("42.4----------------------------------------------")
//
//
//
//
//	// --------------------------------------------------------------------- //
//
//
//
//
//
//
//
//
//
//
//
//
//
//	log.Println("-~-~-~/endpoint\\~-~-~-")
//    log.Println( endpoint )
//    err = worker.DoHttpAction((endpoint).(*models.HttpEndpoint))
//	log.Println("-~-~-~/endpoint\\~-~-~-")
//
//
//
//
//
//
//
//
//
//
//
//
//	// --------------------------------------------------------------------- //
//
//
//
//	log.Println("42.9----------------------------------------------")
//
//	log.Println("Task updation")
//
//	// update in database the task
//	err = worker.dispatcher.updateTask(task)
//
//	if err != nil {
//		log.Println("Error while updating the task result : ", err)
//		return err
//	}
//
//	log.Println("42.10----------------------------------------------")
//
//	log.Println("Task updated")
//
//	return nil
//}
//
//
//
//
//
//
//
////// Response http from each API response
////type HttpResponse struct {
////	Buffer   *JsonB
////	Interval *int
////	Step     *string
////	Comment  *string
////}
//
//
//// do the call on http endpoints
//func (worker *Worker) DoHttpAction(endpoint *models.HttpEndpoint) error {
//
//
//
//
//
//
//
////	// vars
////	var httpResponse HttpResponse
////	var statusCode int
//
////	// actual step process
////	if foundActualStep {
////
////		// do the http call to retrieve API data/informations
////		httpResponse, statusCode, err = worker.CallHttp(task, actualStep)
////
////		if err != nil {
////
////			// error while doing the http call on api's
////			var comment = "Error while doing the http call on api's"
////
////			log.Fatalln(comment)
////
////			// forge an error status code
////			statusCode = 500
////
////			// forge a fake http response
////			httpResponse = HttpResponse{
////				Buffer:  &task.Buffer, // repush the same buffer
////				Comment: &comment}     // forge a comment
////
////		}
////
////    }
//    //	} else {
//    //
//    //		// no step associated for this step name
//    //		var comment = "Error while check the actual step informations"
//    //
//    //		log.Fatalln(comment)
//    //
//    //		// forge an error status code
//    //		statusCode = 600
//    //
//    //		// forge a fake http response
//    //		httpResponse = HttpResponse{
//    //			Buffer:  &task.Buffer, // repush the same buffer
//    //			Comment: &comment}     // forge a comment
//    //
//    //	}
//
////	// default value
////	var statusName = "todo"
////
////	// switch on each status code
////	switch statusCode {
////
////        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
////        //  200  = next                                                         //
////        //______________________________________________________________________//
////        case 200:
////
////            // not an ending step
////            if !actualStep.EndStep {
////
////                var found = false
////
////                // retrieve the next step data
////                for _, s := range worker.Definition.Sequence {
////                    if found {
////                        // set the next step
////                        task.Step = s.Name
////                        break
////                    }
////                    if s.Name == task.Step {
////                        found = true
////                    }
////                }
////
////                // no next step found, error
////                if !found {
////
////                    // error status due to no next step founded
////                    statusName = "error"
////
////                    // forge error comment
////                    var comment = "Impossible to found the next step, maybe a problem in the step sequence"
////
////                    httpResponse.Comment = &comment
////                }
////
////            } else {
////
////                // terminate the task
////                statusName = "done"
////
////                // update the done date
////                var timeNow = time.Now()
////                task.DoneDate = &timeNow
////            }
////
////        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
////        //  301  = next to '...' step or/and next in '...' interval of seconds  //
////        //______________________________________________________________________//
////        case 301:
////
////            // an interval is setup to the next execution laster
////            if httpResponse.Interval != nil {
////
////                // interval send in seconds
////                var interval = *httpResponse.Interval
////
////                // if defined an interval
////                if interval > 0 {
////
////                    // compute the todo date with the interval
////                    var todoDate = *task.TodoDate
////                    var newTodoDate = todoDate.Add(time.Duration(interval) * time.Second)
////
////                    task.TodoDate = &newTodoDate
////
////                    // logger
////                    log.Println("Change the todoDate to '" + newTodoDate.String() + "'")
////                }
////            }
////
////            // a next step definition
////            if httpResponse.Step != nil {
////
////                // new step name
////                var stepName = *httpResponse.Step
////
////                // flag founded
////                var found = false
////
////                // new step exists in the sequence
////                for _, s := range worker.Definition.Sequence {
////                    // found the asked overwritted step
////                    if s.Name == stepName {
////                        found = true
////                        break
////                    }
////                }
////
////                // step founded in the sequence
////                if found {
////
////                    // overwrite the step
////                    task.Step = stepName
////
////                    // logger
////                    log.Println("Change the next step to '" + stepName + "'")
////
////                } else {
////
////                    // error while the overwriting of the next step
////                    statusName = "error"
////
////                    // forge error comment
////                    var comment = "Impossible to found the next step, maybe the overwrite step name doesn't exists"
////
////                    // logger
////                    log.Println(comment)
////
////                    httpResponse.Comment = &comment
////                }
////            }
////
////        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
////        //  420  = cancelled                                                    //
////        //______________________________________________________________________//
////        case 420:
////
////            // Setup the status
////            statusName = "cancelled"
////
////        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
////        //  520  = problem                                                      //
////        //______________________________________________________________________//
////        case 520:
////
////            // Setup the status
////            statusName = "problem"
////
////        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
////        // other = error ( 5XX : auto-retry )                                   //
////        //______________________________________________________________________//
////        default:
////
////            // Error process
////            statusName = "error"
////
////            // Exception for 5XX status code, auto retry if possible
////            if ((statusCode / 100) == 5) && (task.Retry > 1) {
////                statusName = "todo"
////            }
////
////	}
////
////	// buffer key update
////	task.Buffer = *httpResponse.Buffer
////
////	// status key update
////	task.Status = statusName
////
////	// last update key update
////	timeNow = time.Now()
////	task.LastUpdate = &timeNow
////
////	// retry key update
////	task.Retry = task.Retry - 1
////
////	// comment key update
////	if httpResponse.Comment != nil {
////
////		// retrieve the comment string
////		var comment = *httpResponse.Comment // string
////
////		// change it only if not empty
////		if comment != "" {
////			task.Comment = comment
////		}
////	}
////
////	log.Println("Task updation")
//
//	return nil
//}
//
//
/////////////////////////////////////////////////////////////////////////////////
//
////type HttpOut struct {
////	Name      string
////	Arguments JsonB
////	Buffer    JsonB
////}
//
////func (w Worker) CallHttp(task Task, step Step) (httpResponse HttpResponse, statusCode int, err error) {
//func (worker *Worker) CallHttp(task *models.Task, step *models.Step) (error, error) {
//
//	// initialize the http client
//	httpclient := http.Client{}
//
////	// http call data
////	var dataOut = HttpOut{
////		Name:      task.Name,
////		Arguments: task.Arguments,
////		Buffer:    task.Buffer}
////
////	// encode the http call data
////	jsonValue, err := json.Marshal(dataOut)
////
////	if err != nil {
////		log.Fatalln("Error while encode the http call data", err)
////
////		return httpResponse, statusCode, err
////	}
//
//	body := &bytes.Buffer{}
//
//	//	if e.Body != "" {
//	//		body = bytes.NewBuffer([]byte(e.Body))
//	//	}
//	//	if len(e.BasicAuthUser) > 0 || len(e.BasicAuthPassword) > 0 {
//	//	    req.SetBasicAuth(e.BasicAuthUser, e.BasicAuthPassword)
//	//	}
//	//	return req, err
//	//}
//
//
//    var endpoint *models.HttpEndpoint
//
//
//	// create the http request
//	req, err := http.NewRequest(endpoint.Method, endpoint.Url, body)  //  bytes.NewBuffer(jsonValue))
//
//	if err != nil {
//		log.Fatalln("Error while create the http resquest", err)
//	//	return httpResponse, statusCode, err
//        return nil, nil
//	}
//
////	// set some headers
////	req.Header.Set("Content-Type", "application/json")
////	req.Header.Set("X-Custom-Header", "my-custom-header")
//
//	//
//	//	for k, v := range t.Headers {
//	//		req.Header.Set(k, v)
//	//	}
//	//
//	//    tr := &http.Transport{
//	//        TLSClientConfig: &tls.Config{InsecureSkipVerify: t.IgnoreVerifySSL},
//	//    }
//	//    client := &http.Client{Transport: tr}
//	//
//	//	start := time.Now()
//	//	resp, err := client.Do(req)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//
//	// do the http call
//	resp, err := httpclient.Do(req)
//
//	if err != nil {
//		log.Println("Error while do the http call", err)
//		//return httpResponse, statusCode, err
//        return nil, nil
//	}
//
//	//	var bb []byte
//	//	if resp.Body != nil {
//	//		defer resp.Body.Close()
//	//		var errr error
//	//		bb, errr = ioutil.ReadAll(resp.Body)
//	//		if errr != nil {
//	//			return nil, errr
//	//		}
//	//		r.Body = string(bb)
//	//
//	//		bodyJSONArray := []interface{}{}
//	//		if err := json.Unmarshal(bb, &bodyJSONArray); err != nil {
//	//			bodyJSONMap := map[string]interface{}{}
//	//			if err2 := json.Unmarshal(bb, &bodyJSONMap); err2 == nil {
//	//				r.BodyJSON = bodyJSONMap
//	//			}
//	//		} else {
//	//			r.BodyJSON = bodyJSONArray
//	//		}
//	//	}
//	//
//	//	r.StatusCode = resp.StatusCode
//
//	defer resp.Body.Close()
//
//	// read the response body
//    iobody, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		log.Println("Error while read the http body data", err)
//		//return httpResponse, statusCode, err
//        return nil, nil
//	}
//
//	log.Println("===============================")
////	log.Printf("Post data request was '%s'\n", string(jsonValue))
//	log.Println("Response StatusCode:", resp.StatusCode)
//	log.Println("Response Headers:", resp.Header)
//	log.Println("Response Body:", string(iobody))
//	log.Println("-------------------------------")
//
////	// decoding the returned body data
////	err = json.Unmarshal(body, &httpResponse)
////
////	if err != nil {
////		log.Fatalln("Error while decoding the http response body", err)
////
////		return httpResponse, statusCode, err
////	}
////
////	// retrieve the statusCode data
////	statusCode = resp.StatusCode
////
////	log.Println("Http call work fine")
////
////	return httpResponse, statusCode, nil
//
//    return nil, nil
//}
//
//
//
//
//
//
//
//
//
//
//
//

package bot

import (
	"errors"
	"github.com/m49n3t0/boretto/machine/models"
	"strconv"
	"bytes"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	//	"encoding/json"
	//	"gopkg.in/gorp.v2"
	//	"strconv"
	//	"time"

    log "github.com/sirupsen/logrus"
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

	// get robot for this task version
	definition, ok := worker.dispatcher.definitions[task.Version]
	if !ok {
		return errors.New("Robot definition for this version doesn't exists")
	}

	// vars on step
	var step *models.Step

	// get the actual step in the sequence
	for _, s := range definition.Sequence {

		// check with the local task
		if task.Step == s.Name {
			step = &s
			break
		}
	}

	// check step found
	if step == nil {
		return errors.New("Step not found for this version")
	}

	// get the associated endpoint
	endpoint, ok := worker.dispatcher.endpoints[step.EndpointID]
	if !ok {
		return errors.New("Associated endpoint to this step doesn't exists")
	}

	// --------------------------------------------------------------------- //

	// change the task data
	task.Status = models.TaskStatus_DOING

    // update the task informations
    err := worker.updateTask(task, nil)
    if err != nil {
        return err
    }

    log.Println("Task in 'DOING' status")

	// --------------------------------------------------------------------- //

    // do the http calls
	err = worker.HttpCall( task, endpoint )
	if err != nil {
        return worker.updateTask( task, err )
	}

	// --------------------------------------------------------------------- //

    // vars for process function response
    var ( task, mErr ) ( *models.Task, error )
    var ( task, mErr ) ( *models.Task, error )

    // update the local buffer from the API return
    if response.Buffer != nil {
        task.Buffer = *response.Buffer
    }

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

        // ERROR action
        case models.EndpointResponseAction_ERROR:
            task, mErr = worker.processActionError( task, definition, response)

        // action not matched
        default :
            mErr = errors.New("Action '%s' isn't matched by executor process", response.Action)
    }

    return worker.updateTask( task, mErr )
}

func (worker *Worker) updateTask(task *models.Task, mErr error) error {

    // if callback is an error, mistaken the task
    if mErr != nil {

        // status MISTAKE
        task.Status = models.TaskStatus_MISTAKE

        // get the comment from error message
        task.Comment = mErr.Error()

        // always lss a retry when we do an error inside
        task.Retry = task.Retry - 1

        log.Printf("Mistake appear : '%s'", mErr)
    }

    // always update the last update date key
	task.LastUpdate = time.Now()

	// update the database object
    err := dispatcher.db.Update(task)
	if err != nil {
		log.Printf("Error while updating the task result : %s", err)
		return err
	}

	log.Println("Task updated")

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func (worker *Worker) CallHttp(task *models.Task, endpoint *models.Endpoint) (error, error) {

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

func (w Worker) CallHttp(task Task, step Step) (httpResponse HttpResponse, statusCode int, err error) {

	// initialize the http client
	httpclient := http.Client{}

	// http call data
	var dataOut = HttpOut{
		Name:      task.Name,
		Arguments: task.Arguments,
		Buffer:    task.Buffer}

	// encode the http call data
	jsonValue, err := json.Marshal(dataOut)

	if err != nil {
		log.Fatalln("Error while encode the http call data", err)

		return httpResponse, statusCode, err
	}

	// create the http request
	req, err := http.NewRequest("POST", step.Url, bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Fatalln("Error while create the http resquest", err)

		return httpResponse, statusCode, err
	}


	path := fmt.Sprintf("%s%s", e.URL, e.Path)
	method := e.Method
	body := &bytes.Buffer{}

	if e.Body != "" {
		body = bytes.NewBuffer([]byte(e.Body))
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}





















	// set some headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-Header", "my-custom-header")


	if len(e.BasicAuthUser) > 0 || len(e.BasicAuthPassword) > 0 {
	    req.SetBasicAuth(e.BasicAuthUser, e.BasicAuthPassword)
	}

	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}




	// do the http call
	resp, err := httpclient.Do(req)

	if err != nil {
		log.Fatalln("Error while do the http call", err)

		return httpResponse, statusCode, err
	}

	defer resp.Body.Close()


    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: t.IgnoreVerifySSL},
    }
    client := &http.Client{Transport: tr}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	r.TimeSeconds = elapsed.Seconds()
	r.TimeHuman = fmt.Sprintf("%s", elapsed)










	// read the response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("Error while read the http body data", err)

		return httpResponse, statusCode, err
	}

	log.Println("===============================")
	log.Printf("Post data request was '%s'\n", string(jsonValue))
	log.Println("Response StatusCode:", resp.StatusCode)
	log.Println("Response Headers:", resp.Header)
	log.Println("Response Body:", string(body))
	log.Println("-------------------------------")

	// decoding the returned body data
	err = json.Unmarshal(body, &httpResponse)

	if err != nil {
		log.Fatalln("Error while decoding the http response body", err)

		return httpResponse, statusCode, err
	}

	// retrieve the statusCode data
	statusCode = resp.StatusCode




	var bb []byte
	if resp.Body != nil {
		defer resp.Body.Close()
		var errr error
		bb, errr = ioutil.ReadAll(resp.Body)
		if errr != nil {
			return nil, errr
		}
		r.Body = string(bb)

		bodyJSONArray := []interface{}{}
		if err := json.Unmarshal(bb, &bodyJSONArray); err != nil {
			bodyJSONMap := map[string]interface{}{}
			if err2 := json.Unmarshal(bb, &bodyJSONMap); err2 == nil {
				r.BodyJSON = bodyJSONMap
			}
		} else {
			r.BodyJSON = bodyJSONArray
		}
	}

	r.StatusCode = resp.StatusCode























	log.Println("Http call work fine")

	return httpResponse, statusCode, nil
}






































































































































































type WorkerHttpOutput struct {
    ID          int64               `json:"id"`
    Context     string              `json:"context"`
    Arguments   models.JsonB        `json:"arguments"`
    Buffer      models.JsonB        `json:"buffer"`
}

type WorkerHttpInput struct {
	Action models.Action            `json:"action,notnull"`
	Buffer *models.JsonB            `json:"buffer"`
	Data   *WorkerHttpInputData     `json:"data"`
}

type WorkerHttpInputData struct {
	Step        *string             `json:"step"`
	Interval    *int64              `json:"interval"`
	Comment     *string             `json:"comment"`
	NoDecrement *bool               `json:"no_decrement"`
}

func (w Worker) CallHttp(task Task, step Step) (httpResponse HttpResponse, statusCode int, err error) {

	// initialize the http client
	httpclient := http.Client{}

    // body informations
	body := &bytes.Buffer{}

    //

	// http call data
	var dataOut = HttpOut{
		Name:      task.Name,
		Arguments: task.Arguments,
		Buffer:    task.Buffer}


    // output data send to the API
    var output = map[string]interface{}{
        "id": task.ID,
        "context": task.Context,
        "arguments": task.Arguments,
        "buffer": task.Buffer,
    }



	// encode the body data for the call
	jsonBody, err := json.Marshal(output)
	if err != nil {
        return nil, err
	}




	// create the http request
	req, err := http.NewRequest("POST", step.Url, bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Fatalln("Error while create the http resquest", err)

		return httpResponse, statusCode, err
	}


	body := &bytes.Buffer{}

	if e.Body != "" {
		body = bytes.NewBuffer([]byte(e.Body))
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}





















	// set some headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-Header", "my-custom-header")


	if len(e.BasicAuthUser) > 0 || len(e.BasicAuthPassword) > 0 {
	    req.SetBasicAuth(e.BasicAuthUser, e.BasicAuthPassword)
	}

	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}




	// do the http call
	resp, err := httpclient.Do(req)

	if err != nil {
		log.Fatalln("Error while do the http call", err)

		return httpResponse, statusCode, err
	}

	defer resp.Body.Close()


    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: t.IgnoreVerifySSL},
    }
    client := &http.Client{Transport: tr}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	r.TimeSeconds = elapsed.Seconds()
	r.TimeHuman = fmt.Sprintf("%s", elapsed)










	// read the response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("Error while read the http body data", err)

		return httpResponse, statusCode, err
	}

	log.Println("===============================")
	log.Printf("Post data request was '%s'\n", string(jsonValue))
	log.Println("Response StatusCode:", resp.StatusCode)
	log.Println("Response Headers:", resp.Header)
	log.Println("Response Body:", string(body))
	log.Println("-------------------------------")

	// decoding the returned body data
	err = json.Unmarshal(body, &httpResponse)

	if err != nil {
		log.Fatalln("Error while decoding the http response body", err)

		return httpResponse, statusCode, err
	}

	// retrieve the statusCode data
	statusCode = resp.StatusCode




	var bb []byte
	if resp.Body != nil {
		defer resp.Body.Close()
		var errr error
		bb, errr = ioutil.ReadAll(resp.Body)
		if errr != nil {
			return nil, errr
		}
		r.Body = string(bb)

		bodyJSONArray := []interface{}{}
		if err := json.Unmarshal(bb, &bodyJSONArray); err != nil {
			bodyJSONMap := map[string]interface{}{}
			if err2 := json.Unmarshal(bb, &bodyJSONMap); err2 == nil {
				r.BodyJSON = bodyJSONMap
			}
		} else {
			r.BodyJSON = bodyJSONArray
		}
	}

	r.StatusCode = resp.StatusCode























	log.Println("Http call work fine")

	return httpResponse, statusCode, nil
}

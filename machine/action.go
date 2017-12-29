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












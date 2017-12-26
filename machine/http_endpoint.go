package bot

//import (
//	"bytes"
////	"encoding/json"
//	_ "github.com/lib/pq"
//	"github.com/m49n3t0/boretto/models"
////	"gopkg.in/gorp.v2"
//	"io/ioutil"
//	"log"
//	"net/http"
////	"strconv"
////	"time"
//)
//
/////////////////////////////////////////////////////////////////////////////////
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
//

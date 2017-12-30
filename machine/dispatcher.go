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
	"github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////

// to define the dispatcher parameters
type DispatcherParams struct {
	// robot configuration
	Function string
	// dispatcher configuration
	MaxWorker int64
	MaxQueue  int64
	// machine interface object
	Machine *Machine
}

// dispatcher object
type Dispatcher struct {
	// function on what the robot work
	Function string
	// store datas from database
	Definitions map[int64]*models.Definition
	Endpoints   map[string]*models.Endpoint
	// manage the distribution of workflow
	WorkerPool chan chan string
	Queue      chan string
	// manage the quit process
	Signal chan os.Signal
	Quit   chan bool
	// database handler
	DB *pg.DB
	// logger handler
	Logger *logrus.Entry
}

// dispatcher creation handler
func NewDispatcher(params *DispatcherParams) (*Dispatcher, error) {

	// machine interface object
	machine := params.Machine

	// get the logger from interface
	logger, err := machine.GetLogger()
	if err != nil {
		return nil, err
	}
	if logger == nil {
		return nil, errors.New("the logger initialization return an empty object")
	}

	// get the database handler from interface
	database, err := machine.GetDatabase(logger)
	if err != nil {
		return nil, err
	}

	// create the object
	dispatcher := &Dispatcher{
		Function:    params.Function,
		Definitions: make(map[int64]*models.Definition),
		Endpoints:   make(map[string]*models.Endpoint),
		WorkerPool:  make(chan chan string, params.MaxWorker),
		Queue:       make(chan string, params.MaxQueue),
		Signal:      make(chan os.Signal, 2),
		Quit:        make(chan bool),
		DB:          database,
		Logger:      logger.WithField("function", params.Function),
	}

	return dispatcher, nil
}

// stop signals programmatically
func (dispatcher *Dispatcher) Stop() {
	go func() {
		dispatcher.Quit <- true
	}()
}

// stop signals from system
func (dispatcher *Dispatcher) Signal() {
	go func() {

		// link system signal to the dispatcher signal
		signal.Notify(dispatcher.Signal, os.Interrupt, syscall.SIGTERM)

		// when receive syscall signal
		<-dispatcher.Signal

		// do a stopper
		dispatcher.Stop()
	}()
}

///////////////////////////////////////////////////////////////////////////////

// retrieve robot configuration for this function from database
func (dispatcher *Dispatcher) GetRobotConfiguration() error {

	// function logger
	log := dispatcher.Logger

	log.Info("Get the robot configuration")

	// object to fetch
	var robots []*models.Robot

	// get the robot data
	err := dispatcher.DB.
		Model(&robots).
		Where(models.TblRobot_Function+" = ?", dispatcher.Function).
		Where(models.TblRobot_Status+" = ?", "ACTIVE").
		Select()

	if err != nil {
		log.WithError(err).Warn("Select robot configuration error")
		return err
	}

	// no elements, error
	if len(robots) == 0 {
		return errors.New("No robots definition active foudn for this function")
	}

	// store the endpoint to fetch informations
	var stepIDs []string

	// store into dispatcher definitions data
	for _, robot := range robots {

		// remap by version the definitions
		dispatcher.Definitions[robot.Version] = &robot.Definition

		// save the step IDs to fetch after
		for _, step := range robot.Definition.Sequence {
			stepIDs = append(stepIDs, step.EndpointID)
		}
	}

	// object to fetch
	var endpoints []*models.Endpoint

	// get the endpoint data
	err := dispatcher.DB.
		Model(&endpoints).
		Where(models.TblEndpoint_ID+" IN ( ? )", pg.In(stepIDs)).
		Select()

	if err != nil {
		log.WithError(err).Warn("Select robot endpoint steps error")
		return err
	}

	// store locally each endpoints
	for _, endpoint := range endpoints {
		dispatcher.Endpoints[endpoint.ID] = &endpoint
	}

	log.Info("Robot configuration loaded")

	return nil
}

// retrieve the available task ID list
func (dispatcher *Dispatcher) GetTasks() error {

	// function logger
	log := dispatcher.Logger

	log.Info("Reading task ID")

	// where store the ID list
	var IDs []string

	// working on this model
	var task models.Task

	// fetch the available ID list
	err := dispatcher.DB.
		Model(&task).
		Column(models.ColTask_ID).
		Where(models.TblTask_Status+" = ?", models.TaskStatus_TODO).
		Where(models.TblTask_Function+" = ?", dispatcher.Function).
		Where(models.TblTask_Retry+" > ?", 0).
		Where(models.TblTask_TodoDate + " <= NOW()").
		OrderExpr(models.TblTask_TodoDate + " ASC").
		Select(&IDs)

	if err != nil {
		log.WithError(err).Warn("Select task IDs error")
		return err
	}

	// push in the queue the ID informations
	for _, ID := range IDs {
		dispatcher.Queue <- ID
	}

	return nil
}

//// listen the database event channel to do some actions
//func (dispatcher *Dispatcher) Listen() error {
//
//    // function logger
//    log := dispatcher.log
//
//	// build the event chan to listen
//	var eventName = "event_task"
//    //var eventName = fmt.Sprintf("event_task_%s", dispatcher.Function)
//
//	// get the database listener for this robot
//	listener := dispatcher.db.Listen(eventName)
//
//	//    defer ln.Close()
//
//	// get the channel
//	channel := listener.Channel()
//
//	for {
//		select {
//            case event := <-channel:
//
//                // model of the event data
//                var notification models.Notification
//
//                // decode event json data
//                err := json.Unmarshal([]byte(event.Payload), &notification)
//
//                if err != nil {
//                    log.Println("Error while decode the database notification", err)
//                    continue
//                }
//
//                // restart event process
//                if notification.Action == "RESTART" {
//
//                    log.Println("Received event for restart the robot")
//
//                    // XXX: need to develop this part of restart data
//
//                    continue
//                }
//
//                // task event process
//                if notification.Action == "TASK" && notification.Data.TaskID != 0 {
//
//                    log.Printf("Received event for task ID : %d", notification.Data.TaskID)
//
//                    // push task ID event to the channel
//                    dispatcher.queue <- notification.Data.TaskID
//
//                    continue
//                }
//
//            // recursive auto-pull
//            //		case <-time.After(90 * time.Second):
//            case <-time.After(10 * time.Second):
//
//                // logger
//                log.Println("Auto-pull")
//
//                // task pull
//                go dispatcher.getTaskIDs()
//		}
//	}
//}

//// retrieve a task by ID
//func (dispatcher *Dispatcher) getTask(id int64) (*models.Task, error) {
//
//	log.Println("Read task")
//
//	// model to fetch
//	var task models.Task
//
//	// fetch the object
//	err := dispatcher.db.
//		Model(&task).
//		Where(models.TblTask_Id+" = ?", id).
//		Where(models.TblTask_Status+" = ?", "TODO").
//		Where(models.TblTask_Function+" = ?", dispatcher.function).
//		Where(models.TblTask_Retry+" > ?", 0).
//		Where(models.TblTask_TodoDate + " <= NOW()").
//		First()
//
//	if err != nil {
//		if err == pg.ErrNoRows {
//			log.Println("Task not found")
//			return nil, errors.New("Task not found")
//		}
//
//		log.Println("Error while retrieve a task")
//		return nil, err
//	}
//
//	return &task, nil
//}

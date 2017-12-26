package bot

import (
	"encoding/json"
	"errors"
	"github.com/go-pg/pg"
	"github.com/m49n3t0/boretto/models"
	"log"
	"time"
)

///////////////////////////////////////////////////////////////////////////////

// retrieve a database connection
func (dispatcher *Dispatcher) dbConnect() error {

	// build database host address
	var addr = DB_HOST

	if DB_PORT != "" {
		addr += ":" + DB_PORT
	}

	// pg database connector
	db := pg.Connect(&pg.Options{
		Addr:       addr,
		User:       DB_USER,
		Password:   DB_PASSWORD,
		Database:   DB_DATABASE,
		MaxRetries: 2,
	})

	// check connection
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		log.Println("Problem while check database connection")
		return err
	}

	// build the query logger
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		// XXX : maybe only use UnformattedQuery option ( a debug flag ? )
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})

	// reference it into dispatcher
	dispatcher.db = db

	return nil
}

// close the database connection
func (dispatcher *Dispatcher) dbDisconnect() error {

	// disconnect to the database
	err := dispatcher.db.Close()
	if err != nil {
		log.Println("Problem while database disconnect")
		return err
	}

	return nil
}

// retrieve robot configuration for this function from database
func (dispatcher *Dispatcher) getConfiguration() error {

	log.Println("Get the robot configuration")

	// object to fetch
	var robots []*models.Robot

	// get the robot data
	err := dispatcher.db.
		Model(&robots).
		Where(models.TblRobot_Function+" = ?", dispatcher.function).
		Where(models.TblRobot_Status+" = ?", "ACTIVE").
		Select()

	log.Println("1================")

	if err != nil {
		log.Println("Error while select robots", err)
		return err
	}

	log.Println("2================")

	// no elements, error
	if len(robots) == 0 {
		log.Println("Error no robots definition found")
		return errors.New("No robots definition found")
	}

	log.Println("3================")

	// store the endpoint to fetch informations
	var stepIDs []int64

	log.Println("4.0================")

	// store into dispatcher definitions data
	for _, robot := range robots {

		log.Println("4.1.............")
		log.Printf("%+v", robot)

		// remap by version the definitions
		dispatcher.definitions[robot.Version] = &robot.Definition

		log.Println("4.2.............")

		// save the step IDs to fetch after
		for _, step := range robot.Definition.Sequence {
			stepIDs = append(stepIDs, step.EndpointID)
		}
	}

	log.Println("5================")

	// object to fetch
	var endpoints []*models.Endpoint

	// get the endpoint data
	err := dispatcher.db.
		Model(&endpoints).
		Where(models.TblEndpoint_Id+" IN ( ? )", pg.In(stepIDs)).
		Select()

	if err != nil {
		log.Printf("Error while retrieve endpoints", err)
		return err
	}

	// store locally each endpoints
	for _, endpoint := range endpoints {
		dispatcher.endpoints[endpoint.Id] = &endpoint
	}

	log.Println("6================")

	return nil
}

//// retrieve the available task ID list
//func (dispatcher *Dispatcher) getTaskIDs() error {
//
//	log.Println("Read all task IDs")
//
//	// where store the ID list
//	var IDs []int64
//
//	log.Println("9.0=========================")
//
//	// working on this model
//	var task models.Task
//
//	// fetch the available ID list
//	err := dispatcher.db.
//		Model(&task).
//		Column("id").
//		OrderExpr(models.TblTask_Id+" ASC").
//		Where(models.TblTask_Status+" = ?", "TODO").
//		Where(models.TblTask_Function+" = ?", dispatcher.function).
//		Where(models.TblTask_Retry+" > ?", 0).
//		Where(models.TblTask_TodoDate + " <= NOW()").
//		Select(&IDs)
//
//	log.Println("9.1=========================")
//
//	if err != nil {
//		log.Println("Error while fetching the task ids")
//		log.Println(err)
//		return err
//	}
//
//	log.Println("9.2=========================")
//
//	log.Println("All rows:")
//
//	for x, id := range IDs {
//
//		dispatcher.queue <- id
//
//		log.Printf("9.3=====>  %d : %+v\n", x, id)
//	}
//
//	return nil
//}
//
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
//
//// update a task data
//func (dispatcher *Dispatcher) updateTask(task *models.Task) error {
//
//	// last update key
//	task.LastUpdate = time.Now()
//
//	// fetch the object
//	err := dispatcher.db.Update(task)
//
//	if err != nil {
//		log.Println("Error while update the task")
//		return err
//	}
//
//	return nil
//}
//
//// listen the database event channel to do some actions
//func (dispatcher *Dispatcher) listen() error {
//
//	// build the event chan to listen
//	var eventChan = "event_task"
//	//listener := dispatcher.db.Listen("event_task_" + dispatcher.function)
//
//	// get the database listener for this robot
//	listener := dispatcher.db.Listen(eventChan)
//	//listener := dispatcher.db.Listen("event_task_" + dispatcher.function)
//
//	//    defer ln.Close()
//
//	// get the channel
//	channel := listener.Channel()
//
//	// while true
//	for {
//		select {
//		// receive a database event
//		case event := <-channel:
//
//			log.Println("654.0-event-received------------------------------")
//			log.Println(event)
//
//			// model of the event data
//			var notification models.Notification
//
//			// decode event json data
//			err := json.Unmarshal([]byte(event.Payload), &notification)
//
//			if err != nil {
//				log.Println("Error while decode the database notification", err)
//				continue
//			}
//
//			// restart event process
//			if notification.Action == "RESTART" {
//
//				log.Println("Received event for restart the robot")
//
//				// XXX: need to develop this part of restart data
//
//				continue
//			}
//
//			// task event process
//			if notification.Action == "TASK" && notification.Data.TaskID != 0 {
//
//				log.Printf("Received event for task ID : %d", notification.Data.TaskID)
//
//				// push task ID event to the channel
//				dispatcher.queue <- notification.Data.TaskID
//
//				continue
//			}
//
//		// recursive auto-pull
//		//		case <-time.After(90 * time.Second):
//		case <-time.After(10 * time.Second):
//
//			// logger
//			log.Println("Auto-pull")
//
//			// task pull
//			go dispatcher.getTaskIDs()
//		}
//	}
//}

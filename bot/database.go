package bot

import (
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
	var addr = ENV_DB_HOST

	if ENV_DB_PORT != "" {
		addr += ":" + ENV_DB_PORT
	}

	// pg database connector
	db := pg.Connect(&pg.Options{
		Addr:       ENV_DB_HOST + ":" + ENV_DB_PORT,
		User:       ENV_DB_USER,
		Password:   ENV_DB_PASSWORD,
		Database:   ENV_DB_DATABASE,
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
		log.Println("Error while select robots")
		return err
	}

	log.Println("2================")

	// no elemtns, error
	if len(robots) == 0 {
		log.Println("Error no robots definition found")
		return errors.New("No robots definition found")
	}

	log.Println("3================")

	// store the endpoint to fetch
	var stepIDs = make(map[string][]int64)

	log.Println("3.2================")

	// save into dispatcher definitions data
	for _, robot := range robots {

		log.Println("3.3.1.............")
		log.Printf("%+v", robot)

		// remap by version the definitions
		dispatcher.robots[robot.Version] = &robot.Definition

		log.Println("3.3.2.............")

		// save the step IDs to fetch after
		for _, step := range robot.Definition.Sequence {
			stepIDs[step.EndpointType] = append(stepIDs[step.EndpointType], step.EndpointID)
		}
	}

	log.Println("4================")

	// fetch the steps data
	for model, steps := range stepIDs {

		// http case
		if model == "HTTP" {

			// object to fetch
			var endpoints []*models.HttpEndpoint

			// get the endpoint data
			err := dispatcher.db.
				Model(&endpoints).
				Where(models.TblHttpEndpoint_Id+" IN ( ? )", pg.In(steps)).
				Select()

			if err != nil {
				log.Printf("Error while retrieve endpoints : %s", model)
				return err
			}

			// iterate on each endpoints
			for _, endpoint := range endpoints {

				// save endpoints
				dispatcher.endpoints[endpoint.Id] = &endpoint
			}
		}
	}

	log.Println("5================")

	return nil
}

// retrieve the available task ID list
func (dispatcher *Dispatcher) getTaskIDs() error {

	log.Println("Read all task IDs")

	// where store the ID list
	var IDs []int64

	log.Println("9.0=========================")

	// working on this model
	var task models.Task

	// fetch the available ID list
	err := dispatcher.db.
		Model(&task).
		Column("id").
		OrderExpr(models.TblTask_Id+" ASC").
		Where(models.TblTask_Status+" = ?", "TODO").
		Where(models.TblTask_Function+" = ?", dispatcher.function).
		Where(models.TblTask_Retry+" > ?", 0).
		Where(models.TblTask_TodoDate + " <= NOW()").
		Select(&IDs)

	log.Println("9.1=========================")

	if err != nil {
		log.Println("Error while fetching the task ids")
		log.Println(err)
		return err
	}

	log.Println("9.2=========================")

	log.Println("All rows:")

	for x, id := range IDs {

		dispatcher.queue <- id

		log.Printf("9.3=====>  %d : %+v\n", x, id)
	}

	return nil
}

// retrieve a task by ID
func (dispatcher *Dispatcher) getTask(id int64) (*models.Task, error) {

    log.Println("Read task")

	// model to fetch
	var task models.Task

	// fetch the object
	err := dispatcher.db.
		Model(&task).
		Where(models.TblTask_Id+" = ?", id).
		Where(models.TblTask_Status+" = ?", "TODO").
		Where(models.TblTask_Function+" = ?", dispatcher.function).
		Where(models.TblTask_Retry+" > ?", 0).
		Where(models.TblTask_TodoDate + " <= NOW()").
		First()

	if err != nil {
		if err == pg.ErrNoRows {
			log.Println("Task not found")
			return nil, errors.New("Task not found")
		}

		log.Println("Error while retrieve a task")
		return nil, err
	}

    return &task, nil
}

// update a task data
func (dispatcher *Dispatcher) updateTask(task *models.Task) error {

    // last update key
	task.LastUpdate = time.Now()

	// fetch the object
	err := dispatcher.db.Update(task)

	if err != nil {
		log.Println("Error while update the task")
		return err
	}

    return nil
}




















//func (d *Dispatcher) initializeListenerAndListen() {
//
//    _, err := sql.Open("postgres", ConnectionConfiguration)
//
//    if err != nil {
//        panic(err)
//    }
//
//    reportProblem := func(ev pq.ListenerEventType, err error) {
//        if err != nil {
//            log.Println(err.Error())
//        }
//    }
//
//    listener := pq.NewListener(ConnectionConfiguration, 10*time.Second, time.Minute, reportProblem)
//
//    err = listener.Listen("events_task_" + d.Configuration.Function)
//
//    if err != nil {
//        panic(err)
//    }
//
//    log.Println("Start monitoring PostgreSQL...")
//
//    for {
//        d.waitForNotification(listener)
//    }
//}
//
//type DatabaseNotification struct {
//    Table       string
//    Action      string
//    Function    string
//    ID          int64
//}
//
//func (d *Dispatcher) waitForNotification(l *pq.Listener) {
//    for {
//        select {
//            case n := <-l.Notify:
//
//                var notification DatabaseNotification
//
//                err := json.Unmarshal([]byte(n.Extra), &notification)
//
//                if err != nil {
//                    log.Println("error:",err)
//                }
//
//                log.Println("Received data from channel [", n.Channel, "] :")
//
//                log.Printf("%+v \n", notification)
//
//                d.IdQueue <- notification.ID
//
//                log.Println("Data send in task queue")
//
//            case <-time.After(90 * time.Second):
//
//                log.Println("Received no events for 90 seconds, checking connection")
//
//                go l.Ping()
//
//                log.Println("Retreieve ids")
//
//                go d.readTaskIds()
//        }
//    }
//}
//
//func (w *Worker) readOneTask(id int64) (task Task, err error) {
//
//    log.Println("Read one task")
//
//    err = w.connector.SelectOne(
//        &task,
//        "select * from task where status = :status and function = :function and id = :id and retry > 0 and todo_date <= now() limit 1",
//        map[string]interface{}{"status":"todo","function":w.Function,"id":id} )
//
//    if err != nil {
//        log.Fatalln("Select failed", err)
//    }
//
//    return task, err
//}

//import (
//    "log"
//    "time"
//    "encoding/json"
//    "database/sql"
//    "github.com/lib/pq"
//)
//
//func (d *Dispatcher) readSteps() {
//
//    log.Println("Get the dispatcher definition")
//
//    err := d.connector.SelectOne(
//        &d.Definition,
//        "select * from definition where function = :function",
//        map[string]interface{}{"function":d.Configuration.Function} )
//
//    if err != nil {
//        log.Fatalln("Select failed", err)
//    }
//}
//
//func (d *Dispatcher) readTaskIds() {
//
//    log.Println("Read all task IDs")
//
//    var taskIds []int64
//
//    _, err := d.connector.Select(
//        &taskIds,
//        "select id from task where status = :status and function = :function and retry > 0 and todo_date <= now() order by id asc",
//        map[string]interface{}{"status":"todo","function":d.Configuration.Function} )
//
//    if err != nil {
//        log.Fatalln("Select failed", err)
//    }
//
//    log.Println("All rows:")
//
//    for x, id := range taskIds {
//
//        d.IdQueue <- id
//
//        log.Printf("  %d : %v\n", x, id)
//    }
//}
//
//func (d *Dispatcher) initializeListenerAndListen() {
//
//    _, err := sql.Open("postgres", ConnectionConfiguration)
//
//    if err != nil {
//        panic(err)
//    }
//
//    reportProblem := func(ev pq.ListenerEventType, err error) {
//        if err != nil {
//            log.Println(err.Error())
//        }
//    }
//
//    listener := pq.NewListener(ConnectionConfiguration, 10*time.Second, time.Minute, reportProblem)
//
//    err = listener.Listen("events_task_" + d.Configuration.Function)
//
//    if err != nil {
//        panic(err)
//    }
//
//    log.Println("Start monitoring PostgreSQL...")
//
//    for {
//        d.waitForNotification(listener)
//    }
//}
//
//type DatabaseNotification struct {
//    Table       string
//    Action      string
//    Function    string
//    ID          int64
//}
//
//func (d *Dispatcher) waitForNotification(l *pq.Listener) {
//    for {
//        select {
//            case n := <-l.Notify:
//
//                var notification DatabaseNotification
//
//                err := json.Unmarshal([]byte(n.Extra), &notification)
//
//                if err != nil {
//                    log.Println("error:",err)
//                }
//
//                log.Println("Received data from channel [", n.Channel, "] :")
//
//                log.Printf("%+v \n", notification)
//
//                d.IdQueue <- notification.ID
//
//                log.Println("Data send in task queue")
//
//            case <-time.After(90 * time.Second):
//
//                log.Println("Received no events for 90 seconds, checking connection")
//
//                go l.Ping()
//
//                log.Println("Retreieve ids")
//
//                go d.readTaskIds()
//        }
//    }
//}
//
//func (w *Worker) readOneTask(id int64) (task Task, err error) {
//
//    log.Println("Read one task")
//
//    err = w.connector.SelectOne(
//        &task,
//        "select * from task where status = :status and function = :function and id = :id and retry > 0 and todo_date <= now() limit 1",
//        map[string]interface{}{"status":"todo","function":w.Function,"id":id} )
//
//    if err != nil {
//        log.Fatalln("Select failed", err)
//    }
//
//    return task, err
//}
//

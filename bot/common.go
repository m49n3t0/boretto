package bot

import (
	"github.com/go-pg/pg"
	"os"
)

///////////////////////////////////////////////////////////////////////////////

var (
	// robot configuration
	ENV_FUNCTION = os.Getenv("FUNCTION")
	//	ENV_VERSION  = os.Getenv("VERSION")

	// database configuration
	ENV_DB_HOST     = os.Getenv("DB_HOST")
	ENV_DB_PORT     = os.Getenv("DB_PORT")
	ENV_DB_USER     = os.Getenv("DB_USER")
	ENV_DB_PASSWORD = os.Getenv("DB_PASSWORD")
	ENV_DB_DATABASE = os.Getenv("DB_DATABASE")

	// dispatcher configuration
	ENV_MAX_WORKER = 20
	ENV_MAX_QUEUE  = 5
)

//func getDatabaseHandle() *pg.DB {
//	return pg.Connect(&pg.Options{
//		Addr:       ENV_DB_HOST + ":" + ENV_DB_PORT,
//		User:       ENV_DB_USER,
//		Password:   ENV_DB_PASSWORD,
//		Database:   ENV_DB_DATABASE,
//		MaxRetries: 2,
//	})
//}

//var (
//    MaxWorker               = 20        // os.Getenv("MAX_WORKERS")
//    MaxQueue                = 5         // os.Getenv("MAX_QUEUE")
//    MaxLength int64         = 20480
//    ConnectionConfiguration = "postgres://teaas:teaasTEAAS89@ts61115-053.dbaas.ovh.net:35176/teaas"
//)
//
//{
//	// object to fetch
//	var tasks []*credit.Task
//
//	// create query
//	query := rd.WithContext(ctx).
//		Model(&tasks).
//		OrderExpr("id ASC")
//
//	// check optional pagination params
//	query, pErr := common.Paginator(query, params.PerPage, params.Page)
//	if pErr != nil {
//		return nil, 0, pErr
//	}
//
//	// execute query with parameters
//	count, err := query.SelectAndCount()
//
//	if err != nil {
//		return nil, 0, common.RequestError(err)
//	}
//
//	return tasks, int64(count), nil
//}
//
//{
//	// object to fetch
//	var task credit.Task
//
//	// do the query
//	query := rd.WithContext(ctx).
//		Model(&task).
//		Where(credit.TblTask_Id+" = ?", params.TaskID)
//
//	// optional filter
//	if params.CustomerCode != nil && *params.CustomerCode != "" {
//		query = query.Where(credit.TblTask_CustomerCode+" = ?", *params.CustomerCode)
//	}
//
//	// do the query
//	err = query.First()
//
//	if err != nil {
//		if err == pg.ErrNoRows {
//			return nil, models.NewError(http.StatusNotFound, "Task not found")
//		}
//		return nil, common.RequestError(err)
//	}
//
//	return &task, nil
//}

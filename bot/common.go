package bot

import (
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

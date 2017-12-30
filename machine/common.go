package machine

import (
	"os"
)

///////////////////////////////////////////////////////////////////////////////

var (
	// robot configuration
	FUNCTION = os.Getenv("FUNCTION")

	// database configuration
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_DATABASE = os.Getenv("DB_DATABASE")

	// dispatcher configuration
	MAX_WORKER = 20
	MAX_QUEUE  = 5
)

package main

import (
    "os"
    "log"
    "database/sql"
    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"
)

// ================================================= //
// ================================================= //

var (
    ConnectionConfiguration = "postgres://provisionning:totoTOTO89@641a3187-5896-49c9-af7d-d8bed8187f79.pdb.ovh.net:21684/provisionning"
)


func initDb() *gorp.DbMap {
    db, err := sql.Open("postgres", ConnectionConfiguration)
    if err != nil {
        log.Fatalln("sql.Open failed ...", err )
        panic(err)
    }

    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

    // Will log all SQL statements + args as they are run
    // The first arg is a string prefix to prepend to all log messages
    dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))

    // add a table, setting the table name to 'task' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(Worker{}, "workers").SetKeys(true, "Id")

    // drop the tables
    // use a migration tool, or create the tables via scripts
    err = dbmap.DropTablesIfExists()
    if err != nil {
        log.Fatalln("sql.DropTables failed ...", err )
        panic(err)
    }

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    err = dbmap.CreateTablesIfNotExists()
    if err != nil {
        log.Fatalln("sql.CreateTablesIfNotExists failed ...", err )
        panic(err)
    }

    // create a worker
    w := &Worker{
        Name:"aldo",
        Job:"producer",
        Interval:4,
        Produce:"server" }

    err = dbmap.Insert(w)
    if err != nil {
        log.Fatalln("sql.Insert worker failed ...", err )
        panic(err)
    }

    return dbmap
}

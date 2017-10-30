package main

import (
    "log"
    "fmt"
    "time"
)

func main() {

    // OS vars
    var name = "aldo"

    // retrieve the gorp dbmap
    dbmap := initDb()

    // default object
    var w = Worker{}

    // retrieve worker data
    err := dbmap.SelectOne(&w, "select * from workers where name = :name", map[string]interface{}{"name":name})
    if err != nil {
        log.Fatalln("sql.SelectOne worker failed ...", err )
        panic(err)
    }

    // create ticker for work
    ticker := time.NewTicker(time.Second * 4)
    tickerStop := time.NewTicker(time.Second * 10)

    fmt.Println("Worker '", w.Name, "' arrived at work")

    s := true

    for s {

        select {

            case <-ticker.C:
                fmt.Println("Worker '", w.Name, "' produce '1 ", w.Produce, "'")

            case <-tickerStop.C:
                ticker.Stop()
                tickerStop.Stop()
                s = false
        }

    }

    fmt.Println("Worker '", w.Name, "' leave his work")
}
















type Worker struct {
    Id int64 `db:"id"`
    Name string `db:"name"` // name of the person
    Job string `db:"job"` // name of the job
    Interval int `db:"interval"` // how many seconds for do his action
    Produce string `db:"produce"`// this worker produce what ?
}
























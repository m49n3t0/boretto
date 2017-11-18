package main

import (
	"github.com/m49n3t0/teaas/bot"
	"log"
)

// Main function
func main() {

	log.Println("Begin")

    exec, err := bot.New()

    if err != nil {
        panic( err )
    }

    exec.Run()

	log.Println("Done")
}

//package main
//
//import (
//    "log"
// //   "time"
//
//    "github.com/jinzhu/gorm"
//    _ "github.com/lib/pq"
//)
//
//// Task contains informations to work with the robots
//type Task struct {
//    Id int64
//    Domain string
//    Function string
//    Step string
//    Status string
//}
//
//// BotConfig contains informations to parameters the robots
//type BotConfig struct {
//    function string
//    sequence []string
//}
//
//// Bot contains informations for current robots skeleton
//type Bot struct {
//    config BotConfig
//    db gorm.DB
//}
//
//// Build bot struct and initialize all starters
//func New(configuration BotConfig) (*Bot, error) {
//    var bot Bot
//    var err error
//
//    dbmap, err := gorm.Open("postgres", "host=b96a2461-998e-4f9a-8869-b01c43accf20.pdb.ovh.net port=21617 user=my_user dbname=my_data sslmode=require password=azertyOP89")
//
//    if err != nil {
//        log.Println("Error while try to connect to the database", err )
//        return nil, err
//    }
//
//    dbmap.LogMode( true )
//
//    bot = Bot{
//        config: configuration,
//        db: *dbmap,
//    }
//
//    return &bot, err
//}
//
//// Run the executor process
//func (exec *Bot) Run() {
//
//    var tasks []Task
//
//    exec.db.Limit(3).Where("status = ?", "todo").Find(&tasks)
//
//    for _, task := range tasks {
//        log.Printf("working on task : %+T %+v\n", task, task )
//
//        if affected := exec.db.Model(&task).Update("status","doing").RowsAffected; affected != 1 {
//            log.Println("Error while affected a task row", affected)
//            continue
//        }
//
//        for _, step := range exec.config.sequence {
//            log.Printf("id : %d / step : %s\n", task.Id, step)
//        }
//    }
//}
//
//// Main function
//func main() {
//
//    // Build the robot configuration
//    conf := BotConfig{
//        function: "database/delete",
//        sequence: []string{ "starting","checking","onServer","onInterne","ending" },
//    }
//
//    // Build the robot executor
//    exec, err := New( conf )
//
//    if err != nil {
//        log.Fatalln("Error while create the executor robot", err)
//    }
//
//    // Print the informations
//    log.Println("      ---===0O0===---")
//    log.Println( exec )
//    log.Println("      ---===0O0===---")
//
//    exec.Run()
//
//}

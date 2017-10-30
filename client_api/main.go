package main

import (
    "fmt"
    "log"
    "time"
    "net/http"

    "github.com/go-pg/pg"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/loopfz/gadgeto/tonic"
)

// encapsulate the database connection
type Resource struct {
    dbconn pg.DB
}

// retrieve the resource data
func getResource() (*Resource, error){

	dbconn := pg.Connect(&pg.Options{
		User:     "teaas",
		Password: "teaasTEAAS89",
		Database: "teaas",
		Addr:     "ts57878-001.dbaas.ovh.net:35178",
	})

    var n int
    _, err := dbconn.QueryOne(pg.Scan(&n), "SELECT 1")
    if err != nil {
        return nil, err
    }

    return &Resource{dbconn: *dbconn}, nil
}

// clean resource connections and coo
func (r *Resource) clean() {
    err := r.dbconn.Close()
    if err != nil {
        panic(err)
    }
}

// launch api
func main() {

    rsc, err := getResource()
    if err != nil {
        panic(err)
    }
    defer rsc.clean()

    rsc.dbconn.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
        query, err := event.FormattedQuery()
        if err != nil {
            panic(err)
        }

        log.Printf("%s %s", time.Since(event.StartTime), query)
    })

    r := gin.Default()

    r.Use(cors.Default())

    r.GET("/mon/ping", MonPing)

    r.GET("/region", tonic.Handler(rsc.GetRegions, 200))
    r.POST("/region", tonic.Handler(rsc.PostRegion, 200))
    r.GET("/region/:region_id", tonic.Handler(rsc.GetRegion, 200))

    r.Run(":8080")
}

// ========================================================================= //
// ========================================================================= //

// GET #URL#/mon/ping
func MonPing(c *gin.Context) {
    c.JSON( http.StatusOK, gin.H{ "message": "pong" } )
}

// ========================================================================= //
// ========================================================================= //

type Region struct {
    Id int64 `json:"id"`
    Name string `json:"name"`
    RegionId int64 `json:"parent_region_id"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}

func (a Region) String() string {
    return fmt.Sprintf("Region<%d %s %d %s %s>",
        a.Id, a.Name, a.RegionId, a.CreatedAt, a.UpdatedAt)
}

// GET #URL#/region
type GetRegionsInput struct {
    ParentRegionId int64 `query:"parent_region_id"`
}

func (r *Resource) GetRegions(c *gin.Context, in *GetRegionsInput) (*gin.H, error) {
    // define structure
    var regions []Region
    // query
    query := r.dbconn.
        Model(&regions).
        Column("region.*","Childs","Sites")
    // query filter: parent_region_id
    if in.ParentRegionId != 0 {
        query = query.Where("region.region_id = ?", in.ParentRegionId)
    }
    // request
    count, err := query.SelectAndCount()
    if err != nil {
        // log error
        log.Printf("%s %s", time.Now(), err)
        // return error
        return nil, fmt.Errorf("Error while getting regions informations")
    }
    // if no elements, return an empty array
    if count == 0 || len(regions) == 0 {
        return &gin.H{"regions":[]string{},}, nil
    }
    // return result into value
    return &gin.H{"regions":regions,"meta":&gin.H{"total":count,},}, nil
}

// POST #URL#/region
type PostRegionInput struct {
    RegionName string `json:"region_name" binding:"required"`
    ParentRegionId int64 `json:"parent_region_id"`
}

func (r *Resource) PostRegion(c *gin.Context, in *PostRegionInput) (*gin.H, error) {
    // check not null key
    if in.RegionName == "" {
        return nil, fmt.Errorf("Missing region name parameter")
    }
    // define check structure
    var regionCheck Region
    // query
    query := r.dbconn.
        Model(&regionCheck).
        Where("region.name = ?", in.RegionName)
    // request
    count, err := query.Count()
    if err != nil {
        // log error
        log.Printf("%s %s", time.Now(), err)
        // return error
        return nil, fmt.Errorf("Error while creating region")
    }
    // unique name already exists
    if count > 0 {
        return nil, fmt.Errorf("This region name already exists")
    }

    // check parent_id parameter
    if in.ParentRegionId != 0 {
        // query
        query := r.dbconn.
            Model(&regionCheck).
            Where("region.id = ?", in.ParentRegionId)
        // request
        count, err := query.Count()
        if err != nil {
            // log error
            log.Printf("%s %s", time.Now(), err)
            // return error
            return nil, fmt.Errorf("Error while creating region")
        }
        // unique name already exists
        if count != 1 {
            return nil, fmt.Errorf("This region parent id doesn't exists")
        }
    }

    // define structure
    region := &Region{
        Name: in.RegionName,
        RegionId: in.ParentRegionId,
    }
    // insert
    err = r.dbconn.Insert(region)
    if err != nil {
        // log error
        log.Printf("%s %s", time.Now(), err)
        // return error
        return nil, fmt.Errorf("Error while creationg regions")
    }
    // return result into value
    return &gin.H{"region":region,}, nil
}

// GET #URL#/region/:region_id
type GetRegionInput struct {
    RegionId int64 `path:"region_id,required" description:"Region id"`
}

func (r *Resource) GetRegion(c *gin.Context, in *GetRegionInput) (*gin.H, error) {
    // define structure
    var region Region
    // query
    query := r.dbconn.
        Model(&region).
        Column("region.*","Childs","Sites").
        Where("region.id = ?", in.RegionId)
    // request
    count, err := query.SelectAndCount()
    if err != nil {
        // log error
        log.Printf("%s %s", time.Now(), err)
        // return error
        return nil, fmt.Errorf("Error while getting regions informations")
    }
    // if no elements, return a not found error
    if count == 0 {
        return nil, fmt.Errorf("No element found")
    }
    // return result into value
    return &gin.H{"region":region,"meta":&gin.H{"total":count,},}, nil
}


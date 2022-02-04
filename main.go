package main

import (
	"html/template"
	"log"
	"os"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// migrate the db on startup
	m, err := migrate.New("file://sql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	m.Up()

	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	ctrl := &Controller{DB: db}

	router := gin.Default()

	router.Static("/assets", "assets")
	router.HTMLRender = renderer()

	router.GET("/", ctrl.Index)
	router.GET("/new", ctrl.NewJob)
	router.POST("/jobs", ctrl.CreateJob)

	router.Run()
}

func renderer() multitemplate.Renderer {
	funcMap := template.FuncMap{
		"formatAsDate":          formatAsDate,
		"formatAsRfc3339String": formatAsRfc3339String,
	}

	r := multitemplate.NewRenderer()
	r.AddFromFilesFuncs("index", funcMap, "./templates/base.html", "./templates/index.html")
	r.AddFromFiles("new", "./templates/base.html", "./templates/new.html")

	return r
}

package main

import (
	"fmt"
	"go-crawler/args"
	"go-crawler/db"
	"go-crawler/models"
	"go-crawler/redis"
	"go-crawler/routers"
	"go-crawler/tasks"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/judwhite/go-svc/svc"
	"github.com/robfig/cron/v3"
)

type program struct {
	wg   sync.WaitGroup
	quit chan struct{}
}

func (p *program) Init(env svc.Environment) error {
	redis.ConnectRedis()
	return nil
}

func startCronjobs() {
	c := cron.New()
	c.AddFunc("@every 30m", func() {
		tasks.Fetch()
	})
	c.Start()
}

func (p *program) Start() error {
	args.ParseCmd()
	switch args.Cmd.DB {
	case "create":
		fmt.Println("creating database")
		db.Create()
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	case "migrate":
		fmt.Println("migrating tables")
		db.Migrate(args.Cmd.GIN_ENV, &models.Blog{}, &models.Blogger{}, &models.Tag{}, &models.Category{})
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	case "drop":
		fmt.Println("droping database")
		if args.Cmd.TABLE != "" {
			db.Open("")
			db.DB.Migrator().DropTable(args.Cmd.TABLE)
		} else {
			db.Drop()
		}
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	default:
		fmt.Println("server starting...")
		db.Open("")
		//models.GetCategories()
		//tasks.Fetch()
		startCronjobs()
		routers.InitRouter(os.Interrupt)
	}
	return nil
}

func (p *program) Stop() error {
	fmt.Println("\nserver stoping")
	time.Sleep(time.Duration(1) * time.Second)
	return nil
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, os.Interrupt); err != nil {
		log.Fatal(err)
	}
}

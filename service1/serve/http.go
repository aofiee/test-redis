package serve

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"taobin-service/configs"
	"taobin-service/internal/handlers"
	"taobin-service/internal/repositories"
	"taobin-service/internal/services"
	"taobin-service/utils/database/gorm"
	"taobin-service/utils/database/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type config struct {
	Env string
}

func ServeHTTP() error {
	var (
		err error
		cfg config
	)
	app := fiber.New()
	flag.StringVar(&cfg.Env, "env", "", "the environment to use")
	flag.Parse()
	configs.InitConfig("./configs")
	logrus.Info(configs.GetConfig())

	dbConGorm, err := gorm.Connect2Postgres(
		configs.GetConfig().Postgres.Host,
		configs.GetConfig().Postgres.Port,
		configs.GetConfig().Postgres.Username,
		configs.GetConfig().Postgres.Password,
		configs.GetConfig().Postgres.DbName,
		configs.GetConfig().Postgres.SSLMode,
	)
	if err != nil {
		return err
	}

	redisCon, err := redis.Connect2Redis(
		configs.GetConfig().Redis.Host,
		configs.GetConfig().Redis.Port,
		configs.GetConfig().Redis.Password,
	)
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logrus.Info("shutting down ...")
			gorm.DisconnectPostgres(dbConGorm.Postgres)
			redis.DisconnectRedis(redisCon.Redis)
			app.Shutdown()
		}
	}()

	postgresRepo := repositories.NewPostgres(dbConGorm.Postgres, redisCon.Redis)
	srv := services.New(postgresRepo, redisCon.Redis)
	hdl := handlers.New(srv, dbConGorm.Postgres)

	app.Get("/test", hdl.TestCheck)
	taobin := app.Group("/v1/api/machine")
	{
		taobin.Post("/", hdl.CreateMachine)
		taobin.Put("/", hdl.UpdateMachine)
		taobin.Delete("/:id", hdl.DeleteMachine)
		taobin.Get("/:id", hdl.GetMachines)
		taobin.Get("/", hdl.GetMachines)
	}

	err = app.Listen(":" + configs.GetConfig().App.Port)
	if err != nil {
		return err
	}

	fmt.Println("Listening on port: ", configs.GetConfig().App.Port)
	return nil

}

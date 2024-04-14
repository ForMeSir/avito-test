package integration_test

// import (
// 	"testing"
// 	//_ "github.com/lib/pq"
// )

import (
	"avito/internal/cache"
	"avito/internal/handler"
	"avito/internal/repository"
	"avito/internal/service"
	"avito/tests/integration"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB


func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
	}


	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "latest",
			Env: []string{
					"POSTGRES_PASSWORD=test",
					"POSTGRES_USER=test",
					"POSTGRES_DB=DB",
					"listen_addresses = '*'",
			},
	}, func(config *docker.HostConfig) {

			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
			log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://test:test@%s/DB?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120)

	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
			db, err = sqlx.Open("postgres", databaseUrl)
			if err != nil {
					return err
			}
			return db.Ping()
	}); err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
	}

  schema, err := os.ReadFile("20240407122642_init.up.sql")
			if err != nil {
					log.Fatalf("Could not read schema: %s", err)
			}
			if _, err := db.Exec(string(schema)); err != nil {
					log.Fatalf("Could not apply schema: %s", err)
			}

			

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
	}


	os.Exit(code)
}

	
func Test_1(t *testing.T) {
	
simple:=cache.NewSimple()
  cache:= cache.NewCache(simple)
  repos:= repository.NewRepository(db)
  services:=service.NewService(repos, cache)
  handlers:= handler.NewHandler(services)

router := gin.Default()

router.GET("/user_banner",handlers.BannerForUser)
integration.TestOneUserBanner(t,router)
integration.TestTwoUserBanner(t,router)
integration.TestThreeUserBanner(t,router)
integration.TestFourUserBanner(t,router,db)

}


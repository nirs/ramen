package e2e

import (
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

var cleanup = flag.Bool("cleanup", true, "Cleanup app on errors")
var appCount atomic.Uint32

type App struct {
	namespace string
	deployed  bool
}

func NewApp() *App {
	return &App{
		namespace: fmt.Sprintf("app-%d", appCount.Add(1)),
		deployed:  true,
	}
}

func (app App) String() string {
	return app.namespace
}

func (app *App) Deploy(t *testing.T) {
	log.Printf("[%s] Deploying %s", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Waiting until drpc for %s is deployed", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Waiting until %s replicated", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Deploy %s finished successfuly", t.Name(), app)
}

func (app *App) Failover(t *testing.T) {
	log.Printf("[%s] Failing over %s", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Waiting until %s failes over", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Waiting until %s is replicated", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Failover %s finished successfuly", t.Name(), app)
}

func (app *App) Relocate(t *testing.T) {
	log.Printf("[%s] Relocating %s", t.Name(), app)
	t.Fatalf("Testing failed subtest")
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Waiting until %s relocates", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Waiting until %s is replicated", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Relocate %s finished successfuly", t.Name(), app)
}

func (app *App) Undeploy(t *testing.T) {
	log.Printf("[%s] Undeploying %s", t.Name(), app)
	time.Sleep(1 * time.Second)
	log.Printf("[%s] Undeploy %s finished successfully", t.Name(), app)
	app.deployed = false
}

func (app *App) Cleanup(t *testing.T) {
	if app.deployed && *cleanup {
		app.Undeploy(t)
	}
}

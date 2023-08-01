package e2e_test

import (
	"testing"

	"github.com/red-hat-storage/ramen/e2e"
)

func TestBasicDeploy(t *testing.T) {
	t.Parallel()
	app := e2e.NewApp()
	defer app.Cleanup(t)

	if !t.Run("deploy", app.Deploy) {
		return
	}

	t.Run("undeploy", app.Undeploy)
}

func TestBasicLifeCycle(t *testing.T) {
	t.Parallel()
	app := e2e.NewApp()
	defer app.Cleanup(t)

	if !t.Run("deploy", app.Deploy) {
		return
	}

	if !t.Run("failover", app.Failover) {
		return
	}

	if !t.Run("relocate", app.Relocate) {
		return
	}

	t.Run("undeploy", app.Undeploy)
}

func TestBasicFailoverMany(t *testing.T) {
	t.Parallel()
	app := e2e.NewApp()
	defer app.Cleanup(t)

	if !t.Run("deploy", app.Deploy) {
		return
	}

	if !t.Run("failover1", app.Failover) {
		return
	}

	if !t.Run("failover2", app.Failover) {
		return
	}

	t.Run("undeploy", app.Undeploy)
}

func TestBasicRelocateMany(t *testing.T) {
	t.Parallel()
	app := e2e.NewApp()
	defer app.Cleanup(t)

	if !t.Run("deploy", app.Deploy) {
		return
	}

	if !t.Run("relocate1", app.Relocate) {
		return
	}

	if !t.Run("relocate2", app.Relocate) {
		return
	}

	t.Run("undeploy", app.Undeploy)
}

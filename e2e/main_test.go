package e2e_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/red-hat-storage/ramen/e2e"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

var configFile = "./config/config.yaml"

func init() {
	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix("RAMENE2E")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func TestMain(t *testing.T) {
	config := &e2e.Config{}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read configuration file: %v\n", err)
		os.Exit(1)
	}
	if err := viper.UnmarshalExact(config); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse configuration: %v\n", err)
		os.Exit(1)
	}

	config.Validate()

	testContext, err := e2e.NewContext(config)
	if err != nil {
		fmt.Printf("Failed to create TestContext: %v\n", err)
		os.Exit(1)
	}

	defer testContext.Cleanup()

	suite.Run(t, &InfraValidationSuite{testContext: testContext})

	if config.SetupRegionalDRonOCP {
		suite.Run(t, &RegionalDRSuite{testContext: testContext})
	}
}

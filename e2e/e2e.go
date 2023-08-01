package e2e

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Clusters map[string]struct {
		KubeconfigPath string `mapstructure:"kubeconfigpath" required:"true"`
	} `mapstructure:"clusters" required:"true"`
	SetupRegionalDRonOCP bool `mapstructure:"setupregionaldronocp" required:"false"`
}

func (config *Config) Validate() {
	if config.Clusters["hub"].KubeconfigPath == "" {
		fmt.Fprintf(os.Stderr, "Failed to find hub cluster in configuration\n")
		os.Exit(1)
	}

	if config.Clusters["c1"].KubeconfigPath == "" {
		fmt.Fprintf(os.Stderr, "Failed to find c1 cluster in configuration\n")
		os.Exit(1)
	}

	if config.Clusters["c2"].KubeconfigPath == "" {
		fmt.Fprintf(os.Stderr, "Failed to find c2 cluster in configuration\n")
		os.Exit(1)
	}
}

type Context struct {
	Config   *Config
	Clusters map[string]struct {
		k8sClientSet kubernetes.Interface
	}
}

func NewContext(config *Config) (*Context, error) {
	testContext := &Context{
		Config: config,
		Clusters: make(map[string]struct {
			k8sClientSet kubernetes.Interface
		}),
	}
	for clusterName, cluster := range config.Clusters {
		k8sClientSet, err := getClientSetFromKubeConfigPath(cluster.KubeconfigPath)
		if err != nil {
			return nil, err
		}

		testContext.Clusters[clusterName] = struct {
			k8sClientSet kubernetes.Interface
		}{
			k8sClientSet: k8sClientSet,
		}
	}
	return testContext, nil
}

func (ctx *Context) Cleanup() {
	// Add cleanup logic here if needed
}

func (ctx *Context) HubClient() kubernetes.Interface {
	return ctx.Clusters["hub"].k8sClientSet
}

func (ctx *Context) C1Client() kubernetes.Interface {
	return ctx.Clusters["c1"].k8sClientSet
}

func (ctx *Context) C2Client() kubernetes.Interface {
	return ctx.Clusters["c2"].k8sClientSet
}

func getClientSetFromKubeConfigPath(kubeconfigPath string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	k8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k8sClientSet, nil
}

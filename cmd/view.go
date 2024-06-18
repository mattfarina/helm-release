package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

var (
	settings *EnvSettings
)

func newViewCmd(log DebugLog, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a named Helm release",
		Long:  `View the details of a Helm release to stdout as YAML`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleView(log)
		},
	}

	flags := cmd.PersistentFlags()
	flags.Parse(args)
	settings = new(EnvSettings)

	// The Helm plugin system sends the kubeconfig details as an
	// environment variable. Make sure to pick that up for comms
	// to the cluster. The kubecontext is set as an environment variable
	// while the kubeconfig flag is set by Helm to that location.
	if ctx := os.Getenv("HELM_KUBECONTEXT"); ctx != "" {
		settings.KubeContext = ctx
	}

	settings.AddFlags(flags)

	return cmd
}

func handleView(log DebugLog) error {
	// make sure logging is setup for debugging. make it pretty

	// Fetch release

	// display as yaml

	helmDriver := os.Getenv("HELM_DRIVER")
	fmt.Println("helmDriver: ", helmDriver)
	store := setupStorage(log)

	fmt.Printf("store: %+v\n", store)
	return nil

}

// setupStorage provides a means to work with the various types of Helm storage.
// This sets up the storage system to communicate with the proper environment.
func setupStorage(log DebugLog) *storage.Storage {
	helmDriver := os.Getenv("HELM_DRIVER")

	helmsettings := cli.New()

	kc := kube.New(helmsettings.RESTClientGetter())
	kc.Log = log

	lazyClient := &lazyClient{
		namespace: settings.Namespace,
		clientFn:  kc.Factory.KubernetesClientSet,
	}

	var store *storage.Storage
	switch helmDriver {
	case "secret", "secrets", "":
		d := driver.NewSecrets(newSecretClient(lazyClient))
		d.Log = log
		store = storage.Init(d)
	case "configmap", "configmaps":
		d := driver.NewConfigMaps(newConfigMapClient(lazyClient))
		d.Log = log
		store = storage.Init(d)
	case "memory":
		var d *driver.Memory
		if store != nil {
			if mem, ok := store.Driver.(*driver.Memory); ok {
				// This function can be called more than once (e.g., helm list --all-namespaces).
				// If a memory driver was already initialized, re-use it but set the possibly new namespace.
				// We re-use it in case some releases where already created in the existing memory driver.
				d = mem
			}
		}
		if d == nil {
			d = driver.NewMemory()
		}
		d.SetNamespace(settings.Namespace)
		store = storage.Init(d)
	case "sql":
		d, err := driver.NewSQL(
			os.Getenv("HELM_DRIVER_SQL_CONNECTION_STRING"),
			log,
			settings.Namespace,
		)
		if err != nil {
			panic(fmt.Sprintf("Unable to instantiate SQL driver: %v", err))
		}
		store = storage.Init(d)
	default:
		// TODO: real error handling instead of a panic
		panic("Unknown driver in HELM_DRIVER: " + helmDriver)
	}

	return store
}

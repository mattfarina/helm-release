package cmd

import (
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/storage"
)

type DebugLog func(format string, v ...interface{})

type EnvSettings struct {
	KubeConfigFile string
	KubeContext    string
	Namespace      string
	Releases       *storage.Storage
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.KubeConfigFile, "kubeconfig", "", "path to the kubeconfig file")
	fs.StringVar(&s.KubeContext, "kube-context", s.KubeContext, "name of the kubeconfig context to use")
	fs.StringVar(&s.Namespace, "namespace", s.Namespace, "namespace scope of the release")
}

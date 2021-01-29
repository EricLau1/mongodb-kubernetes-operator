package e2eutil

import (
	"fmt"
	"os"
	"testing"

	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	mdbv1 "github.com/mongodb/mongodb-kubernetes-operator/pkg/apis"
)

var TestClient client.Client
var OperatorNamespace string
var WatchNamespace string

type CleanupOptions struct {
	TestContext *Context
}

func (*CleanupOptions) ApplyToCreate(*client.CreateOptions) {}

type Context struct{}

func TestMainEntry(m *testing.M) {
	var cfg *rest.Config
	var testEnv *envtest.Environment
	var err error

	useExistingCluster := true
	testEnv = &envtest.Environment{
		UseExistingCluster:       &useExistingCluster,
		AttachControlPlaneOutput: true,
		KubeAPIServerFlags:       []string{"--authorization-mode=RBAC"},
	}

	fmt.Println("Starting test environment")
	cfg, err = testEnv.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = mdbv1.AddToScheme(scheme.Scheme)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	TestClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Starting test")
	code := m.Run()

	err = testEnv.Stop()
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func NewContext(t *testing.T) *Context {
	return &Context{}
}

func (ctx *Context) Cleanup() {}

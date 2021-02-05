package e2eutil

import (
	"fmt"
	"testing"

	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	mdbv1 "github.com/mongodb/mongodb-kubernetes-operator/pkg/apis"
)

var TestClient client.Client
var OperatorNamespace string

type CleanupOptions struct {
	TestContext *Context
}

func (*CleanupOptions) ApplyToCreate(*client.CreateOptions) {}

type Context struct{}

func RunTest(m *testing.M) (int, error) {
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
		return 1, err
	}

	err = mdbv1.AddToScheme(scheme.Scheme)
	if err != nil {
		return 1, err
	}

	TestClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return 1, err
	}

	fmt.Println("Starting test")
	code := m.Run()

	err = testEnv.Stop()
	if err != nil {
		return code, err
	}

	return code, nil
}

func NewContext(t *testing.T) *Context {
	return &Context{}
}

func (ctx *Context) Cleanup() {}

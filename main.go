package main

import (
	"os"
	"time"

	"github.com/golang/glog"
	flag "github.com/spf13/pflag"

	"github.com/openshift/origin/pkg/util/proc"
	"k8s.io/kubernetes/pkg/client/unversioned"
	kubectl_util "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/util"
)

var (
	flags = flag.NewFlagSet("", flag.ContinueOnError)

	cluster = flags.Bool("use-kubernetes-cluster-service", true, `If true, use the built in kubernetes
        cluster for creating the client`)

	argDomain = flag.String("domain", "cluster.local", "domain under which to create names")

	backend = flag.String("backend", "nsd", "DNS server to use as backend")

	logLevel = flags.Int("v", 1, `verbose output`)
)

func main() {
	clientConfig := kubectl_util.DefaultClientConfig(flags)
	flags.Parse(os.Args)

	var err error
	var kubeClient *unversioned.Client

	if *cluster {
		if kubeClient, err = unversioned.NewInCluster(); err != nil {
			glog.Fatalf("Failed to create client: %v", err)
		}
	} else {
		config, err := clientConfig.ClientConfig()
		if err != nil {
			glog.Fatalf("error connecting to the client: %v", err)
		}
		kubeClient, err = unversioned.New(config)
	}

	namespace, specified, err := clientConfig.Namespace()
	if err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}

	if !specified {
		namespace = ""
	}

	domain := *argDomain
	if !strings.HasSuffix(domain, ".") {
		domain = fmt.Sprintf("%s.", domain)
	}

	
	ks := kube2sky{
		domain: domain,
		backend: 
	}

	ks.endpointsStore = watchEndpoints(kubeClient, &ks)
	ks.servicesStore = watchForServices(kubeClient, &ks)
	ks.podsStore = watchPods(kubeClient, &ks)

	select {}
}

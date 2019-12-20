package main

import (
	"log"
	"os"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Printf("POD CREATED: %s/%s", pod.Namespace, pod.Name)
}

func podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	log.Printf(
		"POD UPDATED. %s/%s %s",
		oldPod.Namespace, oldPod.Name, newPod.Status.Phase,
	)
}

func podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Printf("POD DELETED: %s/%s", pod.Namespace, pod.Name)
}

func main() {
	// Read kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		// If run from inside the container
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create Informer factory
	informerFactory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)

	// Add Pod informer to informer factory
	// Add handler methods to the pod informer
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: podUpdate,
			// Called on resource deletion.
			DeleteFunc: podDelete,
		},
	)

	// Starts all the shared informers that have been created by the factory so
	// far.
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, podInformer.Informer().HasSynced) {
		log.Fatal("Failed to sync")
	}

	<-stopCh
}

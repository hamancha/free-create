package podgetter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func GetClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err.Error())
	}

	return clientset
}

// GetPods returns pod information
func GetPods(w http.ResponseWriter) string {
	client := GetClient()

	// gets pods using selectors
	pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: "app=free-create-app"})
	if err != nil {
		log.Println(err.Error())
	}

	// handling pod information
	for i, pod := range pods.Items {
		podstatusPhase := string(pod.Status.Phase)
		podCreationTime := pod.GetCreationTimestamp()
		age := time.Since(podCreationTime.Time).Round(time.Second)
		labels := pod.GetLabels()

		podInfo := fmt.Sprintf("[%d] Pod: %s, Phase: %s , Created: %s, Age: %s", i, pod.GetName(), podstatusPhase, podCreationTime, age.String())
		fmt.Fprintln(w, podInfo)
		fmt.Fprintln(w, "The labels associated with the pods are:")
		for key, value := range labels {
			podLabels := fmt.Sprintf("%s - %s", key, value)
			fmt.Fprintln(w, podLabels)
		}
	}

	return fmt.Sprintf("There are %d pods in the cluster", len(pods.Items))
}

func WatchPods() {
	client := GetClient()
	watch, err := client.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{LabelSelector: "app=free-create-app"})
	if err != nil {
		log.Println(err.Error())
	}

	for event := range watch.ResultChan() {
		p, ok := event.Object.(*v1.Pod)
		if !ok {
			log.Fatal("unexpected type")
		}
		log.Printf("An event of type: %v occuered on pod: %s\n", event.Type, p.GetName())
		log.Printf("The current status of the pod is %s\n", p.Status.Phase)
	}
}

func WatchPodsWithEvents() {
	client := GetClient()

	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onAdd(obj interface{}) {
	panic("Not implemented")
}

func onUpdate(oldObj, newObj interface{}) {
	panic("Not implemented")
}

func onDelete(obj interface{}) {
	panic("Not implemented")
}

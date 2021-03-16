package podgetter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GetPods returns pod information
func GetPods(w http.ResponseWriter) string {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err.Error())
	}
	// gets pods using selectors
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: "app=free-create-app"})
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

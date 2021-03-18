package podgetter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Application struct {
	name           string
	description    string
	commit         string
	pod_count      int
	oldest_pod_age string
	status         string
}

type Cluster struct {
	name         string
	description  string
	applications map[string]Application
}

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

func buildPodInformation(clusterName string, clusterDescription string, pods *v1.PodList) Cluster {
	var clusterInfo Cluster
	clusterInfo.name = clusterName
	clusterInfo.description = clusterDescription
	for _, pod := range pods.Items {
		podStatusPhase := string(pod.Status.Phase)
		podCreationTime := pod.GetCreationTimestamp()
		podAge := time.Since(podCreationTime.Time).Round(time.Second)
		labels := pod.GetLabels()
		if application, ok := clusterInfo.applications[labels["playfab_application"]]; ok {
			application.pod_count++
			application.commit = labels["pf-main/commit"]
			if podStatusPhase != application.status && podStatusPhase != "Running" {
				application.status = podStatusPhase
			}
			currentOldestDuration, _ := time.ParseDuration(application.oldest_pod_age)
			if podAge > currentOldestDuration {
				application.oldest_pod_age = podAge.String()
			}
		} else {
			unRegisteredApp := Application{
				name:           labels["playfab_application"],
				commit:         labels["pf-main/commit"],
				pod_count:      1,
				oldest_pod_age: podAge.String(),
				status:         podStatusPhase,
			}
			clusterInfo.applications[labels["playfab_application"]] = unRegisteredApp
		}
	}

	return clusterInfo
}

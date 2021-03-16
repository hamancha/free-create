package main

import (
	"fmt"
	"free-create/podgetter"
	"net/http"
)

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/getpods", GetPods)
	http.ListenAndServe("0.0.0.0:9090", handler)
}

// GetPods returns the pods running in the cluster
func GetPods(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Trying to get the pods...")
	data := podgetter.GetPods(w)
	fmt.Fprintln(w, data)
}

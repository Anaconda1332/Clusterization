package alg

import (
	"fmt"
	"time"
)

func Go_simple_clustering(path string, clusters_number int) ([]Point_cluster, time.Duration) {
	var data *[][]float64
	var dim, points_number int
	var err error
	clusters := []Point_cluster{
		{Point: []float64{100, 200}, Cluster: 1},
	}
	data, dim, points_number, err = read(path)
	start_time := time.Now()
	if err == nil {
		clusters = clustering(data, points_number, clusters_number, dim)
	} else {
		fmt.Println("File error")
		fmt.Println(err)
	}
	time_res := time.Since(start_time)
	return clusters, time_res
}

func Go_threaded_clustering(path string, clusters_number int, flow_number int) ([]Point_cluster, time.Duration) {
	var data *[][]float64
	var dim, points_number int
	var err error
	clusters := []Point_cluster{
		{Point: []float64{100, 200}, Cluster: 1},
	}
	data, dim, points_number, err = read(path)
	start_time := time.Now()
	if err == nil {
		clusters = clustering2(data, points_number, clusters_number, dim, flow_number)
	} else {
		fmt.Println("File error")
		fmt.Println(err)
	}
	time_res := time.Since(start_time)
	return clusters, time_res
}

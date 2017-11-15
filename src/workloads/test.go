package workloads

// package main
//
// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
// )
//
// func main() {
// 	var fileDataSize []int
// 	r := rand.New(rand.NewSource(int64(time.Now().Second())))
// 	zipf := rand.NewZipf(r, 2.7, 10485760/10, 10485760) //10485760KB = 10 * 1024 * 1024 = 10G
//
// 	for i := int64(0); i < 100000; i++ {
// 		item := int(zipf.Uint64())
// 		fileDataSize = append(fileDataSize, item)
// 	}
//
// 	for i := 0; i < len(fileDataSize); i++ {
// 		fmt.Println(fileDataSize[i])
// 	}
// }

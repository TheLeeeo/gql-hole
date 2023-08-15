package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"time"
// )

// func main() {
// 	jsonData := map[string]string{
// 		"query": `
// 			{
// 				dailyQoute {
// 					content,
// 					author
// 				}
// 			}
// 		`,
// 	}
// 	jsonValue, _ := json.Marshal(jsonData)
// 	request, err := http.NewRequest("POST", "http://0.0.0.0:6001/graphql", bytes.NewBuffer(jsonValue))
// 	if err != nil {
// 		fmt.Println("error in creating request: ", err)
// 		return
// 	}
// 	request.Header.Set("Content-Type", "application/json")
// 	client := &http.Client{Timeout: time.Second * 10}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		fmt.Println("error in sending request: ", err)
// 		return
// 	}
// 	defer response.Body.Close()
// 	data, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		fmt.Println("error in reading response: ", err)
// 		return
// 	}
// 	fmt.Println(string(data))
// }

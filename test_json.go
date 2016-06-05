package main

import (
  "fmt"
  "encoding/json"
  
)

func main() {
  fmt.Println("Hello, playground")

  data := []byte("[{'field1':'123'}, {'field2':'123'}]")
  fmt.Println(data)
  
  var newDeviceReading []interface{}
  
  json.Unmarshal(data, &newDeviceReading )
  fmt.Println(len(newDeviceReading))

  var deviceReadings []interface{}
  deviceReadings = append(deviceReadings, newDeviceReading...)
  fmt.Println(len(deviceReadings))
  
  jsonAsBytes,_ := json.Marshal(deviceReadings)
  fmt.Println(string(jsonAsBytes))

  var empty1 []string
  var empty2 []interface{}
  empty1 = append(empty1, "test")
  jsonAsBytes1,_ := json.Marshal(empty1)
  jsonAsBytes2,_ := json.Marshal(empty2)
  fmt.Println(string(jsonAsBytes1))
  fmt.Println(string(jsonAsBytes2))


}

package main

import (
  "errors"
  "fmt"
  "encoding/json"
  "github.com/hyperledger/fabric/core/chaincode/shim"
)

type PrintoCentChaincode struct {
}

var devicesIndexStr = "_devicesindex"       //name for the key/value that will store a list of all known marbles

type Device struct{
  DeviceId string `json:"name"`         //the fieldtags are needed to keep case from bouncing around
  Color string `json:"color"`
  Size int `json:"size"`
  User string `json:"user"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
  err := shim.Start(new(PrintoCentChaincode))
  if err != nil {
    fmt.Printf("Error starting PrintoCent chaincode: %s", err)
  }
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *PrintoCentChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
  var err error

  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting welcome message!")
  }

  // Initialize the chaincode
  err = stub.PutState("HILLA_world", []byte(args[0]))
  if err != nil {
      return nil, err
  }
  
  var empty []string
  jsonAsBytes, _ := json.Marshal(empty)               //marshal an emtpy array to clear the index
  err = stub.PutState(devicesIndexStr, jsonAsBytes)
  if err != nil {
    return nil, err
  }
  
  return nil, nil
}

// ============================================================================================================================
// Run - Our entry point
// ============================================================================================================================
func (t *PrintoCentChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
  fmt.Println("run is running " + function)

  // Handle different functions
  if function == "init" {                         //initialize the chaincode state, used as reset
    return t.Init(stub, "init", args)
  } else if function == "write" {                     //writes a value to the chaincode state
    return t.Write(stub, args)
  } else if function == "init_device" {                 //create a new marble
    return t.init_device(stub, args)
  } else if function == "add_reading" {
    return t.add_reading(stub, args)
  }
  fmt.Println("run did not find func: " + function)           //error

  return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *PrintoCentChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
  fmt.Println("query is running " + function)

  // Handle different functions
  if function == "get_state" {                         //read a variable
    return t.get_state(stub, args)
  } else if function == "read" {
    return t.get_state(stub, args)
  }
  fmt.Println("query did not find func: " + function)           //error

  return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *PrintoCentChaincode) get_state(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
  var name, jsonResp string
  var err error

  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
  }

  name = args[0]
  valAsbytes, err := stub.GetState(name)                  //get the var from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
    return nil, errors.New(jsonResp)
  }

  return valAsbytes, nil                          //send it onward
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *PrintoCentChaincode) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
  var name, value string // Entities
  var err error
  fmt.Println("running write()")

  if len(args) != 2 {
    return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
  }

  name = args[0]                              //rename for funsies
  value = args[1]
  err = stub.PutState(name, []byte(value))                //write the variable into the chaincode state
  if err != nil {
    return nil, err
  }
  return nil, nil
}

// ============================================================================================================================
// Init Marble - create a new device, store into chaincode state
// ============================================================================================================================
func (t *PrintoCentChaincode) init_device(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
  var err error

  //     0
  // "device id"
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting 1")
  }

  fmt.Println("- start init device")

  // Initialize a state of the device with empty readings
  deviceid:= args[0]
  var empty []interface{}
  jsonAsBytes, _ := json.Marshal(empty)               //marshal an emtpy array of strings to clear the index
  err = stub.PutState(deviceid, jsonAsBytes)
  if err != nil {
    return nil, err
  }
  
  // get the devices index
  devicesAsBytes, err := stub.GetState(devicesIndexStr)
  if err != nil {
    return nil, errors.New("Failed to get devices index")
  }
  var devicesIndex []string
  json.Unmarshal(devicesAsBytes, &devicesIndex)              //un stringify it aka JSON.parse()
  
  //append
  devicesIndex = append(devicesIndex, deviceid)               //add device name to index list
  fmt.Println("! device index: ", devicesIndex)
  jsonAsBytes, _ = json.Marshal(devicesIndex)
  err = stub.PutState(devicesIndexStr, jsonAsBytes)           //store name of marble

  fmt.Println("- end init device")
  return nil, nil
}

// ============================================================================================================================
// Add device reading
// ============================================================================================================================
func (t *PrintoCentChaincode) add_reading(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
  var err error
  
  //   0              1
  // "deviceid", [r0, r1 ... rN]
  if len(args) < 2 {
    return nil, errors.New("Incorrect number of arguments. Expecting 2")
  }

  if len(args[0]) <= 0 {
    return nil, errors.New("1st argument must be a non-empty string")
  }

  if len(args[1]) <= 0 {
    return nil, errors.New("2nd argument must be an array of readings")
  }

  fmt.Println("- start add reading")
  deviceid := args[0]
  currentReadingsAsBytes, err := stub.GetState(deviceid)
  if err != nil {
    return nil, errors.New("Failed to get old readings")
  }
  
  var currentDeviceReadings []interface{}
  json.Unmarshal(currentReadingsAsBytes, &currentDeviceReadings)  //un stringify it aka JSON.parse()

  var newDeviceReadings []interface{}
  newReadingsAsBytes := []byte(args[1])
  err = json.Unmarshal(newReadingsAsBytes, &newDeviceReadings)   //un stringify it aka JSON.parse()
  if err != nil {
    return nil, errors.New("Failed to parse new readings")
  }

  currentDeviceReadings = append(currentDeviceReadings, newDeviceReadings...)
  
  jsonAsBytes, _ := json.Marshal(currentDeviceReadings)
  err = stub.PutState(deviceid, jsonAsBytes)                      //rewrite the device readings
  if err != nil {
    return nil, err
  }
  
  fmt.Println("- end add reading")
  return nil, nil
}

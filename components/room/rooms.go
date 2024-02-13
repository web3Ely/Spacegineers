package main

func main() {
	controllerAddress := "localhost:53000"
	ctxIdentifier := "compcode"
	componentId := []string{"R0", "R1"}
	componentName := []string{"Room0", "Room1"}

	for index, id := range componentId {
		go NewRoom(controllerAddress, ctxIdentifier, id, componentName[index]).StartServer()
	}

	select {}
}

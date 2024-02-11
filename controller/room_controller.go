/*
Room Controller provide pub/sub pattern for Rooms to send and receive information to and from other Rooms.
Room Controller is a gRPC server and each Room is its gRPC client. It provides bot simple RPC and server-side
streamming function calls for a Room.
*/
package main

import (
	helper "spacegineers_context/controller/utils"
	shipSetting "spacegineers_context/settings"
)

// Controller Server entry point
func main() {
	controller := helper.NewController()
	controller.Address = "localhost:53000"
	controller.ComponentIdentifier = "compcode"
	controller.ConnectionTable = shipSetting.GerShipSettings()
	controller.StartGRPCServer()
}

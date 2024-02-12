/*
This package provides functionalities of reading a yaml file
{Right now we set the value without reading any file}
*/
package ship_setting

func GerShipSettings() map[string][]string {
	// currentlly only returns componentId(string):connectionTable([]string)
	items := make(map[string][]string)
	genrator := []string{"R0"}
	room1 := []string{"G0", "R1"}
	room2 := []string{"R0"}
	items["G0"] = genrator
	items["R0"] = room1
	items["R1"] = room2
	return items
}

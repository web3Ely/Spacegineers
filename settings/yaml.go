/*
This package provides functionalities of reading a yaml file
{Right now we set the value without reading any file}
*/
package fromyaml

func GerShipSettings() map[int32][]int32 {
	items := make(map[int32][]int32)
	room1 := []int32{2}
	room2 := []int32{1, 3}
	room3 := []int32{2}
	items[1] = room1
	items[2] = room2
	items[3] = room3
	return items
}

package ship_setting

/*
Supply type represents incoming and outgoing grpc room resource data. It follows the standard defined in the room protobuf.
  - senderId: is a string that represents a component's identifier, which corresponds to the componentId in component type
  - resouceType: is a string that represents the type of a grpc room resource, it could be 'electricity' or 'air'
  - resouceAvailability: is a boolean that represents the availability of the resource
*/
type Supply struct {
	senderId            string
	resourceType        string
	resouceAvailability bool
}

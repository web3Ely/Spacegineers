package controllerUtil

import (
	"slices"
)

type ComponentManager struct {
	neighbourList []string
	receptionList map[string][]string
}

func NewComponentManager(neighbourList []string) *ComponentManager {
	return &ComponentManager{
		neighbourList: neighbourList,
		receptionList: map[string][]string{},
	}
}

func (m *ComponentManager) NotInReceptList(componentId string, recourceType string) bool {
	return !slices.Contains(m.receptionList[recourceType], componentId)
}

func (m *ComponentManager) AddReceptId(componentId string, resourceType string) {
	if m.NotInReceptList(componentId, resourceType) {
		m.receptionList[resourceType] = append(m.receptionList[resourceType], componentId)
	}
}

/*
func (m *Manager) GetReceptList(resourceType string) []string {
	return m.receptionList[resourceType]
}
fmt.Printf("Room Controller : Component %s receptionList for source type %s is %v \n", componentId, recSupply.ResourceType, component.GetReceptList(recSupply.ResourceType))
*/

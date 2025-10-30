package monitor

import "fmt"

type ResourceInfo struct {
	subscriptionId string
	config         Resource
}

func NewResourceInfo(subscriptionId string, config Resource) *ResourceInfo {
	return &ResourceInfo{
		subscriptionId: subscriptionId,
		config:         config,
	}
}

func (r ResourceInfo) ID() string {
	return r.config.ID
}

func (r ResourceInfo) Name() string {
	return r.config.Name
}

func (r ResourceInfo) Type() string {
	return r.config.Type
}

func (r ResourceInfo) Path() string {

	idString := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s", r.subscriptionId, r.config.ResourceGroup, r.config.Name)
	_ = idString

	return ""
}

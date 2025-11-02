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

func (r ResourceInfo) ResourceGroup() string {
	return r.config.ResourceGroup
}

func (r ResourceInfo) ID() string {
	const resourceIdTemplate = "/subscriptions/%s/resourceGroups/%s/providers/%s/%s"
	return fmt.Sprintf(resourceIdTemplate, r.subscriptionId, r.config.ResourceGroup, r.config.Type, r.config.Name)
}

func (r ResourceInfo) Name() string {
	return r.config.Name
}

func (r ResourceInfo) Type() string {
	return r.config.Type
}

package model

func (project *Project) GetInstanceType(name string) InstanceType {
	for _, entry := range project.Instances {
		if entry.Name == name {
			return entry.InstanceType
		}
	}

	return ""
}

func appendResource(resources []KubernetesResource, newResource KubernetesResource) []KubernetesResource {

	for index := 0; index < len(resources); index++ {
		if resources[index].Name == newResource.Name && resources[index].Type == newResource.Type {
			resources[index].Properties = newResource.Properties
			resources[index].Updated = newResource.Updated
			resources[index].Status = newResource.Status
			return resources
		}
	}

	return append(resources, newResource)
}

func (instance *Instance) GetResource(resourceName string, resourceType ResourceObjectType) *KubernetesResource {

	var resources []KubernetesResource
	if resourceType == ResourceVolumeSnapshot {
		resources = instance.Resources.Snapshots
	} else if resourceType == ResourceSnapshotSchedule {
		resources = instance.Resources.Schedules
	} else {
		resources = instance.Resources.Resources
	}

	for _, resource := range resources {
		if resource.Name == resourceName && resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

func (instance *Instance) GetResourceByName(resourceName string) *KubernetesResource {

	for _, resource := range instance.Resources.Resources {
		if resource.Name == resourceName {
			return &resource
		}
	}

	for _, resource := range instance.Resources.Snapshots {
		if resource.Name == resourceName {
			return &resource
		}
	}

	for _, resource := range instance.Resources.Schedules {
		if resource.Name == resourceName {
			return &resource
		}
	}

	return nil
}

func (instance *Instance) GetResourceByType(resourceType ResourceObjectType) *KubernetesResource {

	var resources []KubernetesResource
	if resourceType == ResourceVolumeSnapshot {
		resources = instance.Resources.Snapshots
	} else if resourceType == ResourceSnapshotSchedule {
		resources = instance.Resources.Schedules
	} else {
		resources = instance.Resources.Resources
	}

	for _, resource := range resources {
		if resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

func (instance *Instance) AddResource(resource KubernetesResource) {

	if resource.Type == ResourceVolumeSnapshot {
		instance.Resources.Snapshots = append(instance.Resources.Snapshots, resource)
	} else if resource.Type == ResourceSnapshotSchedule {
		instance.Resources.Schedules = append(instance.Resources.Schedules, resource)
	} else {
		instance.Resources.Resources = append(instance.Resources.Resources, resource)
	}
}

func (instance *Instance) AddResources(resources ...KubernetesResource) {
	for _, resource := range resources {
		instance.AddResource(resource)
	}
}

func (instance *Instance) HasResources() bool {
	return len(instance.Resources.Resources) > 0 ||
		len(instance.Resources.Snapshots) > 0 ||
		len(instance.Resources.Schedules) > 0
}

func (instance *Instance) GetResources() []KubernetesResource {
	var resources = make([]KubernetesResource, 0)

	resources = append(resources, instance.Resources.Resources...)
	resources = append(resources, instance.Resources.Snapshots...)
	resources = append(resources, instance.Resources.Schedules...)

	return resources
}

func (project *Project) GetNamespace() string {
	return project.Name
}

func (project *Project) GetResource(resourceName string, resourceType ResourceObjectType) *KubernetesResource {
	for _, resource := range project.Resources {
		if resource.Name == resourceName && resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

func (project *Project) GetResourceByName(resourceName string) *KubernetesResource {
	for _, resource := range project.Resources {
		if resource.Name == resourceName {
			return &resource
		}
	}

	return nil
}

func (project *Project) GetResourceByType(resourceType ResourceObjectType) *KubernetesResource {
	for _, resource := range project.Resources {
		if resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

func (project *Project) AddResource(resource KubernetesResource) {
	if project.Resources == nil {
		project.Resources = make([]KubernetesResource, 0)
	}
	project.Resources = appendResource(project.Resources, resource)
}

func (project *Project) AddResources(resources ...KubernetesResource) {
	for _, resource := range resources {
		project.AddResource(resource)
	}
}

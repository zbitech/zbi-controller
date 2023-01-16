package data

import (
	"github.com/zbitech/controller/pkg/model"
)

var (
	Project1 = model.Project{Name: "project1", Network: model.NetworkTypeTest, Owner: Owner1_UserId, TeamId: Owner1Team_Id}
	Project2 = model.Project{Name: "project2", Network: model.NetworkTypeTest, Owner: Owner1_UserId, TeamId: Owner1Team_Id}
	Project3 = model.Project{Name: "project3", Network: model.NetworkTypeTest, Owner: Owner2_UserId, TeamId: Owner2Team_Id}
	Project4 = model.Project{Name: "project4", Network: model.NetworkTypeTest, Owner: Owner2_UserId, TeamId: Owner2Team_Id}

	Projects       = []model.Project{Project1, Project2, Project3, Project4}
	Owner1Projects = []model.Project{Project1, Project2}
	Owner2Projects = []model.Project{Project3, Project4}
)

func AppendProjects(projects []model.Project, _projects ...model.Project) []model.Project {
	return append(projects, _projects...)
}

func CreateProjects(count int, props map[string]interface{}) []model.Project {

	var projects = make([]model.Project, count)
	for index := range projects {
		projects[index] = *CreateProject(props)
	}

	return projects
}

func CreateProject(props map[string]interface{}) *model.Project {
	return &model.Project{
		Name:    getProperty(props, "name", randomString(10)).(string),
		Network: getProperty(props, "network", randomValue(networkTypes)).(model.NetworkType),
		Owner:   getProperty(props, "owner", randomString(10)).(string),
		TeamId:  getProperty(props, "team", randomString(50)).(string),
	}
}

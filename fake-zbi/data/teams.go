package data

var (
	Owner1Team_Id = "1"
	Owner2Team_Id = "2"

	Owner1TeamAdmin_Key = "11"
	Owner1TeamUser_Key  = "12"
	Owner2TeamAdmin_Key = "21"
	Owner2TeamUser_Key  = "22"

	//Owner1Team = entity.NewTeam("My Team", Owner1.Email, Owner1.UserId)
	//Owner2Team = entity.NewTeam("My Team", Owner2.Email, Owner2.UserId)
	//
	//Owner1TeamAdminMbr = entity.NewTeamMember(id.GenerateKey(), Owner1TeamAdmin.Email, true)
	//Owner1TeamUserMbr  = entity.NewTeamMember(id.GenerateKey(), Owner1TeamUser.Email, false)
	//Owner2TeamAdminMbr = entity.NewTeamMember(id.GenerateKey(), Owner1TeamAdmin.Email, true)
	//Owner2TeamUserMbr  = entity.NewTeamMember(id.GenerateKey(), Owner2TeamUser.Email, false)
)

//func AppendTeams(teams []entity.Team, _teams ...entity.Team) []entity.Team {
//	return append(teams, _teams...)
//}

//func CreateTeams(count int, props map[string]interface{}) []entity.Team {
//	var teams = make([]entity.Team, count)
//	for index := range teams {
//		teams[index] = CreateTeam(props)
//	}
//
//	return teams
//}

//func CreateTeam(props map[string]interface{}) entity.Team {
//
//	return entity.NewTeam(
//		getProperty(props, "name", randomString(10)).(string),
//		getProperty(props, "email", randomString(10)+"@zbitech.io").(string),
//		getProperty(props, "userid", randomString(10)).(string))
//}

//func AppendTeamMembers(members []entity.TeamMember, _members ...entity.TeamMember) []entity.TeamMember {
//	return append(members, _members...)
//}

//func CreateTeamMembers(count int, props map[string]interface{}) []entity.TeamMember {
//	var members = make([]entity.TeamMember, 10)
//	for index := range members {
//		members[index] = CreateTeamMember(props)
//	}
//	return members
//}

//func CreateTeamMember(props map[string]interface{}) entity.TeamMember {
//	return entity.NewTeamMember(id.GenerateKey(), getProperty(props, "email", randomString(10)+"@zbitech.io").(string), getProperty(props, "admin", randomValue(boolean)).(bool))
//}

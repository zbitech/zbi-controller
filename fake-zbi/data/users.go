package data

var (
	AdminUser_UserId = "admin"
	Owner1_UserId    = "owner1"
	Owner2_UserId    = "owner2"

	Owner1_TeamAdmin_UserId = "owner1admin"
	Owner1_TeamUser_UserId  = "owner1user"

	Owner2_TeamAdmin_UserId = "owner2admin"
	Owner2_TeamUser_UserId  = "owner2user"

	//Owner1TeamAdmin = entity.NewUser("owner1admin", "Owner1", "Admin", "owner1admin@zbitech.local", ztypes.RoleUser, ztypes.SubscriptionTeamMember)
	//Owner1TeamUser  = entity.NewUser("owner1user", "Owner1", "user", "owner1user@zbitech.local", ztypes.RoleUser, ztypes.SubscriptionTeamMember)
	//Owner2TeamAdmin = entity.NewUser("owner2admin", "Owner2", "Admin", "owner2admin@zbitech.local", ztypes.RoleUser, ztypes.SubscriptionTeamMember)
	//Owner2TeamUser  = entity.NewUser("owner2user", "Owner2", "User", "owner2user@zbitech.local", ztypes.RoleUser, ztypes.SubscriptionTeamMember)

	//AllUsers = []entity.User{*AdminUser, *Owner1, *Owner2, *Owner1TeamAdmin, *Owner1TeamUser, *Owner2TeamAdmin, *Owner2TeamUser}

	//Owner1Password          = "owner1password"
	//Owner2Password          = "owner2password"
	//Owner1TeamAdminPassword = "owner1adminpassword"
	//Owner1TeamUserPassword  = "owner1userpassword"
	//Owner2TeamAdminPassword = "owner2adminpassword"
	//Owner2TeamUserPassword  = "owner2userpassword"

	//AdminToken, _           = jwtutil.GenerateJwtToken(*AdminUser)
	//Owner1Token, _          = jwtutil.GenerateJwtToken(*Owner1)
	//Owner1TeamAdminToken, _ = jwtutil.GenerateJwtToken(*Owner1TeamAdmin)
	//Owner1TeamUserToken, _  = jwtutil.GenerateJwtToken(*Owner1TeamUser)
	//Owner2Token, _          = jwtutil.GenerateJwtToken(*Owner2)
	//Owner2TeamAdminToken, _ = jwtutil.GenerateJwtToken(*Owner2TeamAdmin)
	//Owner2TeamUserToken, _  = jwtutil.GenerateJwtToken(*Owner2TeamUser)
	//InvalidToken            = *Owner1Token + "FAKE"

	//AdminBasicCreds  = "Basic " + utils.Base64EncodeString(AdminUser.UserId+":"+vars.ADMIN_PASSWORD)
	//Owner1BasicCreds = "Basic " + utils.Base64EncodeString(Owner1.UserId+":"+Owner1Password)
	//Owner2BasicCreds = "Basic " + utils.Base64EncodeString(Owner2.UserId+":"+Owner2Password)
)

//func AppendUsers(users []entity.User, _users ...entity.User) []entity.User {
//	return append(users, _users...)
//}

//func CreateBasicCredentials(users []entity.User, passwords []string) []string {
//	var credentials = make([]string, len(users))
//	for index := range credentials {
//		credentials[index] = "Basic " + utils.Base64EncodeString(users[index].UserId+":"+passwords[index])
//	}
//	return credentials
//}

//func CreatePasswords(count int) []string {
//	var passwords = make([]string, count)
//	for index := range passwords {
//		passwords[index] = id.GenerateSecurePassword()
//	}
//
//	return passwords
//}

//func CreateUsers(count int, props map[string]interface{}) []entity.User {
//	var users = make([]entity.User, count)
//	for index := range users {
//		users[index] = *CreateUser(props)
//	}
//	return users
//}

//func CreateUser(props map[string]interface{}) *entity.User {
//	return entity.NewUser(getProperty(props, "userid", randomString(10)).(string),
//		getProperty(props, "firstName", randomString(15)).(string),
//		getProperty(props, "lastName", randomString(15)).(string),
//		getProperty(props, "email", randomString(15)+"@zbitech.local").(string),
//		getProperty(props, "role", randomValue(roleTypes)).(ztypes.Role),
//		getProperty(props, "subscription", randomValue(subscriptionTypes)).(ztypes.SubscriptionLevel))
//}

//func CreateAPKeys(count int, userid string) []entity.APIKey {
//	var keys = make([]entity.APIKey, count)
//	for index := range keys {
//		keys[index] = CreateAPIKey(userid)
//	}
//	return keys
//}

//func CreateAPIKey(userid string) entity.APIKey {
//	return entity.NewAPIKey(userid, vars.AppConfig.Policy.TokenExpirationPolicy)
//}

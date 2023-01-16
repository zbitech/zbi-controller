package request

//type RegistrationInvite struct {
//	Email string                  `json:"email" validate:"required"`
//	Role  model.Role              `json:"role" validate:"required"`
//	Level model.SubscriptionLevel `json:"level" validate:"required"`
//}

//type Registration struct {
//	UserId          string `json:"userid" validate:"required,alphanum"`
//	FirstName       string `json:"firstName" validate:"required,alpha"`
//	LastName        string `json:"lastName" validate:"required,alpha"`
//	Email           string `json:"email" validate:"required,email"`
//	Password        string `json:"password" validate:"required,min=10"`
//	ConfirmPassword string `json:"confirmPassword" validate:"required"`
//	InviteKey       string `json:"inviteKey" validate:"required"`
//}

//type NewUser struct {
//	UserId    string                   `json:"userid" validate:"required,alphanum"`
//	FirstName string                   `json:"firstName" validate:"required,alpha"`
//	LastName  string                   `json:"lastName" validate:"required,alpha"`
//	Email     string                   `json:"email" validate:"required,email"`
//	Role      model.Role              `json:"role" validate:"required"`
//	Level     model.SubscriptionLevel `json:"level" validate:"required"`
//	Team      string                   `json:"team"`
//}

//type RegistrationInvites struct {
//	Invites []struct {
//		Key   string `json:"key" validate:"required"`
//		Email string `json:"email" validate:"required,email"`
//	} `json:"invites" validate:"required"`
//}

//type TeamInvitation struct {
//	Action  string `json:"action" validate:"required"`
//	Members []struct {
//		Key   string `json:"key"`
//		Email string `json:"email" validate:"email"`
//		Admin bool   `json:"admin"`
//	} `json:"members"`
//}

//type TeamSettings struct {
//	Name  string `json:"name" validate:"required,alphanum"`
//	Email string `json:"email" validate:"required,email"`
//}

//type TeamMembership struct {
//	Key    string                  `json:"key" validate:"required"`
//	Status ztypes.InvitationStatus `json:"status" validate:"required"`
//}

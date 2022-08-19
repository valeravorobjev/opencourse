package common

// OpenUser user common model
type OpenUser struct {
	Id     string   `json:"id"`     // User id
	Login  string   `json:"login"`  // User login
	Name   string   `json:"name"`   // User display name
	Email  string   `json:"email"`  // Email user address
	Avatar string   `json:"avatar"` // User avatar image path
	Rating int      `json:"rating"` // User rating
	Roles  []string `json:"roles"`  // User roles
}

// OpenCreateUserQuery model for create user
type OpenCreateUserQuery struct {
	Login    string   `json:"login"`    // User login
	Password string   `json:"password"` // User password
	Email    string   `json:"email"`    // Email user address
	Name     string   `json:"name"`     // User display name
	Avatar   string   `json:"avatar"`   // User avatar image path
	Roles    []string `json:"roles"`    // User roles
}

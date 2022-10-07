package openerrors

// NOTE! In this project, as an experiment, all openerrors are wrapped in special types
// Open - project prefix

// OpenBaseErr base error model
type OpenBaseErr struct {
	File   string // File contains error
	Method string // Function or method that throw error
}

// OpenDefaultErr default error model
type OpenDefaultErr struct {
	BaseErr OpenBaseErr // File contains error
	Msg     string      // Error text message
}

// OpenFieldEmptyErr Error if required field is empty
type OpenFieldEmptyErr struct {
	BaseErr OpenBaseErr // File contains error
	Field   string      // Field name
}

// OpenMinLenErr Error not matching the minimum length
type OpenMinLenErr struct {
	BaseErr OpenBaseErr // File contains error
	Field   string      // Field name
	MinLen  int         // Minimum password length
}

// OpenRoleUnknownErr Error throw if role is not contained in list
type OpenRoleUnknownErr struct {
	BaseErr OpenBaseErr // File contains error
	Role    string      // Enter role
	Roles   []string    // List of available roles
}

// OpenDbErr database errors
type OpenDbErr struct {
	BaseErr OpenBaseErr // File contains error
	DbName  string      // Database's name. Example: mongodb/opencourse
	ConStr  string      // Connection string
	DbErr   string      // Database's error
}

// OpenInvalidIdErr convert id error from user string to db object
type OpenInvalidIdErr struct {
	Default   OpenDefaultErr // File contains error
	Id        string         // User id string
	Converter string         // Converter function
}

package openerrors

// NOTE! In this project, as an experiment, all openerrors are wrapped in special types
// Open - project prefix

// BaseErr base error model
type BaseErr struct {
	File   string // File contains error
	Method string // Function or method that throw error
}

// DefaultErr default error model
type DefaultErr struct {
	BaseErr BaseErr // File contains error
	Msg     string  // Error text message
}

// FieldEmptyErr error if required field is empty in the model
type FieldEmptyErr struct {
	BaseErr BaseErr // File contains error
	Field   string  // Field name
}

// ModelNilOrEmptyErr error if model parameter nil or empty
type ModelNilOrEmptyErr struct {
	BaseErr BaseErr // File contains error
	Model   string  // Model name
}

// MinLenErr Error not matching the minimum length
type MinLenErr struct {
	BaseErr BaseErr // File contains error
	Field   string  // Field name
	MinLen  int     // Minimum password length
}

// RoleUnknownErr Error throw if role is not contained in list
type RoleUnknownErr struct {
	BaseErr BaseErr  // File contains error
	Role    string   // Enter role
	Roles   []string // List of available roles
}

// DbErr database errors
type DbErr struct {
	BaseErr BaseErr // File contains error
	DbName  string  // Database's name. Example: mongodb/opencourse
	ConStr  string  // Connection string
	DbErr   string  // Database's error
}

// InvalidIdErr convert id error from user string to db object
type InvalidIdErr struct {
	Default   DefaultErr // File contains error
	Id        string     // User id string
	Converter string     // Converter function
}

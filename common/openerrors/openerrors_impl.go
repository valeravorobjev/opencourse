package openerrors

import (
	"fmt"
	"strings"
)

// DefaultErr common error implementation
func (err BaseErr) Error() string {
	return fmt.Sprintf("file: %s | method: %s", err.File, err.Method)
}

// DefaultErr common error implementation
func (err DefaultErr) Error() string {
	return fmt.Sprintf("%s | message: %s", err.BaseErr.Error(), err.Msg)
}

// FieldEmptyErr implementation
func (err FieldEmptyErr) Error() string {
	return fmt.Sprintf("%s | message: field %s is empty", err.BaseErr.Error(), err.Field)
}

// ModelNilOrEmptyErr implementation
func (err ModelNilOrEmptyErr) Error() string {
	return fmt.Sprintf("%s | message: field %s is empty", err.BaseErr.Error(), err.Model)
}

// MinLenErr implementation
func (err MinLenErr) Error() string {
	return fmt.Sprintf("%s | message: field %s length is less than %d",
		err.BaseErr.Error(), err.Field, err.MinLen)
}

// RoleUnknownErr implementation
func (err RoleUnknownErr) Error() string {
	return fmt.Sprintf("%s | message: the role %s is not contained in the list of available roles (%s) ",
		err.BaseErr.Error(), err.Role, strings.Join(err.Roles, ","))
}

// DbErr implementation
func (err DbErr) Error() string {
	return fmt.Sprintf("%s | db name: %s | instance: %s | message: %s",
		err.BaseErr.Error(), err.DbName, err.ConStr, err.DbErr)
}

// InvalidIdErr implementation
func (err InvalidIdErr) Error() string {
	return fmt.Sprintf("%s | user id value: %s | converter: %s",
		err.Default.Error(), err.Id, err.Converter)
}

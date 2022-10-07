package openerrors

import (
	"fmt"
	"strings"
)

// OpenDefaultErr common error implementation
func (err OpenBaseErr) Error() string {
	return fmt.Sprintf("file: %s | method: %s", err.File, err.Method)
}

// OpenDefaultErr common error implementation
func (err OpenDefaultErr) Error() string {
	return fmt.Sprintf("%s | message: %s", err.BaseErr.Error(), err.Msg)
}

// OpenFieldEmptyErr implementation
func (err OpenFieldEmptyErr) Error() string {
	return fmt.Sprintf("%s | message: field %s is empty", err.BaseErr.Error(), err.Field)
}

// OpenMinLenErr implementation
func (err OpenMinLenErr) Error() string {
	return fmt.Sprintf("%s | message: field %s length is less than %d",
		err.BaseErr.Error(), err.Field, err.MinLen)
}

// OpenRoleUnknownErr implementation
func (err OpenRoleUnknownErr) Error() string {
	return fmt.Sprintf("%s | message: the role %s is not contained in the list of available roles (%s) ",
		err.BaseErr.Error(), err.Role, strings.Join(err.Roles, ","))
}

// OpenDbErr implementation
func (err OpenDbErr) Error() string {
	return fmt.Sprintf("%s | db name: %s | instance: %s | message: %s",
		err.BaseErr.Error(), err.DbName, err.ConStr, err.DbErr)
}

// OpenInvalidIdErr implementation
func (err OpenInvalidIdErr) Error() string {
	return fmt.Sprintf("%s | user id value: %s | converter: %s",
		err.Default.Error(), err.Id, err.Converter)
}

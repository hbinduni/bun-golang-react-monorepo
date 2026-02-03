package utils

import (
	"go.jetify.com/typeid"
)

// TypeID prefixes for different entity types
const (
	PrefixUser         = "user"
	PrefixItem         = "item"
	PrefixSession      = "sess"
	PrefixOAuthAccount = "oauth"
)

// NewUserID generates a new TypeID for a user
func NewUserID() string {
	tid, _ := typeid.WithPrefix(PrefixUser)
	return tid.String()
}

// NewItemID generates a new TypeID for an item
func NewItemID() string {
	tid, _ := typeid.WithPrefix(PrefixItem)
	return tid.String()
}

// NewSessionID generates a new TypeID for a session
func NewSessionID() string {
	tid, _ := typeid.WithPrefix(PrefixSession)
	return tid.String()
}

// NewOAuthAccountID generates a new TypeID for an OAuth account
func NewOAuthAccountID() string {
	tid, _ := typeid.WithPrefix(PrefixOAuthAccount)
	return tid.String()
}

// ValidateTypeID validates a TypeID string format
func ValidateTypeID(s string) bool {
	// Basic validation - check format prefix_base32
	if len(s) < 3 {
		return false
	}
	// Just check if it has an underscore separator
	for i := 0; i < len(s); i++ {
		if s[i] == '_' {
			return i > 0 && i < len(s)-1
		}
	}
	return false
}

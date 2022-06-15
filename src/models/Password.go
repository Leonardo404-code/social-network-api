package models

// The password represents the password change request format.
type Password struct {
	New string `json:"new"`
	Old string `json:"old"`
}

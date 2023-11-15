package controllers

const (
	invalidState         = "invalid_state"
	invalidNonce         = "invalid_nonce"
	invalidCodeChallenge = "invalid_code_challenge"
	successMessage       = "success"

	userDoesNotHaveThisRole = "user_does_not_have_this_role"
)

var statusCode = map[string]int{
	successMessage: 0,

	// Domain error
	userDoesNotHaveThisRole: 1001,

	// Application error
	invalidState:         2001,
	invalidNonce:         2002,
	invalidCodeChallenge: 2003,
}

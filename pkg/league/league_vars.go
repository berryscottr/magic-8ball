package league

// Data for the bot to track along a request
type Data struct {
	// Err for error tracking
	Err error
	// User data
	User Usr
}

// User for website
type Usr struct {
	// Email
	Email string
	// Password
	Password string
}

// Methods for interacting with the league website
type Methods interface {
	// Login to league website
	Login()
}
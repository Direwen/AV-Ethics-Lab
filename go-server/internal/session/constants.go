package session

type SessionStatus int

const (
	StatusActive    SessionStatus = 1
	StatusCompleted SessionStatus = 2
	StatusExpired   SessionStatus = 3
	StatusAbandoned SessionStatus = 4
)

var SessionStatusErrorMsg = map[SessionStatus]string{
	StatusCompleted: "session is completed",
	StatusExpired:   "session is expired",
	StatusAbandoned: "session is abandoned",
}

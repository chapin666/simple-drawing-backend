package message

const (
	// KindConnected is sent when user connects
	KindConnected = iota + 1
	// KindUserJoined is sent when someone else joins
	KindUserJoined
	// KindUserLeft is sent when someone leaves
	KindUserLeft
	// KindStroke message specifies a drawn stroke by a user
	KindStroke
	// KindClear message is send when a user clears the screen
	KindClear
)

// Point struct.
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// User struct.
type User struct {
	ID    string `json:"id"`
	Color string `json:"color"`
}

// Connected struct.
type Connected struct {
	Kind  int    `json:"kind"`
	Color string `json:"color"`
	Users []User `json:"users"`
}

// NewConnected function
func NewConnected(color string, users []User) *Connected {
	return &Connected{
		Kind:  KindConnected,
		Color: color,
		Users: users,
	}
}

// UserJoined  struct.
type UserJoined struct {
	Kind int  `json:"kind"`
	User User `json:"user"`
}

// NewUserJoined function
func NewUserJoined(userID string, color string) *UserJoined {
	return &UserJoined{
		Kind: KindUserJoined,
		User: User{ID: userID, Color: color},
	}
}

// UserLeft struct.
type UserLeft struct {
	Kind   int    `json:"kind"`
	UserID string `json:"userId"`
}

// NewUserLeft function
func NewUserLeft(userID string) *UserLeft {
	return &UserLeft{
		Kind:   KindUserLeft,
		UserID: userID,
	}
}

// Stroke struct
type Stroke struct {
	Kind   int     `json:"kind"`
	UserID string  `json:"userId"`
	Points []Point `json:"points"`
	Finish bool    `json:"finish"`
}

// Clear struct
type Clear struct {
	Kind   int    `json:"kind"`
	UserID string `json:"userId"`
}

package entity

type User struct {
	Username string `json:"login"`
	Password string `json:"password,omitempty"`
}

func InitNewUser(username, password string) (User, error) {
	u := User{
		Username: username,
		Password: password,
	}

	return u, nil
}

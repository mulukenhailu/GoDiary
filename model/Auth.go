package model

type Auth struct { //basic struct for login and Registration purpose
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

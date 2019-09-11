package auth

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}


func (creds *Credentials) HashPassword() (hash []byte, err error) {

}
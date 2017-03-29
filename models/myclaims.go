package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//MyClaims custom Claims
type MyClaims struct {
	Name  string `json:"name"`  //user name
	DevID string `json:"devid"` //device uuid
	jwt.StandardClaims
}

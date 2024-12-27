package main

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var tokenStr = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJzdWIiOiIxMDEiLCJleHAiOjE3MzUyMDg0MTEsImlhdCI6MTczNTIwNTQxMSwicm9sZXMiOlsiVVNFUiJdfQ.MtufBAyRYgPJV_Lov3pShBIC8eT8ovuHK4v2CcBfMQkRAWe0KzC_sbUN9DzBg4Pts8z7p0cgfFGFpscGTL1pyuG2LkvJ7WtFrXHUwVmFjsAwFroTm3K8KPQbMRxfVuL74N_bJI5UFBZ3VaiSOWBRLS0QH5Eq0nFUg29QWekR-jJTVBDY4gr-nFHuP03sECnmrTiPIp3HaOHUghQkPQqhLNbvKzHbvd2Poqi6UGCnAGWWs1U9o11AuKNd73t2vp-XKKMpf83mUkm6-zauMyIENXEciGWcLZcsEy3XbWc9Cel221jcWM4tDBT8BYfzl8bDrsXXclJiHlUAYUegWqbaIA`

func main() {
	pubKeyPem, err := os.ReadFile(".//..//public.pem")
	if err != nil {
		panic(err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyPem)
	if err != nil {
		panic(err)
	}
	var claims struct {
		jwt.RegisteredClaims
		Roles []string `json:"roles"`
	}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})

	if err != nil {
		panic(err)
	}
	if !token.Valid {
		panic("token is not valid")
	}

	fmt.Println(claims.Subject, ", this user requested this role:", claims.Roles)
}

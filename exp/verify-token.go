package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
)

var token = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkgcHJvamVjdCIsInN1YiI6IjEwMSIsImV4cCI6MTc0MjM2Mjc1NCwiaWF0IjoxNzQyMzU5NzU0fQ.fdBEgySN7ZD09QcUOTJhIOE1t57WS7BnEL7Hrf3fWbC3xWmtf2OaPdjwvrD1Cn0Kjm_z8F6VSSwkpeW2KLVI8UVvwkxak99SgEsTPL6QXXQ9oA0yN0bs15TW7gbm2YTHYylINLEBo8YfMkzgrH-BqHHpMl9yJBTLX6WvI0YXu_82D5xj9CjkobteCeEZWIm_4HrbmRoTksCe0RigwUllFohZ2qWZ1NZZguWonVFE1YRh6BdkGKUL4Hz01ETL22XQ1VkocUFr131WeNAZZIC31FBQJuG0G0v0k6dC57aVevrK8jl-pPRfjZJrMkyUdt3GBbi7uu6tu3eMikhiTp--Hw
`

func main() {
	PublicPEM, err := os.ReadFile("pubkey.pem")
	if err != nil {
		// If there's an error reading the file, print an error message and stop execution
		log.Fatalln("not able to read pem file")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(PublicPEM)
	if err != nil {
		// If there's an error parsing the public key, log the error and stop execution
		log.Fatalln(err)
	}

	var claims jwt.RegisteredClaims
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}
	tkn, err := jwt.ParseWithClaims(token, &claims, keyFunc)
	if err != nil {
		log.Fatalln(err)
	}
	if !tkn.Valid {
		log.Fatalln("token is not valid")
	}
	fmt.Println(claims)
}

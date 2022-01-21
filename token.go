package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

func GetNewToken(userid uint64) (*TokenDetails, error) {
	var err error
	td.AccessToken, err = CreateAccessToken(userid)
	fmt.Printf("access token is : %v", td.AccessToken)
	if err != nil {
		return nil, err
	}
	td.RefreshToken, err = CreateRefreshToken(userid)
	if err != nil {
		return nil, err
	}
	return td, err
}

func CreateAccessToken(userid uint64) (string, error) {
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()
	var err error
	//Creating Access Token
	err = os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	if err != nil {
		log.Fatal(err)
	}
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("access token is : %v", td.AccessToken)
	return td.AccessToken, nil
}

/*
Since the uuid is unique each time it is created, a user can create more than one token. This happens when a user is logged in on different devices.
The user can also log out from any of the devices without them being logged out from all devices. How cool!
*/

// CheckValidToken : Check the validity of this token, whether it is still useful or it has expired
func CheckValidToken(r *http.Request) error {
	token, err := VerifyTokenFromHeaders(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

/*
For a production grade application, it is highly recommended to store JWTs in an HttpOnly cookie.
To achieve this, while sending the cookie generated from the backend to the frontend (client),
a HttpOnly flag is sent along with the cookie, instructing the browser not to display the cookie through the client-side scripts.
Doing this can prevent XSS (Cross Site Scripting) attacks.
JWT can also be stored in browser local storage or session storage.
Storing a JWT this way can expose it to several attacks such as XSS mentioned above, so it is generally less secure when compared to using `HttpOnly cookie technique.
*/

// GetTokenFromHeaders : Extract the token from Headers
func GetTokenFromHeaders(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// VerifyTokenFromHeaders : Verify the headers token with SigningMethodHMAC
func VerifyTokenFromHeaders(r *http.Request) (*jwt.Token, error) {
	tokenString := GetTokenFromHeaders(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// GetTokenMetadata : Get JWT Token Metadata
func GetTokenMetadata(r *http.Request) (*TokenMetadata, error) {
	token, err := VerifyTokenFromHeaders(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &TokenMetadata{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

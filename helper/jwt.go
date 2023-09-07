package helper

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mulukenhailu/Diary_api/model"
)

var PrivateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(user model.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // preparing the content of the token
		"id":  user.ID, // orignally there is no ID field on the user Stuct , so how the user.ID is functioning?
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})

	return token.SignedString(PrivateKey) // signing preared token with the privatekey
}

// from the requestBody extract the jwt part which come after loginin  a user
func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return "" // just return the token whatever it is ...
}

// from currently retrived jwt extract the privatekey part ansd then return

func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)                                        // the token from request
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // if the signing method is fine, get the privateKey field used when signing, but it do not check for validity
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //check if the signing method is ok! or if it is the same as defiend in making of the jwt
			return nil, fmt.Errorf("unexperted signing method: %v", token.Header["alg"])
		}
		return PrivateKey, nil
	})
	return token, err
}

// cross check if the extracted privateKey and true privateKey are of the same kind
func ValidateJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims) // (claim, ok)....check the claim  made, ie check if the token format is in the format specied
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid Token Provided")
}

func CurrentUser(context *gin.Context) (model.User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return model.User{}, err
	}

	token, _ := getToken(context)             // get the token
	claims, _ := token.Claims.(jwt.MapClaims) //get the claim from the token, ...{ id : ""  , iat:""   ,eat:""}
	userId := uint(claims["id"].(float64))

	user, err := model.FindUserById(userId)
	if err != nil {
		return model.User{}, err
	}

	return user, nil

}

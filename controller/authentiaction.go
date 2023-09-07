package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mulukenhailu/Diary_api/helper"
	"github.com/mulukenhailu/Diary_api/model"
)

func Register(context *gin.Context) {
	var input model.Auth
	if err := context.ShouldBindJSON(&input); err != nil { //from the request body extract the json in the  format of model.Auth
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{ // creating user but the content field is not provided yet.
		Username: input.Username,
		Password: input.Password,
	}

	savedUser, err := user.Save() // save the user to the database

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

func Login(context *gin.Context) {
	// 1.extract the json format from the request body and assign it to struct in the form of Auth
	var input model.Auth
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2.retrive the user data using the request body, input

	user, err := model.FindUserByUsername(input.Username)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 3.cross check the password from the retrived user struct and password from request body
	err = user.ValidatePassword(input.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if ok!, generate the jwt for the current loged user for permission of accessing the private apis
	jwt, err := helper.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	context.JSON(http.StatusCreated, gin.H{"jwt": jwt})
}

func AddEntry(context *gin.Context) {

	//1...extracting the json in the format of model Entry from the request body.. it has the content field
	var input model.Entry
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2.get the user but first check if the token claim was valid

	user, err := helper.CurrentUser(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserId = user.ID // fill out the UserId field of the request body

	savedEntry, err := input.Save() //3. call the save method on the entry struct created
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"entry": savedEntry})

}

func GetAllEntry(context *gin.Context) {
	user, err := helper.CurrentUser(context) // Authenticate the user using the jwt
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Entries": user.Entries})
}

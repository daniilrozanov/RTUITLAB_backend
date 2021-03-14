package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	templates "purchases/pkg"
	"strconv"
)

type SignInData struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

type ConfirmData struct {
	Username []byte `json:"name"`
	Password []byte `json:"password"`
}

var (
	usersTransportKey = os.Getenv("USERS_TRANSPORT_KEY")
)

// @Summary SignUp
// @Tags auth
// @Description register in purchases service
// @ID sign-up
// @Accept json
// @Produce json
// @Param input body SignInData true "Sign In Data"
// @Success 200 {string} string "token"
// @Failure default {object} Error
// @Router /signup [post]
func (h *Handler) SignUp (g *gin.Context){
	var input templates.User
	var sid SignInData


	if err := g.BindJSON(&sid); err != nil{ //присвоит значения из json полям с совпадающими тегами
		newErrorResponse(g, http.StatusBadRequest, err.Error())
		return
	}
	input.Name = sid.Username
	input.Password = sid.Password
	id, err := h.buis.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	g.JSON(http.StatusOK, map[string]interface{} {
		"id": id,
	})
}

// @Summary SignIn
// @Tags auth
// @Description get auth token
// @ID sign-in
// @Accept json
// @Produce json
// @Param input body SignInData true "Sign In Data"
// @Success 200 {string} string "token"
// @Failure default {object} Error
// @Router /signin [post]
func (h *Handler) SingIn (g *gin.Context){
	var input SignInData

	if err := g.BindJSON(&input); err != nil{
		newErrorResponse(g, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.buis.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	g.JSON(http.StatusOK, map[string]interface{} {
		"token": token,
	})
}

// @Summary ConfirmUser
// @Tags export auth
// @Description send encrypted user id if exists
// @ID confirm
// @Accept json
// @Produce json
// @Param input body ConfirmData true "Confirm Data"
// @Success 200 {object} map[string]string "token"
// @Failure default {object} Error
// @Router /confirm [post]
func (h *Handler) ConfirmUser (g *gin.Context){
	var input ConfirmData

	if err := g.BindJSON(&input); err != nil{
		newErrorResponse(g, http.StatusBadRequest, err.Error())
		return
	}
	userName, err := templates.Decrypt(input.Username, usersTransportKey)
	if err != nil {
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	userPass, err := templates.Decrypt(input.Password, usersTransportKey)
	if err != nil {
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	id, err := h.buis.Authorization.GetUserId(string(userName), string(userPass))
	if err != nil {
		strNotFound, err := templates.Encrypt([]byte("not found"), usersTransportKey)
		if err != nil {
			newErrorResponse(g, http.StatusInternalServerError, err.Error())
			return
		}
		g.JSON(http.StatusOK, map[string][]byte {
			"id": strNotFound,
		})
		return
	}
	strId, err := templates.Encrypt([]byte(strconv.Itoa(id)), usersTransportKey)
	if err != nil {
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	g.JSON(http.StatusOK, map[string][]byte {
		"id": strId,
	})
}
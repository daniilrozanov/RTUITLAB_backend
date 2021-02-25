package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	templates "purchases/pkg"
	"strconv"
)

type SingInData struct {
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

func (h *Handler) SignUp (g *gin.Context){
	var input templates.User

	if err := g.BindJSON(&input); err != nil{ //присвоит значения из json полям с совпадающими тегами
		newErrorResponse(g, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.buis.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	g.JSON(http.StatusOK, map[string]interface{} {
		"id": id,
	})
}

func (h *Handler) SingIn (g *gin.Context){
	var input SingInData

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
package controllers

import (
	"encoding/json"
	"github.com/GoProjectGroupForEducation/Go-Blog/models"
	"github.com/GoProjectGroupForEducation/Go-Blog/utils"
	"net/http"
)

func GetAllTag(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	res, err := json.Marshal(models.GetAllTags())
	if err != nil {
		panic(err)
	}
	return utils.SendData(w, string(res), "OK", http.StatusOK)
}
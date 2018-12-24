package controllers

import (
	"encoding/json"
	"github.com/GoProjectGroupForEducation/Go-Blog/models"
	"github.com/GoProjectGroupForEducation/Go-Blog/services"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GoProjectGroupForEducation/Go-Blog/utils"
)

func Auth(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	buff, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	body := make(map[string]interface{})

	err = json.Unmarshal(buff, &body)
	if err != nil {
		return err
	}

	if _, ok := body["username"]; ok {
		//password为空
		if _, ok := body["password"]; !ok {
			return utils.SendData(w, string(buff), "Please input password", http.StatusBadRequest)
		}
		//username不为空
		password := string(body["password"].(string))
		username := string(body["username"].(string))
		//从数据库中搜索username，判断password是否匹配
		user := models.GetUserByUsername(username)

		if user != nil {
			if user.Password == password {
				token := services.GenerateAuthToken(user)

				newuser := models.GetUserListByID(user.UserID)
				data, err := json.Marshal(*newuser)
				if err != nil {
					return err
				}
				return utils.SendData(w, `{` +
					`"user":` + string(data) + `,` +
					`"token":"` + string(token.Token) +
					`"}`, "OK", http.StatusOK)
			} else {
				return utils.SendData(w, string(buff), "Wrong password", http.StatusBadRequest)
			}
		} else {
			return utils.SendData(w, string(buff), "User not found", http.StatusBadRequest)
		}

	} else {
		log.Println(req.Method, req.URL.String(), http.StatusBadRequest)
		return utils.SendData(w, string(buff), "Please input username", http.StatusBadRequest)
	}
}

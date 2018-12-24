package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/GoProjectGroupForEducation/Go-Blog/models"
	"github.com/GoProjectGroupForEducation/Go-Blog/services"
	"github.com/GoProjectGroupForEducation/Go-Blog/utils"
)

func GetAllArticles(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	req.ParseForm()
	if len(req.Form["pageNum"]) > 0{
		num, err := strconv.Atoi(req.Form["pageNum"][0])
		if err == nil {
			res, _ := json.Marshal(models.GetArticlesPerPage(num))
			if err != nil {
				return err
			}
			return utils.SendData(w, string(res), "OK", http.StatusOK)
		}
	}
	res, err := json.Marshal(models.GetAllArticles())
	if err != nil {
		return err
	}
	return utils.SendData(w, string(res), "OK", http.StatusOK)
}

func GetArticlesByTag(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	vars := mux.Vars(req)
	tag := vars["tag"]
	res, err := json.Marshal(models.GetArticlesByTag(tag))
	if err != nil {
		return err
	}
	return utils.SendData(w, string(res), "OK", http.StatusOK)
}

func GetArticlesByUserID(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	vars := mux.Vars(req)
	userid, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	res, err := json.Marshal(models.GetArticleByUserID(userid))
	if err != nil {
		return err
	}
	return utils.SendData(w, string(res), "OK", http.StatusOK)
}

func GetConcernArticles(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	author := services.GetCurrentUser(req.Header.Get("Authorization"))
	articleList := []models.ArticleList{}
	for _, one := range author.Following {
		userArticles := models.GetArticleByUserID(one)
		articleList = append(articleList, userArticles...)
	}
	res, err := json.Marshal(articleList)
	if err != nil {
		return err
	}
	return utils.SendData(w, string(res), "OK", http.StatusOK)
}

func CreateArticle(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	var article = models.Article{}
	author := services.GetCurrentUser(req.Header.Get("Authorization"))
	buff, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buff, &article)
	if err != nil {
		panic(err)
	}
	article.Author = author.UserID

	id := models.CreateArticle(article)

	return utils.SendData(w, `{
  "id": `+strconv.Itoa(id)+`
}`, "Create successfully", http.StatusOK)
}

func GetArticleByID(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	article := models.GetArticleByID(id)
	if article == nil {
		return utils.SendData(w, "{}", "Not Found", http.StatusNotFound)
	}
	data, err := json.Marshal(article)
	if err != nil {
		panic(err)
	}
	return utils.SendData(w, string(data), "OK", http.StatusOK)
}

func UpdateArticleByID(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error {
	author := services.GetCurrentUser(req.Header.Get("Authorization"))
	var article = models.Article{}
	buff, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	article.ID = id
	article.Author = author.UserID
	article.UpdatedAt = time.Now()
	err = json.Unmarshal(buff, &article)
	if err != nil {
		return err
	}
	models.UpdateArticleByID(article)
	return utils.SendData(w, "{}", "Update successfully", http.StatusOK)
}
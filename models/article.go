package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/GoProjectGroupForEducation/Go-Blog/utils"
)

type ArticleList struct {
	ID 				  int    		`json:"id"`
	Author    		  UserList    	`json:"author"`
	AuthorId		  int 			`json:"author_id"`
	Content   		  string    	`json:"content"`
	Comments		  []CommentList	`json:"comments"`
	Title     		  string 		`json:"title"`
	Tags	  		  []string  	`json:"tags"`
	UpdatedAt 		  string 		`json:"updated_at"`
}

type Article struct {
	ID 		  int       `json:"id, omitempty"`
	Author    int       `json:"author_id"`
	Title     string    `json:"title"`
	Tags	  []string  `json:"tags"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllArticles() []ArticleList {
	db := &utils.DB{}
	var articles []ArticleList
	var article ArticleList
	var articlesBytes map[string]string
	articlesBytes = db.Scan("article")
	if len(articlesBytes) == 0 {
		return []ArticleList{}
	}
	for _, one := range articlesBytes {
		err := json.Unmarshal([]byte(one), &article)
		article.Author = *GetUserListByID(article.AuthorId)
		article.Comments = GetAllCommentsByArticleID(article.ID)
		if article.Comments == nil {
			fmt.Println(1)
		}
		articles = append(articles, article)
		if err != nil {
			panic(err)
		}
	}
	return articles
}

func CreateArticle(article Article) int {
	db := &utils.DB{}
	id := db.GenerateID("article")
	article.ID = id
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	buff, err := json.Marshal(article)
	if err != nil {
		panic("JSON parsing error")
	}
	db.Set("article", strconv.Itoa(id), string(buff))
	return id
}

func GetArticleByID(id int) *ArticleList {
	db := &utils.DB{}
	buff := db.Get("article", strconv.Itoa(id))
	if len(buff) == 0 {
		return nil
	}
	article := ArticleList{}
	err := json.Unmarshal(buff, &article)
	article.Author = *GetUserListByID(article.AuthorId)
	article.Comments = GetAllCommentsByArticleID(article.ID)
	if err != nil {
		panic(err)
	}
	return &article
}

func GetArticleByUserID(id int) []ArticleList {
	db := &utils.DB{}
	var articles []ArticleList
	var article ArticleList
	var articlesBytes map[string]string
	articlesBytes = db.Scan("article")
	if len(articlesBytes) == 0 {
		return []ArticleList{}
	}
	for _, one := range articlesBytes {
		err := json.Unmarshal([]byte(one), &article)
		if article.Author.UserID == id {
			article.Author = *GetUserListByID(article.AuthorId)
			article.Comments = GetAllCommentsByArticleID(article.ID)
			articles = append(articles, article)
		}
		if err != nil {
			panic(err)
		}
	}
	return articles
}

func UpdateArticleByID(article Article) bool {
	db := &utils.DB{}
	buff := db.Get("article", strconv.Itoa(article.ID))
	if len(buff) == 0 {
		return false
	}
	buff, err := json.Marshal(article)
	if err != nil {
		panic("JSON parsing error")
	}
	db.Set("article", strconv.Itoa(article.ID), string(buff))
	return true
}

func GetArticlesByTag(tag string) []ArticleList {
	db := &utils.DB{}
	var articles []ArticleList
	var article ArticleList
	var articlesBytes map[string]string
	articlesBytes = db.Scan("article")
	if len(articlesBytes) == 0 {
		return []ArticleList{}
	}
	for _, one := range articlesBytes {
		err := json.Unmarshal([]byte(one), &article)
		for _,v := range article.Tags {
			if v == tag {
				article.Author = *GetUserListByID(article.AuthorId)
				article.Comments = GetAllCommentsByArticleID(article.ID)
				articles = append(articles, article)
			}
		}
		if err != nil {
			panic(err)
		}
	}
	return articles
}

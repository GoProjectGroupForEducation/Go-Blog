package models

import (
	"fmt"
	"github.com/GoProjectGroupForEducation/Go-Blog/utils"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type ArticlePage struct {
	Number   int `json: "number"`
	Articles []ArticleList `json: "articles"`
}


type ArticleList struct {
	ID 				  int    		`json:"id"`
	Author    		  UserList    	`json:"author"`
	AuthorId		  int 			`json:"author_id"`
	Content   		  string    	`json:"content"`
	Comments		  []CommentList	`json:"comments"`
	Title     		  string 		`json:"title"`
	Tags	  		  []Tag  		`json:"tags"`
	UpdatedAt 		  time.Time 	`json:"updated_at"`
}

type Article struct {
	ID 		  int       `json:"id, omitempty"`
	Author    int       `json:"author_id"`
	Title     string    `json:"title"`
	Tags	  []Tag		`json:"tags"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllArticles() []ArticleList {
	articles := []ArticleList{}

	// 选择所有的文章id
	row, err := utils.GetConn().Query("SELECT id FROM Articles")
	if err != nil {
		panic(err)
	}
	for row.Next() {
		articleId := -1
		err = row.Scan(&articleId)
		articles = append(articles, *GetArticleByID(articleId))
	}
	return articles
}

func CreateArticle(article Article) int {
	// 向Articles表添加条目
	stmt, err := utils.GetConn().Prepare("insert into Articles(title, content, createdAt, updatedAt, authorId) values (?, ?, ?, ?, ?)")
	if err != nil {
		panic("db insert prepare error")
	}
	res, err := stmt.Exec(article.Title, article.Content, time.Now(), time.Now(), article.Author)
	if err != nil {
		panic("db insert error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic("db insert error")
	}

	for _, v := range article.Tags {
		CreateTag(v.Content)

		// 向postTags表添加条目
		stmt, err = utils.GetConn().Prepare("insert into postTags(ArticleId, TagContent) values (?, ?)")
		if err != nil {
			panic("db insert prepare error")
		}
		_, err = stmt.Exec(id, v.Content)
		if err != nil {
			panic("db insert error")
		}
	}
	return int(id)
}

func GetArticleByID(id int) *ArticleList {
	article := ArticleList{}
	row, err := utils.GetConn().Query("SELECT id,title,content, updatedAt,authorId FROM Articles a WHERE a.id=?", id)
	if err != nil {
		panic(err)
	}
	for row.Next() {
		err = row.Scan(&article.ID, &article.Title, &article.Content, &article.UpdatedAt, &article.AuthorId)
		article.Author = *GetUserListByID(article.AuthorId)

		// 寻找文章拥有的tag
		row2, err := utils.GetConn().Query("SELECT t.content FROM Tags t, postTags p WHERE t.content=p.TagContent and p.ArticleId=?", article.ID)
		if err != nil {
			panic(err)
		}
		tags := []Tag{}
		for row2.Next() {
			tag := Tag{}
			err = row2.Scan(&tag.Content)
			tags = append(tags, tag)
		}
		article.Tags = tags

		// 寻找文章拥有的评论
		row3, err := utils.GetConn().Query("SELECT c.id, c.content, c.createdAt, c.articleId, c.authorId FROM Comments c WHERE c.articleId=?", article.ID)
		if err != nil {
			panic(err)
		}
		comments := []CommentList{}
		for row3.Next() {
			comment := CommentList{}
			err = row3.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.ArticleID, &comment.CreatorID)
			comment.Creator = *GetUserListByID(comment.CreatorID)
			comments = append(comments, comment)
		}
		article.Comments = comments
	}

	return &article
}

func GetArticleByUserID(id int) []ArticleList {
	articles := []ArticleList{}
	row, err := utils.GetConn().Query("SELECT a.id FROM Articles a WHERE a.authorId=?", id)
	if err != nil {
		panic(err)
	}
	for row.Next() {
		articleId := -1
		err = row.Scan(&articleId)
		articles = append(articles, *GetArticleByID(articleId))
	}
	return articles
}

func UpdateArticleByID(article Article) bool {
	stmt, err := utils.GetConn().Prepare("update Articles set title=?, content=?, updatedAt=? where id=?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(article.Title, article.Content, time.Now(), article.ID)
	if err != nil {
		panic(err)
	}

	stmt, err = utils.GetConn().Prepare("delete from postTags where ArticleId = ?")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	_, err = stmt.Exec(article.ID)
	if err != nil {
		panic("db update error")
	}


	for _, v := range article.Tags {
		CreateTag(v.Content)

		// 向postTags表添加条目
		stmt, err = utils.GetConn().Prepare("insert into postTags(ArticleId, TagContent) values (?, ?)")
		if err != nil {
			panic("db insert prepare error")
		}
		_, err = stmt.Exec(article.ID, v.Content)
		if err != nil {
			panic("db insert error")
		}
	}

	return true
}

func GetArticlesByTag(tag string) []ArticleList {
	articles := []ArticleList{}

	// 先从tag表中找到tagcontent
	row, err := utils.GetConn().Query("SELECT p.ArticleId FROM postTags p WHERE p.TagContent=?", tag)
	if err != nil {
		panic(err)
	}
	for row.Next() {
		articleId := -1
		err = row.Scan(&articleId)
		articles = append(articles, *GetArticleByID(articleId))
	}
	return articles
}

func GetArticlesPerPage(pageNum int) ArticlePage {
	articlePage := ArticlePage{}
	allArticles := GetAllArticles()
	articlePage.Number = len(allArticles)
	i := 0
	for range allArticles {
		if i >= (pageNum-1)*6 && i < pageNum * 6 {
			articlePage.Articles = append(articlePage.Articles, allArticles[i])
		}
		i = i + 1
	}
	return articlePage
}

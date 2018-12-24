package models

import (
	"fmt"
	"time"

	"github.com/GoProjectGroupForEducation/Go-Blog/utils"
)

type CommentList struct {
	ID 		  int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ArticleID int       `json:"article_id"`
	Creator   UserList  `json:"creator"`
	CreatorID int       `json:"creator_id"`
}

type Comment struct {
	ID 		  int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ArticleID int       `json:"article_id"`
	CreatorID int       `json:"creator_id"`
}

func CreateComment(comment *Comment) int {
	stmt, err := utils.GetConn().Prepare("insert into Comments values (?, ?, ?, ?, ? , ?)")
	if err != nil {
		panic("db insert prepare error")
	}
	res, err := stmt.Exec(nil, comment.Content, time.Now(), time.Now(), comment.ArticleID, comment.CreatorID)
	if err != nil {
		panic("db insert error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic("db insert error")
	}

	return int(id)
}

func GetAllCommentsByArticleID(articleID int) []CommentList {
	comments := []CommentList{}
	row, err := utils.GetConn().Query("SELECT id, content, createdAt, articleId, authorId FROM Comments c WHERE c.ArticleId = ? ", articleID)
	if err != nil {
		panic(err)
	}

	for row.Next() {
		comment := CommentList{}
		err = row.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.ArticleID, &comment.CreatorID)
		comments = append(comments, comment)
	}

	if err != nil {
		panic(err)
	}
	return comments
}

func UpdateCommentByID(comment Comment) bool {
	stmt, err := utils.GetConn().Prepare("update Comments set content=?, updatedAt=? where id=?")
	if err != nil {
		fmt.Println("error:", err)
	}
	_, err = stmt.Exec(comment.Content, time.Now(), comment.ArticleID)
	if err != nil {
		panic(err)
	}
	return true
}

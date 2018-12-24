package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB
var err error

func init() {
	DB, err = sql.Open("mysql", "blogdb:zsy2720a@tcp(172.17.0.1:3306)/blogdb?charset=utf8&loc=Asia%2FShanghai&parseTime=true")
	err = DB.Ping()
	for {
		if err == nil {
			break
		}
		fmt.Println("Trying...")
		time.Sleep(time.Duration(2)*time.Second)
		err = DB.Ping()
	}

	// Users
	users := `
    CREATE TABLE IF NOT EXISTS user(
        id INTEGER PRIMARY KEY AUTO_INCREMENT,
        username VARCHAR(255) UNIQUE NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL,
        iconPath TEXT NOT NULL
    );
    `


	_ ,err = DB.Exec(users)
	if err != nil {
		panic(err)
	}
	// UserRelations
	userRelations := `
    CREATE TABLE IF NOT EXISTS userRelations(
        UserId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE ,
        followerId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
		PRIMARY KEY (UserId, followerId)
    );
    `

	_ ,err = DB.Exec(userRelations)
	if err != nil {
		panic(err)
	}

	// Tags
	tags := `
    CREATE TABLE IF NOT EXISTS Tags(
        content VARCHAR(255) NOT NULL PRIMARY KEY
    );
    `
	_ ,err = DB.Exec(tags)
	if err != nil {
		panic(err)
	}

	// Articles
	articles := `
    CREATE TABLE IF NOT EXISTS Articles(
        id INTEGER  PRIMARY KEY AUTO_INCREMENT,
        title TEXT,
        content MEDIUMTEXT,
      	createdAt DATETIME,
      	updatedAt DATETIME,
      	authorId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
    );
    `
	_ ,err = DB.Exec(articles)
	if err != nil {
		panic(err)
	}

	// Comments
	comments := `
    CREATE TABLE IF NOT EXISTS Comments(
        id INTEGER  PRIMARY KEY AUTO_INCREMENT,
        content TEXT,
      	createdAt DATETIME,
      	updatedAt DATETIME,
      	articleId INTEGER NOT NULL REFERENCES Articles(id) ON DELETE CASCADE ON UPDATE CASCADE,
      	authorId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
    );
    `
	_ ,err = DB.Exec(comments)
	if err != nil {
		panic(err)
	}

	// postTags
	postTags := `
    CREATE TABLE IF NOT EXISTS postTags(
      	ArticleId INTEGER NOT NULL REFERENCES Articles(id) ON DELETE CASCADE ON UPDATE CASCADE,
      	TagContent VARCHAR(255) NOT NULL REFERENCES Tags(content) ON DELETE CASCADE ON UPDATE CASCADE,
		PRIMARY KEY (ArticleId, TagContent)
    );
    `

	_ ,err = DB.Exec(postTags)
	if err != nil {
		panic(err)
	}

	fmt.Println("Create database successfully!")
}

func GetConn() *sql.DB {
	return DB;
}


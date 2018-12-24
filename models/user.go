package models

import (
	"fmt"
	"github.com/GoProjectGroupForEducation/Go-Blog/utils"
)

type UserList struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	Email string 	`json:"email"`
	Followers []int `json:"followers"`
	Following []int `json:"following"`
	Iconpath string `json:"iconpath"`
}

type UserDetail struct {
	UserID   int    	`json:"id"`
	Username string 	`json:"username"`
	Email string 		`json:"email"`
	Followers []int 	`json:"followers"`
	Following []int 	`json:"following"`
	Articles  []ArticleList `json:"articles"`
	Iconpath string `json:"iconpath"`
}


type User struct {
	UserID    int       `json:"id, omitempty"`
	Username  string    `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
	Followers []int `json:"followers,[]"`
	Following []int `json:"following,[]"`
	Iconpath string `json:"iconpath"`
}

func CreateUser(user User) int {
	stmt, err := utils.GetConn().Prepare("insert into user(username, email, password, iconpath) values (?, ?, ?, ?)")
	if err != nil {
		panic("db insert prepare error")
	}
	res, err := stmt.Exec(user.Username, user.Email, user.Password, user.Iconpath)
	if err != nil {
		fmt.Println(err)
		panic("db insert error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic("db insert error")
	}

	return int(id)
}

func Follow(userid int, followerid int) bool {
	stmt, err := utils.GetConn().Prepare("insert into userRelations (UserId, followerId) values (?, ?)")
	if err != nil {
		panic("db insert prepare error")
	}
	fmt.Println(userid, followerid)
	_, err = stmt.Exec(userid, followerid)
	if err != nil {
		panic("db insert error")
	}
	return true
}

func Unfollow(userid int, followerid int) bool {
	stmt, err := utils.GetConn().Prepare("delete from userRelations where Userid=? and followerId=?")
	if err != nil {
		fmt.Println(err)
		panic("db insert prepare error")
	}
	_, err = stmt.Exec(userid, followerid)
	if err != nil {
		panic("db insert error")
	}
	return true
}

func GetUserByID(id int) *User {
	user := User{}
	row, err := utils.GetConn().Query("SELECT * FROM user WHERE user.id = ? ", id)
	if err != nil {
		fmt.Println("error:", err)
	}
	user.Following = []int{};
	user.Followers = []int{};
	if row.Next()  {
		err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Iconpath)
		for _, v := range GetUserFollowers(id) {
			user.Followers = append(user.Followers, v.UserID)
		}
		for _, v := range GetUserFollowing(id) {
			user.Following = append(user.Following, v.UserID)
		}
	}

	if err != nil {
		panic(err)
	}
	return &user
}

func GetUserDetailByID(id int) *UserDetail {
	user := UserDetail{}
	row, err := utils.GetConn().Query("SELECT id, username, email, iconpath FROM user WHERE user.id = ? ", id)
	if err != nil {
		panic(err)
	}
	user.Following = []int{};
	user.Followers = []int{};
	if row.Next() {
		err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.Iconpath)
		for _, v := range GetUserFollowers(id) {
			user.Followers = append(user.Followers, v.UserID)
		}
		for _, v := range GetUserFollowing(id) {
			user.Following = append(user.Following, v.UserID)
		}
		user.Articles = GetArticleByUserID(id)
		if err != nil {
			panic(err)
		}
		return &user
	}
	return nil
}

func GetUserByUsername(username string) *User {
	user := User{}
	id := -1
	row, err := utils.GetConn().Query("SELECT * FROM user WHERE user.username = ? ", username)
	if err != nil {
		fmt.Println("error:", err)
	}
	user.Following = []int{};
	user.Followers = []int{};
	if row.Next()  {
		err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Iconpath)
		id = user.UserID
		for _, v := range GetUserFollowers(id) {
			user.Followers = append(user.Followers, v.UserID)
		}
		for _, v := range GetUserFollowing(id) {
			user.Following = append(user.Following, v.UserID)
		}
	}

	if err != nil {
		panic(err)
	}

	if id == -1 {
		return nil
	}

	return &user
}

func GetUserListByID(id int) *UserList {
	user := UserList{}
	row, err := utils.GetConn().Query("SELECT id, username, email, iconpath FROM user WHERE user.id = ? ", id)
	if err != nil {
		panic(err)
	}
	user.Following = []int{};
	user.Followers = []int{};
	if row.Next() {
		err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.Iconpath)
		for _, v := range GetUserFollowers(id) {
			user.Followers = append(user.Followers, v.UserID)
		}
		for _, v := range GetUserFollowing(id) {
			user.Following = append(user.Following, v.UserID)
		}
	}

	if err != nil {
		panic(err)
	}

	return &user
}

func UpdateUserByID(id int, user User) bool {
	stmt, err := utils.GetConn().Prepare("update user set username=?, email=?, password=? where id=?")
	if err != nil {
		fmt.Println("error:", err)
	}
	_, err = stmt.Exec(user.Username, user.Email, user.Password, id)
	if err != nil {
		panic(err)
	}

	return true
}

func UpdateUserIconByID(id int, iconPath string) bool {
	stmt, err := utils.GetConn().Prepare("update user set iconPath=? where id=?")
	if err != nil {
		fmt.Println("error:", err)
	}
	_, err = stmt.Exec(iconPath, id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return true
}

func GetUserFollowing(id int) []UserList {
	users := []UserList{}

	rows, err := utils.GetConn().Query("SELECT B.id, B.username, B.email, B.iconPath FROM userRelations A, user B WHERE A.followerId = ? and A.UserId = B.id", id)
	if err != nil {
		fmt.Println("error:", err)
	}
	for rows.Next() {
		user := UserList{}
		err = rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Iconpath)
		users = append(users, user)
	}

	return users
}

func GetUserFollowers(id int) []UserList {
	users := []UserList{}

	rows, err := utils.GetConn().Query("SELECT B.id, B.username, B.email, B.iconPath FROM userRelations A, user B WHERE A.UserId = ? and A.followerId = B.id", id)
	if err != nil {
		fmt.Println("error:", err)
	}
	for rows.Next() {
		user := UserList{}
		err = rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Iconpath)
		users = append(users, user)
	}

	return users
}
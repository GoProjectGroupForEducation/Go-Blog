package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/GoProjectGroupForEducation/Go-Blog/services"
	"github.com/GoProjectGroupForEducation/Go-Blog/utils"
	"github.com/gorilla/mux"
	"github.com/GoProjectGroupForEducation/Go-Blog/models"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	Filename string `json:"imgpath"`
}

func DownloadFile(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error  {
	vars := mux.Vars(req)

	filename := strings.Split(vars["filename"],".")

	file, _, err := req.FormFile("uploadFile")
	if err != nil {
		return utils.SendData(w, "{}","INVALID_FILE", http.StatusBadRequest)
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return utils.SendData(w, "{}","INVALID_FILE", http.StatusBadRequest)
	}
	//判定文件类型
	//filetype := http.DetectContentType(fileBytes)

	//产生随机的token用作文件名
	//var checksum = []byte("heng-is-a-very-handsome-boy")
	//h := md5.New()
	//token := fmt.Sprintf("%x", h.Sum(checksum))

	//要登录后才能使用上传文件
	user := services.GetCurrentUser(req.Header.Get("Authorization"))
	timeStr:=time.Now().Format("20060102150405")  //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	fileName := user.Username + "at" + timeStr + "." + filename[1]

	models.UpdateUserIconByID(user.UserID, fileName)
	//fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		return utils.SendData(w, "{}","CANT_READ_FILE_TYPE", http.StatusInternalServerError)
	}
	root, _ := os.Getwd()
	downloadFilePath := root + "/static/"

	newPath := downloadFilePath + fileName

	// fmt.Printf("FileType: %s, File: %s\n", filename[1], newPath)
	newFile, err := os.Create(newPath)
	if err != nil {
		return utils.SendData(w, "{}","CANT_WRITE_FILE", http.StatusInternalServerError)
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		return utils.SendData(w, "{}","CANT_WRITE_FILE", http.StatusInternalServerError)
	}

	fileInfo := FileInfo{fileName}
	//getMultiPart(req)
	data, err := json.Marshal(fileInfo)

	if err != nil {
		return utils.SendData(w, "{}","Error when marshal!", http.StatusInternalServerError)
	}
	return utils.SendData(w, string(data), "Upload file successfully", http.StatusOK)
}

func DownloadPostFile(w http.ResponseWriter, req *http.Request, next utils.NextFunc) error  {
	vars := mux.Vars(req)

	filename := strings.Split(vars["filename"],".")

	file, _, err := req.FormFile("uploadFile")
	if err != nil {
		return utils.SendData(w, "{}","INVALID_FILE", http.StatusBadRequest)
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return utils.SendData(w, "{}","INVALID_FILE", http.StatusBadRequest)
	}
	//判定文件类型
	//filetype := http.DetectContentType(fileBytes)

	//产生随机的token用作文件名
	//var checksum = []byte("heng-is-a-very-handsome-boy")
	//h := md5.New()
	//token := fmt.Sprintf("%x", h.Sum(checksum))

	//要登录后才能使用上传文件
	user := services.GetCurrentUser(req.Header.Get("Authorization"))
	timeStr:=time.Now().Format("20060102150405")  //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	fileName := user.Username + "at" +timeStr + "." + filename[1]

	//fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		return utils.SendData(w, "{}","CANT_READ_FILE_TYPE", http.StatusInternalServerError)
	}
	root, _ := os.Getwd()
	downloadFilePath := root + "/static/"

	newPath := downloadFilePath + fileName

	newFile, err := os.Create(newPath)
	if err != nil {
		return utils.SendData(w, "{}","CANT_WRITE_FILE", http.StatusInternalServerError)
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		return utils.SendData(w, "{}","CANT_WRITE_FILE", http.StatusInternalServerError)
	}

	fileInfo := FileInfo{fileName}
	//getMultiPart(req)
	data, err := json.Marshal(fileInfo)

	if err != nil {
		return utils.SendData(w, "{}","Error when marshal!", http.StatusInternalServerError)
	}
	return utils.SendData(w, string(data), "Upload file successfully", http.StatusOK)
}


//字节解析multi-part

//通过MultipartReader

func getMultiPart(r *http.Request)()  {



	mr,err := r.MultipartReader()

	if err != nil{

		fmt.Println("r.MultipartReader() err,",err)

		return

	}

	form ,_ := mr.ReadForm(128)

	getFormData(form)

}


func getFormData(form *multipart.Form) {
	//获取 multi-part/form body中的form value
	for k, v := range form.Value {
		fmt.Println("value,k,v = ", k, ",", v)
	}

	fmt.Println()
	//获取 multi-part/form中的文件数据
	for _, v := range form.File {
		for i := 0; i < len(v); i++ {
			f, _ := v[i].Open()
			buf, _ := ioutil.ReadAll(f)
			fmt.Println("file-content", string(buf))
			fmt.Println()
		}
	}
}
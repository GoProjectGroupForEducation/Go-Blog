package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type sendDataType struct {
	status int
	msg    string
	data   string
	time   string
}

func GetStoreTimeNow() string {
	time := time.Now().UnixNano() / 1e6
	return fmt.Sprintf("%v", time)
	// return time
}

func GetISOTimeNow() string {
	return time.Now().Format(time.RFC3339)
}

func SendData(w http.ResponseWriter, data string, msg string, status int) error {
	var buff []byte
	header := w.Header()
	header.Add("Content-type", "application/json")
	res := sendDataType{status, msg, data, time.Now().Format(time.RFC3339)}
	if(status != 200){
		strbuff := strings.Join([]string{`{"status":`, strconv.Itoa(res.status), `,"msg":"`, res.msg, `","data":`, "{}", `,"time":"`, res.time, `"}`}, "")
		buff = make([]byte, len(strbuff))
		for i, one := range strbuff {
			buff[i] = byte(one)
		}
	}else {
		strbuff := strings.Join([]string{`{"status":`, strconv.Itoa(res.status), `,"msg":"`, res.msg, `","data":`, res.data, `,"time":"`, res.time, `"}`}, "")
		buff = make([]byte, len(strbuff))

		for i, one := range strbuff {
			buff[i] = byte(one)
		}
		fmt.Println(buff)
	}
	w.WriteHeader(status)
	_, err := w.Write(buff)
	if err != nil {
		return err
	}
	return nil
}
package controller

import (
	"encoding/json"
)

type Controller struct {
}

// NewApp creates a new App application struct
func NewController() *Controller {
	return &Controller{}
}

type resp struct {
	Data any    `json:"data"`
	Err  string `json:"err"`
}

func Resp(data any, err error) string {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	result := resp{
		Data: data,
		Err:  errMsg,
	}
	//fmt.Println("api resp:", data, err)
	js, _ := json.Marshal(result)
	return string(js)
}

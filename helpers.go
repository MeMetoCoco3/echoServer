package main

const msgKey = "msg"

type Response struct {
	Content    interface{}
	IsLoggedIn bool
}

func SetIsLogged(user interface{}, respose *Response) {
	if user != nil {
		respose.IsLoggedIn = true
	}
}

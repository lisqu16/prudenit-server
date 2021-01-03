package user

type User struct {
	Email		string `json:"email"`
	Username	string `json:"username"`
}

type Response struct {
	Ok		bool `json:"ok"`
	Msgs 	[]string `json:"messages"`	
	Data 	map[string]interface{} `json:"data"`
}
package forms

type LoginForm struct {
	Username string
	Password string
}

type GetCode struct {
	Email string `json:"email,omitempty"`
}
type CheckCode struct {
	Email string `json:"email,omitempty"`
	Code  string `json:"code"`
}
type Form struct {
	Text string `json:"text"`
}

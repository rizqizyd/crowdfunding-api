package helper

import "github.com/go-playground/validator/v10"

type Response struct {
	Meta Meta        `json:"meta"` // strukturnya paten, beri penanda `json:"meta"` supaya outputnya huruf kecil semua
	Data interface{} `json:"data"` // nilainya bisa flexibel
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	// menangani validasi error, membuat array dan menambah data array (error) melalui perulangan
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

/*
Output Response:
meta: {
	message: "Your account has been created",
	code: 200,
	status: "success"
},
data: {
	id: 1,
	name: "MasonMount",
	occupation: "content creator",
	email: "mount@gmail.com",
	token: "masonmount"
}
*/

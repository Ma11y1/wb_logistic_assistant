package response

import (
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

type CaptchaGetTask models.CaptchaTask

type CaptchaVerifyAnswer struct {
	Token string `json:"wb-captcha-short-token"`
}

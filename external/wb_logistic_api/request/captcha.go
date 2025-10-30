package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// CaptchaGetTaskRequest Gets the captcha task
//
// URL: https://pow.wildberries.ru/api/v1/short/get-task
type CaptchaGetTaskRequest struct {
	BaseRequest
	infoJWTToken string
}

func NewCaptchaGetTaskRequest(client *transport.BaseHTTPClient) *CaptchaGetTaskRequest {
	return &CaptchaGetTaskRequest{BaseRequest: *NewRequest(client, "https://pow.wildberries.ru/api/v1/short/get-task")}
}

func (r *CaptchaGetTaskRequest) Do(ctx context.Context) (response response.CaptchaGetTask, err error) {
	res, err := r.PostData(ctx, bytes.NewBuffer([]byte(r.infoJWTToken)))
	if err != nil {
		err = fmt.Errorf("CaptchaGetTaskRequest.Do: %s", err)
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		err = fmt.Errorf("CaptchaGetTaskRequest.Do: error decoding JSON %s", err)
	}

	return
}

func (r *CaptchaGetTaskRequest) ClientID(id string) {
	r.queryParameters.Set("client_id", id)
}

// BrowserInfo Accepts JWT token received from models.CaptchaBrowserInfo
func (r *CaptchaGetTaskRequest) BrowserInfo(token string) {
	r.infoJWTToken = token
}

// CaptchaVerifyAnswerRequest Submits a task with an answer and receives a captcha token
//
// URL: https://pow.wildberries.ru/api/v1/short/verify-answer
type CaptchaVerifyAnswerRequest struct {
	BaseRequest
}

func NewCaptchaVerifyAnswerRequest(client *transport.BaseHTTPClient) *CaptchaVerifyAnswerRequest {
	return &CaptchaVerifyAnswerRequest{BaseRequest: *NewRequest(client, "https://pow.wildberries.ru/api/v1/short/verify-answer")}
}

func (r *CaptchaVerifyAnswerRequest) Do(ctx context.Context) (response response.CaptchaGetTask, err error) {
	err = r.PostUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("CaptchaVerifyAnswerRequest.Do: %s", err)
	}
	return
}

func (r *CaptchaVerifyAnswerRequest) Answers(answers []int) {
	r.parameters.Set("answers", answers)
}

func (r *CaptchaVerifyAnswerRequest) Task(task *models.CaptchaTask) {
	r.parameters.Set("task", task)
}

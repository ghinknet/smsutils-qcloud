package qcloud

import (
	"github.com/ghinknet/smsutils/v3/errors"
	"github.com/ghinknet/smsutils/v3/model"
	"github.com/ghinknet/smsutils/v3/utils"
	"github.com/ghinknet/toolbox/pointer"
	smsv20210111 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func (c Client) SendMessage(dest string, sender string, template string, vars model.Vars) error {
	// Try to parse number
	dest, _, _, _, err := utils.ProcessNumberForChinese(dest)

	// Preprocess vars
	params := make([]*string, len(vars))
	for k, v := range vars {
		params[k] = &v.Value
	}

	// Construct a request
	request := smsv20210111.NewSendSmsRequest()

	// Set request params
	request.PhoneNumberSet = []*string{&dest}
	request.SmsSdkAppId = &c.SmsSdkAPPID
	request.TemplateId = &template
	request.SignName = &sender
	request.TemplateParamSet = params

	// Send requests
	response, err := c.Client.SendSms(request)
	if err != nil {
		return err
	}
	if response != nil && response.Response != nil {
		for _, status := range response.Response.SendStatusSet {
			if pointer.SafeDeref(status.Code) != "Ok" {
				return errors.ErrDriverSendFailed.
					WithDriverName(Name).
					WithDriverCode(pointer.SafeDeref(status.Code)).
					WithDriverMessage(pointer.SafeDeref(status.Message)).
					WithDriverRequestID(pointer.SafeDeref(response.Response.RequestId)).
					WithDriverResponse(response.Response)
			}
		}
	} else {
		return errors.ErrDriverSendFailed.
			WithDriverName(Name).
			WithDriverResponse(response)
	}

	return nil
}

package claude

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"one-api/dto"
	relaychannel "one-api/relay/channel"
	relaycommon "one-api/relay/common"
	"one-api/service"
)

type Adaptor struct {
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo, request dto.GeneralOpenAIRequest) {

}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	return fmt.Sprintf("%s/v1/complete", info.BaseUrl), nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, info *relaycommon.RelayInfo) error {
	relaychannel.SetupApiRequestHeader(info, c, req)
	req.Header.Set("x-api-key", info.ApiKey)
	anthropicVersion := c.Request.Header.Get("anthropic-version")
	if anthropicVersion == "" {
		anthropicVersion = "2023-06-01"
	}
	req.Header.Set("anthropic-version", anthropicVersion)
	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	return request, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error) {
	return relaychannel.DoApiRequest(a, c, info, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage *dto.Usage, err *dto.OpenAIErrorWithStatusCode) {
	if info.IsStream {
		var responseText string
		err, responseText = claudeStreamHandler(c, resp)
		usage = service.ResponseText2Usage(responseText, info.UpstreamModelName, info.PromptTokens)
	} else {
		err, usage = claudeHandler(c, resp, info.PromptTokens, info.UpstreamModelName)
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return ChannelName
}
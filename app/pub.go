package app

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	logBeat "github.com/Ankr-network/dccn-common/protos/logbeat"
	"go.uber.org/zap"
)

type requestBodyCon struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int64  `json:"id"`
	Method  string `json:"method"`
}

var (
	rgx = regexp.MustCompile(`\{(.*?)\}`)
)

func TransferMetricsToProtobuf(logText string) *logBeat.MetricsV2 {

	var (
		remoteAddr  string
		timeLocal   string
		request     string
		status      string
		requestBody string
		apiId       string
	)

	resList := strings.Split(logText, " ")
	if len(resList) >= 15 {
		remoteAddr = resList[0]
		timeLocal = resList[4] + " " + resList[5]
		request = resList[7] + " " + resList[8] + " " + resList[9]
		status = resList[10]
	} else {
		return nil
	}

	requestBodyList := rgx.FindStringSubmatch(logText)
	if len(request) == 0 {
		return nil
	}
	requestBody = requestBodyList[0]
	requestBodyStr := strings.Replace(requestBody, "\\x22", `"`, -1)
	var requestBodyRes requestBodyCon
	fmt.Println("requestBodyStr", requestBodyStr)
	err := json.Unmarshal([]byte(requestBodyStr), &requestBodyRes)
	if err != nil {
		log.Error("Pub", zap.Error(err))
		return nil
	}
	fmt.Println("requestBodyRes", requestBodyRes)

	sentTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeLocal)
	if err != nil {
		log.Error("Pub", zap.Error(err))
		return nil
	}

	requestList := strings.Split(request, `/`)
	if len(requestList) >= 5 && len(requestList[1]) == 32 {
		apiId = requestList[1]
	}

	statusInt64, err := strconv.ParseInt(status, 10, 64)
	if err != nil {
		log.Error("Pub", zap.Error(err))
		return nil
	}

	res := &logBeat.MetricsV2{
		XReadIp:    remoteAddr,
		Index:      sentTime.UnixNano(),
		XUserID:    apiId,
		MethodName: requestBodyRes.Method,
		Request:    []byte(requestBodyStr),
		JsonrpcID:  requestBodyRes.Id,
		Code:       int32(statusInt64),
	}

	fmt.Println(res)

	return res
}

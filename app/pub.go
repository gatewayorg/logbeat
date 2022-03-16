package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	logBeat "github.com/Ankr-network/dccn-common/protos/logbeat"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
	if len(requestBodyList) == 0 {
		return nil
	}
	requestBody = requestBodyList[0]
	requestBodyStr := strings.Replace(requestBody, "\\x22", `"`, -1)
	var requestBodyRes requestBodyCon
	// fmt.Println("requestBodyStr", requestBodyStr)
	err := json.Unmarshal([]byte(requestBodyStr), &requestBodyRes)
	if err != nil {
		log.Error("Pub", zap.Error(err))
		return nil
	}
	// fmt.Println("requestBodyRes", requestBodyRes)

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

	log.Info("Pub", zap.Any("pub metrics", res))

	return res
}

type PubMetrics struct {
	ProducerMap  map[string]*nsq.Producer
	MetricsTopic string
}

func NewPubMetrics(nsqdAddress []string) *PubMetrics {
	producerMap := make(map[string]*nsq.Producer)
	// init nsqd
	if len(nsqdAddress) == 0 {
		log.Error("Pub", zap.Error(errors.New("no nsqd")))
		return nil
	}
	config := nsq.NewConfig()
	for _, address := range nsqdAddress {
		fmt.Println("address:", address)
		produce, _ := nsq.NewProducer(address, config)
		err := produce.Ping()
		if err == nil {
			producerMap[address] = produce
		} else {
			log.Error("Pub", zap.Error(err))
		}
	}
	log.Info("Pub", zap.Any("init mq success", producerMap))

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			pingAndkeepalive(producerMap)
		}
	}()

	return &PubMetrics{
		ProducerMap:  producerMap,
		MetricsTopic: logBeat.LogbeatMetricsTopic,
	}

}

func pingAndkeepalive(producerMap map[string]*nsq.Producer) {
	// ping and keep live
	log.Info("keep alive")
	for addressKey, producerConn := range producerMap {
		err := producerConn.Ping()
		if err != nil {
			log.Error("Pub", zap.Error(err))
			delete(producerMap, addressKey)
		}
	}
	time.Sleep(5 * time.Minute)

}

func (pub *PubMetrics) ProducerPub(message *logBeat.MetricsV2) error {
	if len(pub.ProducerMap) != 0 {

		r := rand.Intn(len(pub.ProducerMap))
		for k, conn := range pub.ProducerMap {
			if r == 0 {
				if message != nil {
					body, err := proto.Marshal(message)
					if err != nil {
						log.Error("Pub", zap.Error(err))
						return err
					}
					err = conn.Publish(pub.MetricsTopic, body)
					if err != nil {
						log.Error("Pub", zap.Error(err))
						return err
					}
					log.Info(fmt.Sprintf("Pub success, meaage is %v, mq address is %s", message, k))
				}
			}
			r--
		}

	}
	return nil
}

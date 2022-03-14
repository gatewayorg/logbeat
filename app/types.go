package app

type AccessLog struct {
	remoteAddr           string
	remoteUser           string
	timeLocal            string
	request              string
	status               string
	bodyBytesSent        string
	httpReferer          string
	httpUserAgent        string
	httpXForwardedFor    string
	upStreamAddr         string
	requestTime          string
	upstreamResponseTime string
	upstreamAddr         string
	requestBody          string
}

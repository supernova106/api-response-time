package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	// MongoDB collection
	CollectionAccessLog = "nginx"
)

// AccessLog model
type AccessLog struct {
	Id          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	UserIp      string        `json:"user_ip" bson:"user_ip"`
	APIName     string        `json:"apiname" bson:"apiname"`
	Request     string        `json:"request" bson:"request"`
	Code        int64         `json:"code" bson:"code"`
	RequestTime float64       `bson:"request_time"`
	Time        time.Time     `bson:"time"`
}

type AccessLogs []AccessLog

// AccesslogSummary
type AccessLogSummary struct {
	APIName         string  `json:"_id" bson:"_id"`
	AvgResponseTime float64 `bson:"avgresponsetime"`
}

type AccessLogSummaries []AccessLogSummary

type StatusCodeSummary struct {
	Code  int64 `json:"_id" bson:"_id"`
	Count int64 `bson:"count"`
}

type StatusCodeSummaries []StatusCodeSummary

type RequestRateSummary struct {
	APIName string `json:"_id" bson:"_id"`
	Count   int64  `bson:"count"`
}

type RequestRateSummaries []RequestRateSummary

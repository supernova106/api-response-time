package accesslog

import (
	"api-response-time/app/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

// List all accesslog
func GetAll(c *gin.Context) {
	db := c.MustGet("mdb").(*mgo.Database)

	apiName := c.DefaultQuery("apiname", "all")
	timeRange := c.DefaultQuery("time", "3")
	status := c.DefaultQuery("status", "false")
	timeint, _ := strconv.ParseInt((timeRange), 10, 64)
	now := time.Now()
	timestamp := now.Add(-time.Duration(timeint) * time.Minute)

	var operations []bson.M

	if status == "code" {
		statusCodeSummaries := models.StatusCodeSummaries{}
		if apiName == "all" {
			operations = []bson.M{
				{"$match": bson.M{"time": bson.M{"$gt": timestamp}}},
				{"$group": bson.M{"_id": "$code", "count": bson.M{"$sum": 1}}},
			}
		} else {
			operations = []bson.M{
				{"$match": bson.M{"apiname": apiName, "time": bson.M{"$gt": timestamp}}},
				{"$group": bson.M{"_id": "$code", "count": bson.M{"$sum": 1}}},
			}
		}

		_ = db.C(models.CollectionAccessLog).Pipe(operations).All(&statusCodeSummaries)
		c.JSON(200, statusCodeSummaries)
	} else if status == "rate" {
		requestRateSummaries := models.RequestRateSummaries{}
		if apiName == "all" {
			operations = []bson.M{
				{"$match": bson.M{"time": bson.M{"$gt": timestamp}}},
				{"$group": bson.M{"_id": "$apiname", "count": bson.M{"$sum": 1}}},
			}
		} else {
			operations = []bson.M{
				{"$match": bson.M{"apiname": apiName, "time": bson.M{"$gt": timestamp}}},
				{"$group": bson.M{"_id": "$apiname", "count": bson.M{"$sum": 1}}},
			}
		}

		_ = db.C(models.CollectionAccessLog).Pipe(operations).All(&requestRateSummaries)
		c.JSON(200, requestRateSummaries)
	} else {
		accessLogSummaries := models.AccessLogSummaries{}
		if apiName == "all" {
			operations = []bson.M{
				{"$match": bson.M{"time": bson.M{"$gt": timestamp}, "code": bson.M{"$lt": 500}}},
				{"$group": bson.M{"_id": "$apiname", "avgresponsetime": bson.M{"$avg": "$request_time"}}},
			}
		} else {
			operations = []bson.M{
				{"$match": bson.M{"apiname": apiName, "time": bson.M{"$gt": timestamp}, "code": bson.M{"$lt": 500}}},
				{"$group": bson.M{"_id": "$apiname", "avgresponsetime": bson.M{"$avg": "$request_time"}}},
			}
		}

		_ = db.C(models.CollectionAccessLog).Pipe(operations).All(&accessLogSummaries)
		c.JSON(200, accessLogSummaries)
	}
}

func Check(c *gin.Context) {
	c.String(200, "Hello! It's running!")
}

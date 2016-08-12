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

	if status == "true" {
		if apiName == "all" {
			c.JSON(400, gin.H{"status_code": 400, "error_msg": "apiname is not specified"})
			return
		}
		statusCodeSummaries := models.StatusCodeSummaries{}
		operations = []bson.M{
			{"$match": bson.M{"apiname": apiName, "time": bson.M{"$gt": timestamp}}},
			{"$group": bson.M{"_id": "$code", "count": bson.M{"$sum": 1}}},
		}
		_ = db.C(models.CollectionAccessLog).Pipe(operations).All(&statusCodeSummaries)
		c.JSON(200, statusCodeSummaries)
	} else {
		accessLogSummaries := models.AccessLogSummaries{}
		if apiName == "all" {
			operations = []bson.M{
				{"$match": bson.M{"time": bson.M{"$gt": timestamp}}},
				{"$group": bson.M{"_id": "$apiname", "avgresponsetime": bson.M{"$avg": "$request_time"}}},
			}
		} else {
			operations = []bson.M{
				{"$match": bson.M{"apiname": apiName, "time": bson.M{"$gt": timestamp}}},
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

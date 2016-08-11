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
	timeint, _ := strconv.ParseInt((timeRange), 10, 64)
	accessLogSummaries := models.AccessLogSummaries{}
	now := time.Now()
	timestamp := now.Add(-time.Duration(timeint) * time.Minute)

	var operations []bson.M

	if apiName == "all" {
		operations = []bson.M{
			{"$match": bson.M{"time": bson.M{"$gt": timestamp}}},
			{"$group": bson.M{"_id": "$apiname", "avgresponsetime": bson.M{"$avg": "$request_time"}}},
		}
	} else {
		operations = []bson.M{
			{"$match": bson.M{"time": bson.M{"$gt": timestamp}, "apiname": apiName}},
			{"$group": bson.M{"_id": "$apiname", "avgresponsetime": bson.M{"$avg": "$request_time"}}},
		}
	}

	_ = db.C(models.CollectionAccessLog).Pipe(operations).All(&accessLogSummaries)
	c.JSON(200, accessLogSummaries)
}

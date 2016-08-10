package accesslog

import (
	"api-response-time/app/db"
	"api-response-time/app/models"
	"github.com/gin-gonic/gin"
	//	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// List all accesslog
func List(c *gin.Context) {
	db := c.MustGet("mdb").(*db.DB)

	d := db.Copy()
	defer d.Close()

	//accessLog := models.AccessLog{}
	//query := bson.M{"user_ip": "184.72.68.45"}
	//_ = d.FindOne(CollectionAccessLog, query, &accessLog)
	accessLogSummary := models.AccessLogSummary{}
	o1 := []bson.M{"$match": bson.M{"time": bson.M{"$gt": bson.Now()}}}
	o2 := []bson.M{"$group": bson.M{"_id": "$apiname", "avgResponseTime": bson.M{"$avg": "$request_time"}}}
	operations := []bson.M{o1, o2}
	_ = d.Pipe("nginx", operations, &accessLogSummary)
	c.JSON(200, accessLogSummary)

}

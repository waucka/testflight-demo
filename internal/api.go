package api

import (
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
	"net/http"
)

type Config struct {
	session *mgo.Session
	db *mgo.Database
	usercoll *mgo.Collection
	chancoll *mgo.Collection
}

func NewConfig(session *mgo.Session, dbname string) *Config {
	db := session.DB(dbname)
	return &Config{
		session,
		db,
		db.C("users"),
		db.C("channels"),
	}
}

func (self *Config) GetRouter() *gin.Engine {
	router := gin.New()

	router.Use(AjaxErrorGuard())
	router.Use(MiddlewareAuth(self.usercoll))
	router.GET("/channel", self.GetChannelList)
	router.GET("/channel/:slug", self.GetChannelInfo)
	router.GET("/channel/:slug/item", self.GetChannelItemList)
	router.POST("/channel/:slug/item", self.CreateChannelItem)
	router.GET("/channel/:slug/item/:itemSlug", self.GetChannelItem)

	return router
}

func getUsername(c *gin.Context) string {
	username, err := c.Get("USERNAME")
	if err != nil {
		return ""
	}
	return username.(string)
}

func (self *Config) GetChannelList(c *gin.Context) {
	_ = forceAuth(c)
	chanIter := self.chancoll.Find(nil).Iter()
	channelData := make(map[string]string)
	var chanRec ChannelDBRecord
	for chanIter.Next(&chanRec) {
		channelData[chanRec.Slug] = chanRec.Title
	}
	c.JSON(http.StatusOK, channelData)
}

func (self *Config) GetChannelInfo(c *gin.Context) {
	InternalError("Not implemented")
}

func (self *Config) GetChannelItemList(c *gin.Context) {
	InternalError("Not implemented")
}

func (self *Config) CreateChannelItem(c *gin.Context) {
	InternalError("Not implemented")
}

func (self *Config) GetChannelItem(c *gin.Context) {
	InternalError("Not implemented")
}

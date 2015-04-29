package api

import (
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
	"net/http"
	"time"
)

type Config struct {
	session  *mgo.Session
	db       *mgo.Database
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
	router.POST("/channel", self.CreateChannel)
	router.GET("/channel/:slug", self.GetChannelInfo)
	router.GET("/channel/:slug/item", self.GetChannelItemList)
	router.POST("/channel/:slug/item", self.CreateChannelItem)
	router.GET("/channel/:slug/item/:itemSlug", self.GetChannelItem)

	return router
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

func (self *Config) CreateChannel(c *gin.Context) {
	_ = forceAuth(c)
	slug := c.Request.FormValue("slug")
	title := c.Request.FormValue("title")
	if slug == "" {
		BadRequest("slug cannot be empty")
	}
	if title == "" {
		BadRequest("title cannot be empty")
	}
	err := self.chancoll.Insert(&ChannelDBRecord{
		Slug:  slug,
		Title: title,
		Items: make([]ItemDBRecord, 0),
	})
	if err == nil {
		c.String(http.StatusNoContent, "")
	} else {
		BadRequest(err.Error())
	}
}

func (self *Config) GetChannelInfo(c *gin.Context) {
	_ = forceAuth(c)
	slug := c.Params.ByName("slug")
	var chanRec ChannelDBRecord
	err := self.chancoll.FindId(slug).One(&chanRec)
	if err == mgo.ErrNotFound {
		NotFound("No such channel " + slug)
	} else if err != nil {
		InternalError("Could not fetch channel info from database")
	}
	c.JSON(http.StatusOK, chanRec.ToJSON())
}

func (self *Config) GetChannelItemList(c *gin.Context) {
	_ = forceAuth(c)

	slug := c.Params.ByName("slug")
	var chanRec ChannelDBRecord
	err := self.chancoll.FindId(slug).One(&chanRec)
	if err == mgo.ErrNotFound {
		NotFound("No such channel " + slug)
	} else if err != nil {
		InternalError("Could not fetch channel info from database")
	}

	itemData := make(map[string]string)
	for _, item := range chanRec.Items {
		itemData[item.Slug] = item.Title
	}
	c.JSON(http.StatusOK, itemData)
}

func (self *Config) CreateChannelItem(c *gin.Context) {
	username := forceAuth(c)
	chanSlug := c.Params.ByName("slug")
	title := c.Params.ByName("title")
	b64data := c.Params.ByName("b64data")
	itemSlug := c.Request.FormValue("itemSlug")
	var chanrec ChannelDBRecord
	err := self.chancoll.FindId(chanSlug).One(&chanrec)
	if err != nil {
		InternalError("Cannot fetch channel info from database")
	}
	itemrec := &ItemDBRecord{
		Slug:         itemSlug,
		Title:        title,
		DateUploaded: time.Now(),
		Data:         b64data,
		Uploader:     username,
	}
	chanrec.Items = append(chanrec.Items, *itemrec)
	err = self.chancoll.UpdateId(chanSlug, chanrec)
	if err != nil {
		InternalError("Cannot update channel info in database")
	}
	c.String(http.StatusOK, "/channel/"+chanSlug+"/item/"+itemSlug)
}

func (self *Config) GetChannelItem(c *gin.Context) {
	_ = forceAuth(c)

	slug := c.Params.ByName("slug")
	itemSlug := c.Params.ByName("itemSlug")
	var chanRec ChannelDBRecord
	err := self.chancoll.FindId(slug).One(&chanRec)
	if err == mgo.ErrNotFound {
		NotFound("No such channel " + slug)
	} else if err != nil {
		InternalError("Could not fetch channel info from database")
	}

	for _, item := range chanRec.Items {
		if item.Slug == itemSlug {
			c.JSON(http.StatusOK, item.ToJSON())
			return
		}
	}
	NotFound("Channel " + slug + " has no item " + itemSlug)
}

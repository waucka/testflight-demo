package api

import (
	"github.com/drewolson/testflight"
	"labix.org/v2/mgo"
	. "gopkg.in/check.v1"
	"time"
	"io"
	"io/ioutil"
	"os"
	"net/http"
	"path/filepath"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ApiSuite struct{
	apiConfig *Config
	tokens map[string]string
	chan1Rec *ChannelJSONRecord
	chan2Rec *ChannelJSONRecord
	item1Rec *ItemJSONRecord
	item2Rec *ItemJSONRecord
	user1 *UserJSONRecord
	user2 *UserJSONRecord
	baduser1 *UserJSONRecord
}

var _ = Suite(&ApiSuite{
	nil,
	make(map[string]string),
	nil,
	nil,
	nil,
	nil,
	&UserJSONRecord{
		Username: "testuser1",
		Subscriptions: make([]string, 0),
	},
	&UserJSONRecord{
		Username: "testuser2",
		Subscriptions: make([]string, 0),
	},
	&UserJSONRecord{
		Username: "baduser1",
		Subscriptions: make([]string, 0),
	},
})
var (
	zeroTime time.Time
)

func createUser(username string, usercoll *mgo.Collection) (*UserDBRecord, error) {
	userrec := &UserDBRecord{
		Username: username,
		Subscriptions: make([]string, 0),
	}
	err := usercoll.Insert(userrec)
	return userrec, err
}

func createChannel(slug, title string, chancoll *mgo.Collection) (*ChannelDBRecord, error) {
	chanrec := &ChannelDBRecord{
		Slug: slug,
		Title: title,
		Items: make([]ItemDBRecord, 0),
	}
	err := chancoll.Insert(chanrec)
	return chanrec, err
}

func createItem(chanSlug, slug, title string, dateUploaded time.Time, b64data, uploader string, chancoll *mgo.Collection) (*ItemDBRecord, error) {
	var chanrec ChannelDBRecord
	err := chancoll.FindId(chanSlug).One(&chanrec)
	if err != nil {
		return nil, err
	}
	itemrec := &ItemDBRecord{
		Slug: slug,
		Title: title,
		DateUploaded: dateUploaded,
		Data: b64data,
		Uploader: uploader,
	}
	chanrec.Items = append(chanrec.Items, *itemrec)
	err = chancoll.UpdateId(chanSlug, chanrec)
	return itemrec, err
}

func (self *ApiSuite) loadTestData(c *C) {
	c.Log("Loading test data...")
	dataDir := os.Getenv("TEST_DATADIR")

	_, err := createUser(self.user1.Username, self.apiConfig.usercoll)
	c.Assert(err, IsNil)
	_, err = createUser(self.user2.Username, self.apiConfig.usercoll)
	c.Assert(err, IsNil)
	//Intentionally omitted; we don't want this user to exist.
	//createUser(self.baduser1.Username, self.apiConfig.usercoll)

	self.chan1Rec = &ChannelJSONRecord{
		Slug: "test-channel-1",
		Title: "Test Channel 1",
		Items: make([]*ItemJSONRecord, 0),
	}
	createChannel(
		self.chan1Rec.Slug,
		self.chan1Rec.Title,
		self.apiConfig.chancoll,
	)

	self.chan2Rec = &ChannelJSONRecord{
		Slug: "test-channel-2",
		Title: "Test Channel 2",
		Items: make([]*ItemJSONRecord, 0),
	}
	createChannel(
		self.chan2Rec.Slug,
		self.chan2Rec.Title,
		self.apiConfig.chancoll,
	)

	self.item1Rec = &ItemJSONRecord{
		Slug: "test-item-1",
		Title: "Test Item 1",
		DateUploaded: time.Now(),
		Uploader: self.user1.Username,
	}
	rawDataItem1, err := ioutil.ReadFile(filepath.Join(dataDir, "item1.jpg"))
	c.Assert(err, IsNil)
	b64DataItem1 := base64.StdEncoding.EncodeToString(rawDataItem1)
	createItem(self.chan1Rec.Slug, self.item1Rec.Slug, self.item1Rec.Title,
		self.item1Rec.DateUploaded, b64DataItem1, self.item1Rec.Uploader,
		self.apiConfig.chancoll)

	self.item2Rec = &ItemJSONRecord{
		Slug: "test-item-1",
		Title: "Test Item 1",
		DateUploaded: time.Now(),
		Uploader: self.user1.Username,
	}
	rawDataItem2, err := ioutil.ReadFile(filepath.Join(dataDir, "item2.jpg"))
	c.Assert(err, IsNil)
	b64DataItem2 := base64.StdEncoding.EncodeToString(rawDataItem2)
	createItem(self.chan1Rec.Slug, self.item2Rec.Slug, self.item2Rec.Title,
		self.item2Rec.DateUploaded, b64DataItem2, self.item2Rec.Uploader,
		self.apiConfig.chancoll)


	c.Log("Test data loaded!")
}

func (self *ApiSuite) TearDownSuite(c *C) {
	err := self.apiConfig.db.DropDatabase()
	c.Assert(err, IsNil)
	c.Log("Dropped testing database")
}

func (self *ApiSuite) SetUpSuite(c *C) {
	mongoHost := os.Getenv("MONGO_HOST")
	mongoHostLen := len(mongoHost)
	c.Assert(mongoHostLen > 0, Equals, true)
	dburl := "mongodb://" + mongoHost
	session, err := mgo.Dial(dburl)
	c.Assert(err, IsNil)

	self.apiConfig = NewConfig(session, "testing")
	self.loadTestData(c)
}

func (self *ApiSuite) badAuthDo(requester *testflight.Requester, verb, route, authHeader, authHeaderContents string) (*testflight.Response, error) {
	req, err := http.NewRequest(verb, route, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(authHeader, authHeaderContents)
	return requester.Do(req), nil
}

func (self *ApiSuite) authDo(requester *testflight.Requester, username, verb, route string, body []byte) (*testflight.Response, error) {
	var bodyReader io.Reader = nil
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest(verb, route, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer SUP3R_S33CR37:" + username)
	return requester.Do(req), nil
}

func (self *ApiSuite) authGet(requester *testflight.Requester, username, route string) (*testflight.Response, error) {
	return self.authDo(requester, username, "GET", route, nil)
}

func (self *ApiSuite) authPost(requester *testflight.Requester, username, route string, body []byte) (*testflight.Response, error) {
	return self.authDo(requester, username, "POST", route, body)
}

func (self *ApiSuite) TestGetChannelListBadUser(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.authGet(r, self.baduser1.Username, "/channel")
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusUnauthorized)
    })
}

func (self *ApiSuite) TestGetChannelListBadAuthHeaderName(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.badAuthDo(r, "GET", "/channel", "Auhorization", "Bearer SUP3R_S33CR37:" + self.user1.Username)
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusUnauthorized)
    })
}

func (self *ApiSuite) TestGetChannelListBadAuthType(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.badAuthDo(r, "GET", "/channel", "Authorization", "FNORD SUP3R_S33CR37:" + self.user1.Username)
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusBadRequest)
    })
}

func (self *ApiSuite) TestGetChannelListBadAuthSecret(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.badAuthDo(r, "GET", "/channel", "Authorization", "Bearer SUP3R_S3CR37:" + self.user1.Username)
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusBadRequest)
    })
}

func (self *ApiSuite) TestGetChannelListBadAuthNoSecret(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.badAuthDo(r, "GET", "/channel", "Authorization", "TROLOLOL:" + self.user1.Username)
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusBadRequest)
    })
}

func (self *ApiSuite) TestGetChannelListBadAuthNoUser(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.badAuthDo(r, "GET", "/channel", "Authorization", "Bearer SUP3R_S33CR37")
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusBadRequest)
    })
}

func (self *ApiSuite) TestGetChannelListGoodAuth(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.authGet(r, self.user1.Username, "/channel")
	    c.Assert(err, IsNil)
	    c.Log(response.Body)
	    c.Assert(response.StatusCode, Equals, http.StatusOK)
    })
}

func (self *ApiSuite) TestGetChannelInfoBadAuth(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.authGet(r, self.baduser1.Username, "/channel/" + self.chan1Rec.Slug)
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusUnauthorized)
    })
}

func (self *ApiSuite) TestGetChannelInfoGoodAuth(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.authGet(r, self.user1.Username, "/channel/" + self.chan1Rec.Slug)
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusOK)
    })
}

func (self *ApiSuite) TestGetChannelInfoBadChannel(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.authGet(r, self.user1.Username, "/channel/nosuchchannel")
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusNotFound)
    })
}


func (self *ApiSuite) TestGetItemListBadAuth(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.authGet(r, self.baduser1.Username, "/channel/" + self.chan1Rec.Slug + "/item")
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusUnauthorized)
    })
}

func (self *ApiSuite) TestGetItemListGoodAuth(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.authGet(r, self.user1.Username, "/channel/" + self.chan1Rec.Slug + "/item")
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusOK)
	    results := make(map[string]string)
	    err = json.Unmarshal(response.RawBody, &results)
	    c.Assert(err, IsNil)

	    title, ok := results[self.item1Rec.Slug]
	    c.Assert(ok, Equals, true)
	    c.Assert(title, Equals, self.item1Rec.Title)
	    title, ok = results[self.item2Rec.Slug]
	    c.Assert(ok, Equals, true)
	    c.Assert(title, Equals, self.item2Rec.Title)
    })
}

func (self *ApiSuite) TestGetItemListBadChannel(c *C) {
    testflight.WithServer(self.apiConfig.GetRouter(), func(r *testflight.Requester) {
	    response, err := self.authGet(r, self.user1.Username, "/channel/nosuchchannel/item")
	    c.Assert(err, IsNil)
	    c.Assert(response.StatusCode, Equals, http.StatusNotFound)
    })
}

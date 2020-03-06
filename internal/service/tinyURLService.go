package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"tinyURL/internal/repository"

	"github.com/gin-gonic/gin"
)

type tinyURL struct {
	URL   string `json:"url"`
	Error error  `json:"error,omitempty"`
	repo  repository.Repository
}

type tinyURLService struct {
	Service
}

func NewTinyURLService(s Service) *tinyURLService {
	return &tinyURLService{s}
}

func (t *tinyURL) Get(c *gin.Context) (string, error) {
	hashval := c.Param("hashval")
	record, err := t.repo.GetByPrimaryKey(hashval)
	if err != nil {
		return "", err
	}

	return record, err
}

func (t *tinyURL) Post(c *gin.Context) (string, error) {
	url := c.PostForm("url")
	URLHash := hashURL(url)
	m := map[string]string{"URL": url, "Hashval": URLHash}

	newRecord := t.repo.New(m)
	err := newRecord.InsertDB()
	//repo.InsertDB()
	if err != nil {
		return "", err
	}

	response, err := json.Marshal(*(NewtinyURL(URLHash, nil, nil)))
	if err != nil {
		return "", err
	}
	return string(response), err
}

func NewtinyURL(hashval string, r repository.Repository, err error) *tinyURL {
	newURL := domainName + "/" + hashval
	return &tinyURL{newURL, err, r}
}

func hashURL(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hashvalue := hex.EncodeToString(hasher.Sum(nil))
	return hashvalue[:8]
}

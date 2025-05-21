package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"divar_recommender/internal/models"
)

const (
	getPostRoute        = "%s/v1/open-platform/finder/post/%s"
	getPostsSearchRoute = "%s/v2/open-platform/finder/post"
)

type DivarService struct {
	token string
	uri   string
}

func NewDivarService(uri string, token string) *DivarService {
	return &DivarService{
		token: token,
		uri:   uri,
	}
}

func (s *DivarService) GetPost(postToken string) (models.Post, error) {
	requestUrl := fmt.Sprintf(getPostRoute, s.uri, postToken)
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return models.Post{}, err
	}

	req.Header.Add("X-Api-Key", s.token)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Post{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Post{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return models.Post{}, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var data models.Post
	err = json.Unmarshal(body, &data)
	if err != nil {
		return models.Post{}, fmt.Errorf("JSON unmarshal error: %v, body: %s", err, string(body))
	}

	return data, nil
}

func (s *DivarService) GetPosts(requestModel models.GetPostsRequestModel) ([]models.PostItem, error) {
	requestUrl := fmt.Sprintf(getPostsSearchRoute, s.uri)

	body, err := json.Marshal(requestModel)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", s.token)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data models.PostItems
	err = json.Unmarshal(result, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, err
}

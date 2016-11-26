// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http:#www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// Package backend exposes a REST API to manage posts stored in the Google
// Cloud Datastore using the Cloud Endpoints feature of App Engine.
package backend

import (
	"net/url"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
)

// PostsAPI defines all the endpoints of the posts API.
type PostsAPI struct{}

// A Post contains all the information related to a post.
type Post struct {
	UID      *datastore.Key `json:"uid" datastore:"-"`
	Text     string         `json:"text"`
	Username string         `json:"username"`
	Avatar   string         `json:"avatar"`
	Favorite bool           `json:"favorite"`
}

// Posts contains a slice of posts. This type is needed because go-endpoints
// only supports pointers to structs as input and output types.
type Posts struct {
	Posts []Post `json:"posts"`
}

// List returns a list of all the existing posts.
func (PostsAPI) List(c context.Context) (*Posts, error) {
	if err := checkReferer(c); err != nil {
		return nil, err
	}

	posts := []Post{}
	keys, err := datastore.NewQuery("Post").GetAll(c, &posts)
	if err != nil {
		return nil, err
	}
	for i, k := range keys {
		posts[i].UID = k
	}
	return &Posts{posts}, nil
}

// AddRequest contains all the fields needed to create a new Post.
type AddRequest struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// Add creates a new post given the fields in AddRequest, stores it in the
// datastore, and returns it.
func (PostsAPI) Add(c context.Context, r *AddRequest) (*Post, error) {
	if err := checkReferer(c); err != nil {
		return nil, err
	}

	k := datastore.NewIncompleteKey(c, "Post", nil)
	t := &Post{Text: r.Text, Username: r.Username, Avatar: r.Avatar}
	k, err := datastore.Put(c, k, t)
	if err != nil {
		return nil, err
	}
	t.UID = k
	return t, nil
}

// SetFavoriteRequest contains the information needed to change the favorite
// status of a post.
type SetFavoriteRequest struct {
	UID      *datastore.Key
	Favorite bool
}

// SetFavorite changes the favorite status of a post given its UID.
func (PostsAPI) SetFavorite(c context.Context, r *SetFavoriteRequest) error {
	if err := checkReferer(c); err != nil {
		return err
	}

	return datastore.RunInTransaction(c, func(c context.Context) error {
		var p Post
		if err := datastore.Get(c, r.UID, &p); err == datastore.ErrNoSuchEntity {
			return endpoints.NewNotFoundError("post not found")
		} else if err != nil {
			return err
		}
		p.Favorite = r.Favorite
		_, err := datastore.Put(c, r.UID, &p)
		return err
	}, nil)
}

// checkReferer returns an error if the referer of the HTTP request in the
// given context is not allowed.
//
// The allowed referer is the appspot domain for the application, such as:
//   my-project-id.appspot.com
// and all domains are accepted when running locally on dev app server.
func checkReferer(c context.Context) error {
	if appengine.IsDevAppServer() {
		return nil
	}

	r := endpoints.HTTPRequest(c).Referer()
	u, err := url.Parse(r)
	if err != nil {
		log.Infof(c, "malformed referer detected: %q", r)
		return endpoints.NewUnauthorizedError("couldn't extract domain from referer")
	}

	if u.Host != appengine.AppID(c)+".appspot.com" {
		log.Infof(c, "unauthorized referer detected: %q", r)
		return endpoints.NewUnauthorizedError("referer unauthorized")
	}
	return nil
}

func init() {
	// register the posts API with cloud endpoints.
	api, err := endpoints.RegisterService(PostsAPI{}, "posts", "v1", "posts api", true)
	if err != nil {
		panic(err)
	}

	// adapt the name, method, and path for each method.
	info := api.MethodByName("List").Info()
	info.Name, info.HTTPMethod, info.Path = "getPosts", "GET", "posts"

	info = api.MethodByName("SetFavorite").Info()
	info.Name, info.HTTPMethod, info.Path = "setFavorite", "PUT", "posts"

	info = api.MethodByName("Add").Info()
	info.Name, info.HTTPMethod, info.Path = "addPost", "POST", "posts"

	// start handling cloud endpoint requests.
	endpoints.HandleHTTP()
}

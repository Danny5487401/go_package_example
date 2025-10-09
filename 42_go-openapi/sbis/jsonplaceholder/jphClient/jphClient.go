package jphClient

import (
	"github.com/Danny5487401/go_package_example/42_go-openapi/sbis/jsonplaceholder/gen/client"
	"github.com/Danny5487401/go_package_example/42_go-openapi/sbis/jsonplaceholder/gen/client/operations"
	"github.com/Danny5487401/go_package_example/42_go-openapi/sbis/jsonplaceholder/gen/models"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

func New() *client.Jsonplaceholder {
	jsonPlaceHolderHost := "jsonplaceholder.typicode.com"
	transport := httptransport.New(jsonPlaceHolderHost, "/", nil)
	client := client.New(transport, strfmt.Default)
	return client
}

func GetPosts(client *client.Jsonplaceholder) (*operations.GetPostsOK, error) {
	params := operations.NewGetPostsParams()
	ok, err := client.Operations.GetPosts(params)
	if err != nil {
		return nil, err
	}
	return ok, nil
}

func PostPost(postObj *models.NewJSONPlaceholderPost, client *client.Jsonplaceholder) (*operations.PostPostCreated, error) {
	params := operations.NewPostPostParams()
	params.PostObject = postObj
	ok, err := client.Operations.PostPost(params)
	if err != nil {
		return nil, err
	}
	return ok, nil
}

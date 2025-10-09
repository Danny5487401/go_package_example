package handlers

import (
	"fmt"
	"github.com/Danny5487401/go_package_example/42_go-openapi/nbi/gen/server/models"
	"github.com/Danny5487401/go_package_example/42_go-openapi/nbi/gen/server/restapi/operations"
	jphModels "github.com/Danny5487401/go_package_example/42_go-openapi/sbis/jsonplaceholder/gen/models"
	"github.com/Danny5487401/go_package_example/42_go-openapi/sbis/jsonplaceholder/jphClient"
	"github.com/go-openapi/runtime/middleware"
	"os"
	"runtime"
)

//-------------------------------
// HostInfo service
//-------------------------------

func GetHostInfo(params operations.GetHostInfoParams) middleware.Responder {
	host, _ := os.Hostname()
	numCpu := runtime.NumCPU()

	// these two may only work if Go is installed on the host!
	// unless these values are captured at compile time
	arch := runtime.GOARCH
	rtOs := runtime.GOOS

	//fmt.Printf("Host name: %s; os: %s; arch: %s; num CPUs: %d\n", host, os, arch, numCpu)

	info := models.HostInfo{}
	info.HostName = host
	info.Architecture = arch
	info.OsName = rtOs
	info.NumCpus = int64(numCpu)

	return operations.NewGetHostInfoOK().WithPayload(&info)
}

//-------------------------------
// PhoneBook service
//-------------------------------

func GetPhoneBook(params operations.GetPhoneBookParams) middleware.Responder {

	list := PhoneBookDb.Entries()
	return operations.NewGetPhoneBookOK().WithPayload(list)
}

func AddPhoneBookEntry(params operations.AddPhoneBookEntryParams) middleware.Responder {

	entry := params.Entry

	PhoneBookDb.AddEntry(entry)

	return operations.NewAddPhoneBookEntryOK().WithPayload(entry)
}

func GetPhoneBookEntry(params operations.GetPhoneBookEntryParams) middleware.Responder {

	entry := PhoneBookDb.GetEntry(params.First, params.Last)

	if entry == nil {
		// Response for 404
		errMsg := fmt.Sprintf("No entry for %s-%s.", params.First, params.Last)
		return operations.NewGetPhoneBookEntryNotFound().WithPayload(errMsg)
	}

	// Response for success (200)
	return operations.NewGetPhoneBookEntryOK().WithPayload(entry)
}

//-------------------------------
// Json PlaceHolder service
//-------------------------------

func GetPostTitles(params operations.GetPostTitlesParams) middleware.Responder {

	client := jphClient.New()

	getResponse, err := jphClient.GetPosts(client)
	if err != nil {
		fmt.Println(err.Error())
		return operations.NewGetPostTitlesInternalServerError().WithPayload(err.Error())
	}

	// read SBI struct and create similar NBI struct from it

	titles := make([]*models.PostTitle, 0, len(getResponse.Payload))
	for _, aPost := range getResponse.Payload {

		title := &models.PostTitle{ID: aPost.ID, Title: aPost.Title}
		titles = append(titles, title)
	}

	return operations.NewGetPostTitlesOK().WithPayload(titles)
}

func GetPostsByUser(params operations.GetPostsByUserParams) middleware.Responder {

	client := jphClient.New()

	getResponse, err := jphClient.GetPosts(client)
	if err != nil {
		fmt.Println(err.Error())
		return operations.NewGetPostsByUserInternalServerError().WithPayload(err.Error())
	}

	// read SBI struct and create similar NBI struct from it

	// get capacity to cover worst case, we could also start small ans let the slice grow as needed
	titles := make([]*models.UserPost, 0, len(getResponse.Payload))
	for _, aPost := range getResponse.Payload {

		if aPost.UserID == params.User {
			title := &models.UserPost{ID: aPost.ID, Title: aPost.Title, Body: aPost.Body}
			titles = append(titles, title)
		}
	}

	return operations.NewGetPostsByUserOK().WithPayload(titles)
}

func AddPostObject(params operations.AddPostObjectParams) middleware.Responder {
	// in this function we convert a NBI object to its SBI equivalent
	// and a SBI object to its NBI equivalent.
	// In this demo the data going across the NBI and the SBI are very similar (while being defined in different packages).
	// In more realistic applications the data at each API could require more significant adaptation or be totally different.

	client := jphClient.New()

	// convert NBI object to POST to SBI equivalent object
	postObject := params.Post

	postObj := &jphModels.NewJSONPlaceholderPost{}
	postObj.UserID = postObject.UserID
	postObj.Title = postObject.Title
	postObj.Body = postObject.Body

	postResponse, err := jphClient.PostPost(postObj, client)
	if err != nil {
		return operations.NewAddPostObjectInternalServerError().WithPayload(err.Error())
	}

	// convert SBI response object into NBI object
	sbiPost := postResponse.Payload

	po := &models.PostObject{}
	po.Body = sbiPost.Body
	po.ID = sbiPost.ID
	po.Title = sbiPost.Title
	po.UserID = sbiPost.UserID

	return operations.NewAddPostObjectOK().WithPayload(po)
}

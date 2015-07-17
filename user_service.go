/*
Used to consume external authorization.
Given a key parameter, return a user's ability to download the project refs.
Uses https://github.com/bndr/gopencils
 */
package main
import (
	"fmt"
	"log"
	"encoding/json"
	"bytes"
	"errors"
)

// declare the type of UserAccessGetter
// The url must conform to consumers_spec user_service spec
//type UserAccessGetter func(url string) string

type UserAccessResponse struct {
	Access      bool `json:"access"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	RawResponse []byte
	Filled      bool
}

type UserService struct {
	Downloader                *Downloader
	Username, Project, Action string
	UserAccessResponse        *UserAccessResponse
}

var AllowedActions = []string{"download", "push", "force_push", "admin"}

func (us *UserService) vetAction() bool {
	for _, b := range AllowedActions {
		if b == us.Action {
			return true
		}
	}
	return false
}

func NewUserService(base string, username string, project string, action string) *UserService {
	url := fmt.Sprintf("%s?username=%s&project=%s&action=%s", base, username, project, action)
	us := &UserService{Downloader: NewDownloader(url), Username: username, Project: project, Action: action}
	us.UserAccessResponse = &UserAccessResponse{Filled: false}
	if us.vetAction() != true {
		log.Println(action, "is not in AllowedActions")
		us.UserAccessResponse.Message = fmt.Sprintf("%s is not in AllowedActions", us.Action)
	}
	return us
}

// fills UserAccessResponse.RawResponse and pushes json into struct
func (us *UserService) GetResponse() (error) {
	us.Downloader.GetPage()
	buf := bytes.NewBuffer(us.Downloader.Response)
	us.UserAccessResponse.RawResponse = buf.Bytes()
	rdr := bytes.NewReader(us.Downloader.Response)
	// read the json into our struct
	if err := json.NewDecoder(rdr).Decode(&us.UserAccessResponse); err == nil {
		us.UserAccessResponse.Filled = true
		return nil
	} else {
		errors.New(fmt.Sprintf("STATUS ERR:", err))
		return err
	}

}

func (us *UserService) Can() bool {
	if !us.UserAccessResponse.Filled {
		us.GetResponse()
	}
	return us.UserAccessResponse.Access
}
package req

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
)

type GithubClient struct {
	*req.Client
}

func NewGithubClient() *GithubClient {
	return &GithubClient{
		Client: req.C().
			SetBaseURL("https://api.github.com").
			SetCommonErrorResult(&ErrorMessage2{}).
			EnableDumpEachRequest().
			OnAfterResponse(func(client *req.Client, resp *req.Response) error {
				if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
					return nil
				}
				if errMsg, ok := resp.ErrorResult().(*ErrorMessage2); ok {
					resp.Err = errMsg // Convert api error into go error
					return nil
				}
				if !resp.IsSuccessState() {
					// Neither a success response nor a error response, record details to help troubleshooting
					resp.Err = fmt.Errorf(
						"bad status: %s\nraw content:\n%s",
						resp.Status,
						resp.Dump(),
					)
				}
				return nil
			}),
	}
}

type UserProfile struct {
	Name string `json:"name"`
	Blog string `json:"blog"`
}

// GetUserProfile_Style1 returns the user profile for the specified user.
// Github API doc: https://docs.github.com/en/rest/users/users#get-a-user
func (c *GithubClient) GetUserProfile_Style1(
	ctx context.Context,
	username string,
) (user *UserProfile, err error) {
	_, err = c.R().
		SetContext(ctx).
		SetPathParam("username", username).
		SetSuccessResult(&user).
		Get("/users/{username}")
	return
}

// GetUserProfile_Style2 returns the user profile for the specified user.
// Github API doc: https://docs.github.com/en/rest/users/users#get-a-user
func (c *GithubClient) GetUserProfile_Style2(
	ctx context.Context,
	username string,
) (user *UserProfile, err error) {
	err = c.Get("/users/{username}").
		SetPathParam("username", username).
		Do(ctx).
		Into(&user)
	return
}

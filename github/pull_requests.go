package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// TODO DRY these functions as they're all doing mostly the same thing

// TODO refactor to accept a struct instead of param list

// UpdatePullRequest updates an existing pull request
func UpdatePullRequest(repo string, number int, title, body, state, base string) (*PullRequest, error) {
	reqURL := fmt.Sprintf("%s/repos/%s/pulls/%d", PullRequestsURL, repo, number)

	ppr := PatchPullRequest{title, body, state, base}
	patch, err := json.Marshal(ppr)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	req, err := http.NewRequest("PATCH", reqURL, bytes.NewBuffer(patch))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	acceptHeaders := []string{"application/vnd.github.symmetra-preview+json", "application/vnd.github.sailor-v-preview+json"}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("OAUTH_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var msg *Response
		json.NewDecoder(res.Body).Decode(&msg)
		res.Body.Close()
		return nil, fmt.Errorf("could not update PR %d from %s: %s", number, repo, msg.Message)
	}

	var pr *PullRequest
	if err := json.NewDecoder(res.Body).Decode(&pr); err != nil {
		res.Body.Close()
		return nil, err
	}

	return pr, nil
}

// MergePullRequest will merge an existing pull request
func MergePullRequest(repo string, number int, title, message, method string) (*MergeStatus, error) {
	reqURL := fmt.Sprintf("%s/repos/%s/pulls/%d/merge?", PullRequestsURL, repo, number)

	if len(title) > 0 {
		reqURL = fmt.Sprintf("%s&commit_title=%s", reqURL, title)
	}

	if len(message) > 0 {
		reqURL = fmt.Sprintf("%s&commit_message=%s", reqURL, message)
	}

	if len(method) > 0 {
		reqURL = fmt.Sprintf("%s&merge_method=%s", reqURL, method)
	}

	req, err := http.NewRequest("PUT", reqURL, nil)
	if err != nil {
		return nil, err
	}

	acceptHeaders := []string{"application/vnd.github.symmetra-preview+json", "application/vnd.github.sailor-v-preview+json"}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("OAUTH_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var msg *Response
		json.NewDecoder(res.Body).Decode(&msg)
		res.Body.Close()

		return nil, fmt.Errorf("could not merge PR %d to %s: %s", number, repo, msg.Message)
	}

	var status *MergeStatus
	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		res.Body.Close()
		return nil, err
	}

	res.Body.Close()
	return status, nil
}

// CreatePullRequest opens a new pull request
func CreatePullRequest(repo, title, body, head, base string) (*PullRequest, error) {
	reqURL := fmt.Sprintf("%s/repos/%s/pulls", PullRequestsURL, repo)

	npr := NewPullRequest{title, body, head, base}
	prs, err := json.Marshal(npr)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(prs))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	acceptHeaders := []string{"application/vnd.github.symmetra-preview+json", "application/vnd.github.sailor-v-preview+json"}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("OAUTH_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		var error *Error
		json.NewDecoder(res.Body).Decode(&error)
		res.Body.Close()

		if len(error.Errors) > 0 {
			return nil, fmt.Errorf("Error creating PR %s", error.Errors[0].Message)
		}

		return nil, fmt.Errorf("Create PR failed: %s", res.Status)
	}

	var pr *PullRequest
	if err := json.NewDecoder(res.Body).Decode(&pr); err != nil {
		res.Body.Close()
		return nil, err
	}

	return pr, nil
}

// ListPullRequestStatuses returns the statuses from a pull request
func ListPullRequestStatuses(reqURL string) ([]*Status, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	acceptHeaders := []string{"application/vnd.github.symmetra-preview+json", "application/vnd.github.sailor-v-preview+json"}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("OAUTH_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var msg *Response
		json.NewDecoder(res.Body).Decode(&msg)
		res.Body.Close()
		return nil, fmt.Errorf("listing PR statuses failed: %s", msg.Message)
	}

	var statuses []*Status
	if err := json.NewDecoder(res.Body).Decode(&statuses); err != nil {
		res.Body.Close()
		return nil, err
	}

	res.Body.Close()
	return statuses, nil
}

// ListPullRequestComments returns comments dfor a particular pull request
func ListPullRequestComments(repo string, number int) ([]*Comment, error) {
	reqURL := fmt.Sprintf("%s/repos/%s/issues/%d/comments", PullRequestsURL, repo, number)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	acceptHeaders := []string{"application/vnd.github.symmetra-preview+json", "application/vnd.github.sailor-v-preview+json"}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("OAUTH_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var msg *Response
		json.NewDecoder(res.Body).Decode(&msg)
		res.Body.Close()
		return nil, fmt.Errorf("List PR %d commits from %s failed: %s", number, repo, msg.Message)
	}

	var comments []*Comment
	if err := json.NewDecoder(res.Body).Decode(&comments); err != nil {
		res.Body.Close()
		return nil, err
	}

	res.Body.Close()
	return comments, nil
}

// ListPullRequestCommits returns commits for a particular pull request
func ListPullRequestCommits(repo string, number int) ([]*Commit, error) {
	reqURL := fmt.Sprintf("%s/repos/%s/pulls/%d/commits", PullRequestsURL, repo, number)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	acceptHeaders := []string{"application/vnd.github.symmetra-preview+json", "application/vnd.github.sailor-v-preview+json"}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("OAUTH_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var msg *Response
		json.NewDecoder(res.Body).Decode(&msg)
		res.Body.Close()
		return nil, fmt.Errorf("could not list PR %d commits from %s: %s", number, repo, msg.Message)
	}

	var commits []*Commit
	if err := json.NewDecoder(res.Body).Decode(&commits); err != nil {
		res.Body.Close()
		return nil, err
	}

	res.Body.Close()
	return commits, nil
}

// GetPullRequest returns a single pull request given a repo and a number
func GetPullRequest(repo string, number int) (*PullRequest, error) {
	reqURL := fmt.Sprintf("%s/repos/%s/pulls/%d", PullRequestsURL, repo, number)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	acceptHeaders := []string{"application/vnd.github.symmetra-preview+json", "application/vnd.github.sailor-v-preview+json"}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("OAUTH_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var msg *Response
		json.NewDecoder(res.Body).Decode(&msg)
		res.Body.Close()
		return nil, fmt.Errorf("could not get PR %d from %s: %s", number, repo, msg.Message)
	}

	var result *PullRequest
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		res.Body.Close()
		return nil, err
	}

	res.Body.Close()
	return result, nil
}

// ListPullRequests returns all pull requests for a repo
// https://developer.github.com/v3/pulls/#list-pull-requests
func ListPullRequests(repo string, allIssues bool) (*PullRequestsResult, error) {
	state := "state=open"
	if allIssues {
		state = "state=all"
	}
	reqURL := fmt.Sprintf("%s/repos/%s/pulls?%s", PullRequestsURL, repo, state)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	acceptHeaders := []string{"application/vnd.github.symmetra-preview+json", "application/vnd.github.sailor-v-preview+json"}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("OAUTH_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var msg *Response
		json.NewDecoder(res.Body).Decode(&msg)
		res.Body.Close()
		return nil, fmt.Errorf("could not list PRs from %s: %s", repo, msg.Message)
	}

	var pulls []*PullRequest
	if err := json.NewDecoder(res.Body).Decode(&pulls); err != nil {
		res.Body.Close()
		return nil, err
	}

	result := new(PullRequestsResult)
	result.PullRequests = pulls

	res.Body.Close()
	return result, nil
}

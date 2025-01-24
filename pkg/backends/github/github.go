package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"issues_sdk/pkg/issuetracker"
	"net/http"
)

type GitHubClient struct{
	Owner string
	Token string
	Repo string
}

func NewGitHubClient(owner, token, repo string) *GitHubClient{
	return &GitHubClient{Owner:owner,Token: token,Repo: repo}
}

func (g *GitHubClient) CreateIssue(issue issuetracker.Issue)(string, error){
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", g.Owner, g.Repo)
	payload := map[string]interface{}{
		"title": issue.Title,
		"body": issue.Description,
		"labels": issue.Labels,
		"assignees": issue.Assignees,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal issue payload: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "token "+ g.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	
	if resp.StatusCode != http.StatusCreated{
		return "", fmt.Errorf("GitHub API error : %s", res["message"])
	}

	return res["html_url"].(string), nil

}

func (g *GitHubClient) GetIssue(issueID string) (*issuetracker.Issue, error){
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", g.Owner, g.Repo, issueID)
	
	req, err := http.NewRequest("GET",url,nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "token "+ g.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 410{
		return nil, fmt.Errorf("No issue exists with the number:%d", issueID)
	}else if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("Failed to fetch issue: status code %d", resp.StatusCode)
	}

	body, err:= io.ReadAll(resp.Body)
	if err != nil{
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var issue issuetracker.Issue	
	if err := json.Unmarshal(body, &issue); err != nil{
		return nil, fmt.Errorf("error decoding json response: %v", err)
	}
	return &issue, nil

}
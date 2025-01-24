//to instantiate the appropriate backend
package issuetracker

import(
	"fmt"
	"issues_sdk/pkg/backends/github"
	"issues_sdk/pkg/backends/jira"
)

func NewIssueTracker(trackerType string, config map[string]string)(IssueTracker, error){
	switch trackerType {
	case "github":
		return github.NewGitHubClient(config["owner"],config["token"],config["repo"]),nil
	case "jira":
		return jira.NewJiraClient(config["baseURL"],config["username"],config["token"]),nil
	default:
		return nil, fmt.Errorf("Unsupported tracker type: %s", trackerType)
	}
}
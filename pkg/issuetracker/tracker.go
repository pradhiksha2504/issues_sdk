package issuetracker

type IssueTracker interface{
	CreateIssue(issue Issue)(string, error)
	GetIssue(issueID string)(*Issue, error)
	CloseIssue(issueID string)error
}
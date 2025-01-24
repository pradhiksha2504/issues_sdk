package issuetracker

type Issue struct{
	Title string `json:"title"`
	Description string `json:"desc"`
	Labels []string `json:"labels"`
	Assignees []string `json:"assignees"`
	Metadata map[string]string `json:"metadata"`
}
package model

type CustomAction struct {
	ProjectID		int		`json:"project_id"`
	MergeRequestID	int		`json:"iid"`
	Clicked			bool	`json:"clicked"`
}

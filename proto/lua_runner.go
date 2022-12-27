package proto

type ScriptRunnerRequest struct {
	Id      int32  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

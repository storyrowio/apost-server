package models

type AiChoice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Role string `json:"role"`
}

type AiUsage struct {
	CompletionToken int `json:"completionToken"`
	PromptToken     int `json:"promptToken"`
	TotalToken      int `json:"totalToken"`
}

type AiResult struct {
	Choices []AiChoice `json:"choices"`
	Usage   AiUsage    `json:"usage"`
}

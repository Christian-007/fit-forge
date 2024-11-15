package domains

type EmailAddressOptions struct {
	Email string `json:"email"`
	Name string `json:"name"`
}

type EmailWithTemplateRequest struct {
	From EmailAddressOptions `json:"from"`
	To []EmailAddressOptions `json:"to"`
	TemplateUuid string `json:"template_uuid"`
	TemplateVariables map[string]any `json:"template_variables"`
}
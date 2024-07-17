package main

// GenericReport represents a generic report structure
type Report struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Data        interface{} `json:"data"`
}

// NewReport creates a new report
func NewReport(id, name, description, reportType string, data interface{}) *Report {
	return &Report{
		ID:          id,
		Name:        name,
		Description: description,
		Type:        reportType,
		Data:        data,
	}
}

// GetID returns the report ID
func (r *Report) GetID() string {
	return r.ID
}

// GetName returns the report name
func (r *Report) GetName() string {
	return r.Name
}

//UpdateReport updates the report data
func (r *Report) UpdateReport(data interface{}) {
	r.Data = data
}

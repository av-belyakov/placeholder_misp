package datamodels

func NewEventReports() *EventReports {
	return &EventReports{}
}

// SetEventReportsName устанавливает значение для Name
func (report *EventReports) SetEventReportsName(v string) {
	report.Name = v
}

// GetEventReportsName возвращает значение Name
func (report *EventReports) GetEventReportsName() string {
	return report.Name
}

// SetEventReportsContent устанавливает значение для Content
func (report *EventReports) SetEventReportsContent(v string) {
	report.Content = v
}

// GetEventReportsContent возвращает значение Content
func (report *EventReports) GetEventReportsContent() string {
	return report.Content
}

// SetEventReportsDistribution устанавливает значение для Distribution
func (report *EventReports) SetEventReportsDistribution(v string) {
	report.Distribution = v
}

// GetEventReportsDistribution возвращает значение Distribution
func (report *EventReports) GetEventReportsDistribution() string {
	return report.Distribution
}

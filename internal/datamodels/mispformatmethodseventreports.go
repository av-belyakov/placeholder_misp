package datamodels

func NewEventReports() *EventReports {
	return &EventReports{}
}

// SetName устанавливает значение для Name
func (report *EventReports) SetName(v string) {
	report.Name = v
}

// GetName возвращает значение Name
func (report *EventReports) GetName() string {
	return report.Name
}

// SetContent устанавливает значение для Content
func (report *EventReports) SetContent(v string) {
	report.Content = v
}

// GetContent возвращает значение Content
func (report *EventReports) GetContent() string {
	return report.Content
}

// SetDistribution устанавливает значение для Distribution
func (report *EventReports) SetDistribution(v string) {
	report.Distribution = v
}

// GetDistribution возвращает значение Distribution
func (report *EventReports) GetDistribution() string {
	return report.Distribution
}

// Comparison выполняет сравнение двух объектов типа EventReports
func (report *EventReports) Comparison(newReport *EventReports) bool {
	if report.Name != newReport.Name {
		return false
	}

	if report.Content != newReport.Content {
		return false
	}

	if report.Distribution != newReport.Distribution {
		return false
	}

	return true
}

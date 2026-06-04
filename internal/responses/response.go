package responses

// ******* для ответов на запрос о местоположении и принадлежности сенсоров *******
// ResponseSensorsInformation ответ от внешнего сервиса на запрос о местоположении
// и принадлежности сенсоров
type ResponseSensorsInformation struct {
	Informations []SensorInformation `json:"found_information"`
	Source       string              `json:"source"`
	TaskId       string              `json:"task_id"`
	Error        string              `json:"error"`
}

// SensorInformation подробная информация о местоположении и принадлежности сенсора
type SensorInformation struct {
	INN                      string `json:"inn"`
	GeoCode                  string `json:"geo_code"`
	HomeNet                  string `json:"home_net"`
	SensorID                 string `json:"sensor_id"`
	ObjectArea               string `json:"object_area"`
	SpecialSensorID          string `json:"special_sensor_id"`
	OrganizationName         string `json:"organization_name"`
	FullOrganizationName     string `json:"full_organization_name"`
	SubjectRussianFederation string `json:"subject_russian_federation"`
	NetboxTenantGroup        string `json:"netbox_tenant_group"`
	Error                    string `json:"error"`
}

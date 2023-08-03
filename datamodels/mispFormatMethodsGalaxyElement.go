package datamodels

func (gmisp GalaxyElementMispFormat) GetGalaxyElementMisp() GalaxyElementMispFormat {
	return gmisp
}

func (gmisp *GalaxyElementMispFormat) SetValueIdGalaxyElementMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gmisp.Id = data

		isSuccess = true
	}

	return isSuccess
}

func (gmisp *GalaxyElementMispFormat) SetValueGalaxyClusterIdGalaxyElementMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gmisp.GalaxyClusterId = data

		isSuccess = true
	}

	return isSuccess
}

func (gmisp *GalaxyElementMispFormat) SetValueKeyGalaxyElementMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gmisp.Key = data

		isSuccess = true
	}

	return isSuccess
}

func (gmisp *GalaxyElementMispFormat) SetValueValueGalaxyElementMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gmisp.Value = data

		isSuccess = true
	}

	return isSuccess
}

package models

type ExactTimeMarkerEvent struct {
	Type        string `json:"type"`
	Time        string `json:"time,omitempty"`
	IsStretched bool   `json:"isStretched"`
}

type MacroEvent struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Time string `json:"time"`
}

type ExecuteCommandEvent struct {
	UUID string `json:"uuid"`
}

type EffectiveTransitions struct {
	Runtime  string `json:"runtime"`
	Duration string `json:"duration"`
}

type AssetCutEvent struct {
	EffectiveTransitions EffectiveTransitions `json:"effectiveTransitions"`
}

type AssetEvent struct {
	UUID                 string               `json:"uuid"`
	Type                 string               `json:"type"`
	CustomAssetTypeName  string               `json:"customAssetTypeName"`
	ExternalID           string               `json:"externalId"`
	EffectiveTransitions EffectiveTransitions `json:"effectiveTransitions"`
	IsStretched          bool                 `json:"isStretched"`
}

type RotationEvent struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	ExternalID string `json:"externalId"`
}

type CommentEvent struct {
	Text string `json:"text"`
}

type SpotBlockEvent struct {
	FillLength         string                `json:"fillLength"`
	LogEventCollection []*LogEventCollection `json:"logEventCollection"`
}

type LogEventCollection struct {
	UUID                 string                `json:"uuid"`
	Type                 string                `json:"type"`
	VerifyID             string                `json:"verifyID"`
	DisplayText          string                `json:"displayText"`
	Chain                string                `json:"chain"`
	StatusCode           string                `json:"statusCode"`
	EditCode             string                `json:"editCode"`
	AssetCutEvent        *AssetCutEvent        `json:"assetCutEvent,omitempty"`
	ExactTimeMarkerEvent *ExactTimeMarkerEvent `json:"exactTimeMarkerEvent,omitempty"`
	MacroEvent           *MacroEvent           `json:"macroEvent,omitempty"`
	ExecuteCommandEvent  *ExecuteCommandEvent  `json:"executeCommandEvent,omitempty"`
	AssetEvent           *AssetEvent           `json:"assetEvent,omitempty"`
	SpotBlockEvent       *SpotBlockEvent       `json:"spotBlockEvent,omitempty"`
	RotationEvent        *RotationEvent        `json:"rotationEvent,omitempty"`
	CommentEvent         *CommentEvent         `json:"commentEvent,omitempty"`
}

type HourGroupCollection struct {
	Hour               int                   `json:"hour"`
	LogEventCollection []*LogEventCollection `json:"logEventCollection"`
}

type LogDataObject struct {
	StationUUID         string                 `json:"stationUUID"`
	Date                string                 `json:"date"`
	HourGroupCollection []*HourGroupCollection `json:"hourGroupCollection"`
	UUID                string                 `json:"uuid"`
}

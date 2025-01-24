package godm

import (
	"time"
)

type OPM struct {
	Header OPMHeader
	// MetaData OPMMetaData
	// Data     OPMData
}

type OPMHeader struct {
	CcsdsOpmVers   string    `odm:"CCSDS_OPM_VERS,required"`
	Comments       []string  `odm:"COMMENT"`
	Classification string    `odm:"CLASSIFICATION"`
	CreationDate   time.Time `odm:"CREATION_DATE,required"`
	Originator     string    `odm:"ORIGINATOR,required"`
	MessageId      string    `odm:"MESSAGE_ID"`
}

type OPMMetaData struct {
	// Optional
	Comments []string `json:"COMMENT" xml:"COMMENT"`

	// Mandatory
	ObjectName string `json:"OBJECT_NAME" xml:"OBJECT_NAME"`

	// Mandatory
	ObjectId string `json:"OBJECT_ID" xml:"OBJECT_ID"`

	// Mandatory
	CenterName string `json:"CENTER_NAME" xml:"CENTER_NAME"`

	// Mandatory
	RefFrame string `json:"REF_FRAME" xml:"REF_FRAME"`

	// Conditional, if it is not intrinsic to the reference frame
	RefFrameEpoch time.Time `json:"REF_FRAME_EPOCH" xml:"REF_FRAME_EPOCH"`

	// Mandatory
	TimeSystem string `json:"TIME_SYSTEM" xml:"TIME_SYSTEM"`
}

type OPMData struct {
	// Mandatory
	StateVector StateVector

	// Optional, none or all
	OsculatingKeplerianElements OsculatingKeplerianElements

	// Optional, mass required if maneuver specified
	SpacecraftParameters SpacecraftParameters

	// Optional
	CovarianceMatrix CovarianceMatrix

	// Optional, repeats for each maneuver
	ManeuverParametersList []ManeuverParameters

	// Optional, defined in an ICD, key must start with "USER_DEFINED_"
	UserDefinedParameters UserDefinedParameters
}

type StateVector struct {
	// Optional
	Comments []string `json:"COMMENT" xml:"COMMENT"`

	// Mandatory
	Epoch time.Time `json:"EPOCH" xml:"EPOCH"`

	// Mandatory, km
	X float64 `json:"X" xml:"X"`

	// Mandatory, km
	Y float64 `json:"Y" xml:"Y"`

	// Mandatory, km
	Z float64 `json:"Z" xml:"Z"`

	// Mandatory, km/s
	XDOT float64 `json:"X_DOT" xml:"X_DOT"`

	// Mandatory, km/s
	YDOT float64 `json:"Y_DOT" xml:"Y_DOT"`

	// Mandatory, km/s
	ZDOT float64 `json:"Z_DOT" xml:"Z_DOT"`
}

type OsculatingKeplerianElements struct {
	// Optional
	Comments []string `json:"COMMENT" xml:"COMMENT"`

	// Conditional, none or all, km
	SemiMajorAxis float64 `json:"SEMI_MAJOR_AXIS" xml:"SEMI_MAJOR_AXIS"`

	// Conditional, none or all
	Eccentricity float64 `json:"ECCENTRICITY" xml:"ECCENTRICITY"`

	// Conditional, none or all, deg
	Inclination float64 `json:"INCLINATION" xml:"INCLINATION"`

	// Conditional, none or all, deg
	RaOfAscNode float64 `json:"RA_OF_ASC_NODE" xml:"RA_OF_ASC_NODE"`

	// Conditional, none or all, deg
	ArgOfPericenter float64 `json:"ARG_OF_PERICENTER" xml:"ARG_OF_PERICENTER"`

	// Conditional, none or all, deg. Either this or MEAN_ANOMALY must be provided.
	TrueAnomaly float64 `json:"TRUE_ANOMALY" xml:"TRUE_ANOMALY"`
	MeanAnomaly float64 `json:"MEAN_ANOMALY" xml:"MEAN_ANOMALY"`

	// Conditional, none or all, km**3/s**2
	GM float64 `json:"GM" xml:"GM"`
}

type SpacecraftParameters struct {
	// Optional
	Comments []string `json:"COMMENT" xml:"COMMENT"`

	// Conditional, required if maneuver specified, kg
	Mass float64 `json:"MASS" xml:"MASS"`

	// Optional, m**2
	SolarRadArea float64 `json:"SOLAR_RAD_AREA" xml:"SOLAR_RAD_AREA"`

	// Optional
	SolarRadCoeff float64 `json:"SOLAR_RAD_COEFF" xml:"SOLAR_RAD_COEFF"`

	// Optional, m**2
	DragArea float64 `json:"DRAG_AREA" xml:"DRAG_AREA"`

	// Optional
	DragCoeff float64 `json:"DRAG_COEFF" xml:"DRAG_COEFF"`
}

type CovarianceMatrix struct {
	// Optional
	Comments []string `json:"COMMENT" xml:"COMMENT"`

	// Conditional, may be omitted if same as REF_FRAME
	CovRefFrame string `json:"COV_REF_FRAME" xml:"COV_REF_FRAME"`

	// Conditional, none or all, km**2
	CXX float64 `json:"CX_X" xml:"CX_X"`

	// Conditional, none or all, km**2
	CYX float64 `json:"CY_X" xml:"CY_X"`

	// Conditional, none or all, km**2
	CYY float64 `json:"CY_Y" xml:"CY_Y"`

	// Conditional, none or all, km**2
	CZX float64 `json:"CZ_X" xml:"CZ_X"`

	// Conditional, none or all, km**2
	CZY float64 `json:"CZ_Y" xml:"CZ_Y"`

	// Conditional, none or all, km**2
	CZZ float64 `json:"CZ_Z" xml:"CZ_Z"`

	// Conditional, none or all, km**2/s
	CXDOTX float64 `json:"CX_DOT_X" xml:"CX_DOT_X"`

	// Conditional, none or all, km**2/s
	CXDOTY float64 `json:"CX_DOT_Y" xml:"CX_DOT_Y"`

	// Conditional, none or all, km**2/s
	CXDOTZ float64 `json:"CX_DOT_Z" xml:"CX_DOT_Z"`

	// Conditional, none or all, km**2/s**2
	CXDOTXDOT float64 `json:"CX_DOT_X_DOT" xml:"CX_DOT_X_DOT"`

	// Conditional, none or all, km**2/s
	CYDOTX float64 `json:"CY_DOT_X" xml:"CY_DOT_X"`

	// Conditional, none or all, km**2/s
	CYDOTY float64 `json:"CY_DOT_Y" xml:"CY_DOT_Y"`

	// Conditional, none or all, km**2/s
	CYDOTZ float64 `json:"CY_DOT_Z" xml:"CY_DOT_Z"`

	// Conditional, none or all, km**2/s**2
	CYDOTXDOT float64 `json:"CY_DOT_X_DOT" xml:"CY_DOT_X_DOT"`

	// Conditional, none or all, km**2/s**2
	CYDOTYDOT float64 `json:"CY_DOT_Y_DOT" xml:"CY_DOT_Y_DOT"`

	// Conditional, none or all, km**2/s
	CZDOTX float64 `json:"CZ_DOT_X" xml:"CZ_DOT_X"`

	// Conditional, none or all, km**2/s
	CZDOTY float64 `json:"CZ_DOT_Y" xml:"CZ_DOT_Y"`

	// Conditional, none or all, km**2/s
	CZDOTZ float64 `json:"CZ_DOT_Z" xml:"CZ_DOT_Z"`

	// Conditional, none or all, km**2/s**2
	CZDOTXDOT float64 `json:"CZ_DOT_X_DOT" xml:"CZ_DOT_X_DOT"`

	// Conditional, none or all, km**2/s**2
	CZDOTYDOT float64 `json:"CZ_DOT_Y_DOT" xml:"CZ_DOT_Y_DOT"`

	// Conditional, none or all, km**2/s**2
	CZDOTZDOT float64 `json:"CZ_DOT_Z_DOT" xml:"CZ_DOT_Z_DOT"`
}

type ManeuverParameters struct {
	// Optional
	Comments []string `json:"COMMENT" xml:"COMMENT"`

	// Optional
	ManEpochIgnition time.Time `json:"MAN_EPOCH_IGNITION" xml:"MAN_EPOCH_IGNITION"`

	// Optional, s
	ManDuration float64 `json:"MAN_DURATION" xml:"MAN_DURATION"`

	// Optional, <0, kg
	ManDeltaMass float64 `json:"MAN_DELTA_MASS" xml:"MAN_DELTA_MASS"`

	// Optional
	ManRefFrame string `json:"MAN_REF_FRAME" xml:"MAN_REF_FRAME"`

	// Optional, km/s
	ManDV1 float64 `json:"MAN_DV_1" xml:"MAN_DV_1"`

	// Optional, km/s
	ManDV2 float64 `json:"MAN_DV_2" xml:"MAN_DV_2"`

	// Optional, km/s
	ManDV3 float64 `json:"MAN_DV_3" xml:"MAN_DV_3"`
}

type UserDefinedParameters map[string]string

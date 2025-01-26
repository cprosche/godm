package godm

import (
	"strings"
	"time"
)

// TODO: json & xml support
func ParseOPM(s string) (OPM, error) {
	result := OPM{}
	kvs, err := parseIntoKVs(s)
	if err != nil {
		return OPM{}, err
	}

	// handle user defined parameters
	for _, kv := range kvs {
		if strings.HasPrefix(kv.Key, "USER_DEFINED_") {
			if result.Data.UserDefinedParameters == nil {
				result.Data.UserDefinedParameters = map[string]string{}
			}
			result.Data.UserDefinedParameters[kv.Key] = kv.Value
		}
	}

	fields, err := getODMFields(&result)
	if err != nil {
		return OPM{}, err
	}

	err = populateODMFields(fields, kvs)
	if err != nil {
		return OPM{}, err
	}

	// correct comment placement
	if len(result.Data.CovarianceMatrix.Comments) > 0 && result.Data.CovarianceMatrix.CovRefFrame == "" && len(result.Data.ManeuverParametersList) > 0 {
		result.Data.ManeuverParametersList[0].Comments = append(result.Data.ManeuverParametersList[0].Comments, result.Data.CovarianceMatrix.Comments...)
		result.Data.CovarianceMatrix.Comments = nil
	}

	return result, nil
}

// TODO: handle raw
type OPM struct {
	Header   OPMHeader
	MetaData OPMMetaData
	Data     OPMData

	Raw string
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
	Comments      []string  `odm:"COMMENT"`
	ObjectName    string    `odm:"OBJECT_NAME,required"`
	ObjectId      string    `odm:"OBJECT_ID,required"`
	CenterName    string    `odm:"CENTER_NAME,required"`
	RefFrame      string    `odm:"REF_FRAME,required"`
	RefFrameEpoch time.Time `odm:"REF_FRAME_EPOCH"` // TODO: add validation for this Conditional, if it is not intrinsic to the reference frame
	TimeSystem    string    `odm:"TIME_SYSTEM,required"`
}

type OPMData struct {
	StateVector                 StateVector                 // Mandatory
	OsculatingKeplerianElements OsculatingKeplerianElements // TODO: add validation for this
	SpacecraftParameters        SpacecraftParameters        // TODO: add validation for this
	CovarianceMatrix            CovarianceMatrix            // TODO: add validation for this
	ManeuverParametersList      []ManeuverParameters        // TODO: add validation for this
	UserDefinedParameters       map[string]string
}

type StateVector struct {
	Comments []string  `odm:"COMMENT"`
	Epoch    time.Time `odm:"EPOCH,required"`
	X        float64   `odm:"X,required"`
	Y        float64   `odm:"Y,required"`
	Z        float64   `odm:"Z,required"`
	XDOT     float64   `odm:"X_DOT,required"`
	YDOT     float64   `odm:"Y_DOT,required"`
	ZDOT     float64   `odm:"Z_DOT,required"`
}

type OsculatingKeplerianElements struct {
	Comments        []string `odm:"COMMENT"`
	SemiMajorAxis   float64  `odm:"SEMI_MAJOR_AXIS"`
	Eccentricity    float64  `odm:"ECCENTRICITY"`
	Inclination     float64  `odm:"INCLINATION"`
	RaOfAscNode     float64  `odm:"RA_OF_ASC_NODE"`
	ArgOfPericenter float64  `odm:"ARG_OF_PERICENTER"`
	TrueAnomaly     float64  `odm:"TRUE_ANOMALY"`
	MeanAnomaly     float64  `odm:"MEAN_ANOMALY"`
	GM              float64  `odm:"GM"`
}

type ManeuverParameters struct {
	Comments         []string  `odm:"COMMENT"`
	ManEpochIgnition time.Time `odm:"MAN_EPOCH_IGNITION"`
	ManDuration      float64   `odm:"MAN_DURATION"`
	ManDeltaMass     float64   `odm:"MAN_DELTA_MASS"`
	ManRefFrame      string    `odm:"MAN_REF_FRAME"`
	ManDV1           float64   `odm:"MAN_DV_1"`
	ManDV2           float64   `odm:"MAN_DV_2"`
	ManDV3           float64   `odm:"MAN_DV_3"`
}

package godm

import (
	"strings"
	"time"
)

func ParseOMM(s string) (OMM, error) {
	result := OMM{}
	kvs, err := parseIntoKVs(s)
	if err != nil {
		return OMM{}, err
	}

	for _, kv := range kvs {
		println(kv.Key, kv.Value)
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
		return OMM{}, err
	}

	for _, f := range fields {
		println(f.Name, f.ReflectVal.Type().Name())
	}

	err = populateODMFields(fields, kvs)
	if err != nil {
		return OMM{}, err
	}

	return result, nil
}

// TODO: handle raw
type OMM struct {
	Header   OMMHeader
	MetaData OMMMetaData
	Data     OMMData

	Raw string
}

type OMMHeader struct {
	CcsdsOmmVers   string    `odm:"CCSDS_OMM_VERS,required"`
	Comments       []string  `odm:"COMMENT"`
	Classification string    `odm:"CLASSIFICATION"`
	CreationDate   time.Time `odm:"CREATION_DATE,required"`
	Originator     string    `odm:"ORIGINATOR,required"`
	MessageId      string    `odm:"MESSAGE_ID"`
}

type OMMMetaData struct {
	Comments          []string  `odm:"COMMENT"`
	ObjectName        string    `odm:"OBJECT_NAME,required"`
	ObjectId          string    `odm:"OBJECT_ID,required"`
	CenterName        string    `odm:"CENTER_NAME,required"`
	RefFrame          string    `odm:"REF_FRAME,required"`
	RefFrameEpoch     time.Time `odm:"REF_FRAME_EPOCH"` // TODO: add validation
	TimeSystem        string    `odm:"TIME_SYSTEM,required"`
	MeanElementTheory string    `odm:"MEAN_ELEMENT_THEORY"`
}

// TODO: fill this
type OMMData struct {
	MeanKeplerianElements MeanKeplerianElements
	SpacecraftParameters  SpacecraftParameters
	TLERelatedParameters  TLERelatedParameters
	CovarianceMatrix      CovarianceMatrix
	UserDefinedParameters map[string]string
}

type MeanKeplerianElements struct {
	Comments        []string  `odm:"COMMENT"`
	Epoch           time.Time `odm:"EPOCH,required"`
	SemiMajorAxis   float64   `odm:"SEMI_MAJOR_AXIS"` // TODO: add validation
	MeanMotion      float64   `odm:"MEAN_MOTION"`
	Eccentricity    float64   `odm:"ECCENTRICITY,required"`
	Inclination     float64   `odm:"INCLINATION,required"`
	RaOfAscNode     float64   `odm:"RA_OF_ASC_NODE,required"`
	ArgOfPericenter float64   `odm:"ARG_OF_PERICENTER,required"`
	MeanAnomaly     float64   `odm:"MEAN_ANOMALY,required"`
	Gm              float64   `odm:"GM"`
}

// TODO: some validation for this
type TLERelatedParameters struct {
	Comments           []string `odm:"COMMENT"`
	EphemerisType      string   `odm:"EPHEMERIS_TYPE"`
	ClassificationType string   `odm:"CLASSIFICATION_TYPE"`
	NoradCatId         int      `odm:"NORAD_CAT_ID"`
	ElementSetNo       int      `odm:"ELEMENT_SET_NO"`
	RevAtEpoch         int      `odm:"REV_AT_EPOCH"`
	BStar              float64  `odm:"BSTAR"`
	BTerm              float64  `odm:"BTERM"`
	MeanMotionDot      float64  `odm:"MEAN_MOTION_DOT"`
	MeanMotionDdot     float64  `odm:"MEAN_MOTION_DDOT"`
	Agom               float64  `odm:"AGOM"`
}

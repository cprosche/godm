package godm

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func ParseOPM(s string) (*OPM, error) {
	result := OPM{}
	result.Raw = strings.TrimSpace(s)
	lineEnding := detectLineEnding(result.Raw)
	lines := strings.Split(result.Raw, lineEnding)

	// remove empty lines
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			lines = append(lines[:i], lines[i+1:]...)
			i--
		}
	}

	// Parse header, Table 3-1
	foundMessageId := false
	// see if message contains MESSAGE_ID
	if strings.Contains(result.Raw, "MESSAGE_ID") {
		foundMessageId = true
	}
	header := OPMHeader{}
	i := 0
	for {
		line := lines[i]
		i++
		if line == "" {
			continue
		}

		k, v, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		switch k {
		case "CCSDS_OPM_VERS":
			header.CcsdsOpmVers = v
		case "COMMENT":
			header.Comments = append(header.Comments, v)
		case "CLASSIFICATION":
			header.Classification = v
		case "CREATION_DATE":
			if !strings.Contains(v, "Z") {
				v += "Z"
			}
			header.CreationDate, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
		case "ORIGINATOR":
			header.Originator = v
		case "MESSAGE_ID":
			header.MessageId = v
		}

		if foundMessageId && header.MessageId != "" {
			break
		}

		if !foundMessageId && header.Originator != "" {
			break
		}
	}

	result.Header = header

	// Parse metadata, Table 3-2
	metaData := OPMMetaData{}
	for {
		line := lines[i]
		i++

		if line == "" {
			continue
		}

		k, v, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		switch k {
		case "COMMENT":
			metaData.Comments = append(metaData.Comments, v)
		case "OBJECT_NAME":
			metaData.ObjectName = v
		case "OBJECT_ID":
			metaData.ObjectId = v
		case "CENTER_NAME":
			metaData.CenterName = v
		case "REF_FRAME":
			metaData.RefFrame = v
		case "REF_FRAME_EPOCH":
			if !strings.Contains(v, "Z") {
				v += "Z"
			}
			metaData.RefFrameEpoch, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
		case "TIME_SYSTEM":
			metaData.TimeSystem = v
		}

		if metaData.TimeSystem != "" {
			break
		}
	}

	result.MetaData = metaData

	// Parse data, Table 3-3
	// State Vector
	stateVector := StateVector{}
	for {
		line := lines[i]
		i++

		if line == "" {
			continue
		}

		k, v, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		foundZdot := false
		switch k {
		case "COMMENT":
			stateVector.Comments = append(stateVector.Comments, v)
		case "EPOCH":
			if !strings.Contains(v, "Z") {
				v += "Z"
			}
			stateVector.Epoch, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
		case "X":
			stateVector.X, err = parseFloat(v)
			if err != nil {
				return nil, err
			}
		case "Y":
			stateVector.Y, err = parseFloat(v)
			if err != nil {
				return nil, err
			}
		case "Z":
			stateVector.Z, err = parseFloat(v)
			if err != nil {
				return nil, err
			}
		case "X_DOT":
			stateVector.XDOT, err = parseFloat(v)
			if err != nil {

				return nil, err
			}
		case "Y_DOT":
			stateVector.YDOT, err = parseFloat(v)
			if err != nil {
				return nil, err
			}
		case "Z_DOT":
			stateVector.ZDOT, err = parseFloat(v)
			if err != nil {
				return nil, err
			}
			foundZdot = true
		}

		if foundZdot {
			break
		}
	}
	result.Data.StateVector = stateVector

	// Osculating Keplerian Elements
	if strings.Contains(result.Raw, "SEMI_MAJOR_AXIS") {
		kep := OsculatingKeplerianElements{}
		foundGm := false
		for {
			line := lines[i]
			i++

			if line == "" {
				continue
			}

			k, v, err := parseLine(line)
			if err != nil {
				return nil, err
			}

			switch k {
			case "COMMENT":
				kep.Comments = append(kep.Comments, v)
			case "SEMI_MAJOR_AXIS":
				kep.SemiMajorAxis, err = parseFloat(v)
				if err != nil {
					return nil, err
				}
			case "ECCENTRICITY":
				kep.Eccentricity, err = parseFloat(v)
				if err != nil {
					return nil, err
				}
			case "INCLINATION":
				kep.Inclination, err = parseFloat(v)
				if err != nil {
					return nil, err
				}
			case "RA_OF_ASC_NODE":
				kep.RaOfAscNode, err = parseFloat(v)
				if err != nil {
					return nil, err
				}
			case "ARG_OF_PERICENTER":
				kep.ArgOfPericenter, err = parseFloat(v)
				if err != nil {
					return nil, err
				}
			case "TRUE_ANOMALY":
				kep.TrueAnomaly, err = parseFloat(v)
				if err != nil {
					return nil, err
				}
			case "MEAN_ANOMALY":
				kep.MeanAnomaly, err = parseFloat(v)
				if err != nil {
					return nil, err
				}
			case "GM":
				kep.GM, err = parseFloat(v)
				if err != nil {
					return nil, err
				}
				foundGm = true
			}

			if foundGm {
				break
			}
		}

		if kep.MeanAnomaly != 0 && kep.TrueAnomaly != 0 {
			return nil, errors.New("either TRUE_ANOMALY or MEAN_ANOMALY must be provided, not both")
		}

		result.Data.OsculatingKeplerianElements = kep
	}

	return &result, nil
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

type OPM struct {
	Header   OPMHeader
	MetaData OPMMetaData
	Data     OPMData
	Raw      string
}

type OPMHeader struct {
	// Mandatory
	CcsdsOpmVers string `json:"CCSDS_OPM_VERS" xml:"CCSDS_OPM_VERS"`

	// Optional
	Comments []string `json:"COMMENT" xml:"COMMENT"`

	// Optional
	Classification string `json:"CLASSIFICATION" xml:"CLASSIFICATION"`

	// Mandatory
	CreationDate time.Time `json:"CREATION_DATE" xml:"CREATION_DATE"`

	// Mandatory
	Originator string `json:"ORIGINATOR" xml:"ORIGINATOR"`

	// Optional
	MessageId string `json:"MESSAGE_ID" xml:"MESSAGE_ID"`
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

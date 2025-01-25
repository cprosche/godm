package godm

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func ParseOPM(s string) (OPM, error) {
	result := OPM{}
	kvs, err := parseIntoKVs(s)
	if err != nil {
		return OPM{}, err
	}

	fields, err := getODMFields(&result)
	if err != nil {
		return OPM{}, err
	}

	kvIndex := 0
	for _, f := range fields {
		if kvIndex >= len(kvs) {
			if f.Required {
				return OPM{}, fmt.Errorf("expected key %s, got none, more fields than KVs", f.Name)
			}
			continue
		}

		kv := kvs[kvIndex]
		if f.Name != kv.Key && f.Required {
			return OPM{}, fmt.Errorf("expected key %s, got %s", f.Name, kv.Key)
		}

		if f.Name == kv.Key {
			if f.ReflectVal.Type().Name() == "Time" {
				t, err := parseTime(kv.Value)
				if err != nil {
					return OPM{}, err
				}
				f.ReflectVal.Set(reflect.ValueOf(t))
				kvIndex++
			} else {
				switch f.ReflectVal.Kind() {
				case reflect.String:
					f.ReflectVal.SetString(kv.Value)
					kvIndex++
				case reflect.Slice:
					for kvs[kvIndex].Key == f.Name {
						f.ReflectVal.Set(reflect.Append(f.ReflectVal, reflect.ValueOf(kvs[kvIndex].Value)))
						kvIndex++
					}
				case reflect.Float64:
					n, err := parseFloat(kv.Value)
					if err != nil {
						return OPM{}, err
					}
					f.ReflectVal.SetFloat(n)
					kvIndex++
				default:
					return OPM{}, fmt.Errorf("unsupported type %s", f.ReflectVal.Kind())
				}
			}
		}
	}

	return result, nil
}

type Field struct {
	Name       string
	ReflectVal reflect.Value
	Required   bool
}

func getODMFields(v interface{}) ([]Field, error) {
	var (
		fields = []Field{}
		val    = reflect.ValueOf(v).Elem()
	)

	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		if field.Type.Name() == "Time" {
			tag, required := parseOdmTag(field)
			if tag == "" {
				continue
			}
			fields = append(fields, Field{
				Name:       tag,
				ReflectVal: val.Field(i),
				Required:   required,
			})
			continue
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			f := val.Field(i).Addr().Interface()
			subFields, err := getODMFields(f)
			if err != nil {
				return nil, err
			}
			fields = append(fields, subFields...)
		default:
			tag, required := parseOdmTag(field)
			if tag == "" {
				continue
			}
			fields = append(fields, Field{
				Name:       tag,
				ReflectVal: val.Field(i),
				Required:   required,
			})
		}
	}

	return fields, nil
}

func parseOdmTag(f reflect.StructField) (string, bool) {
	tag := f.Tag.Get("odm")
	if tag == "" {
		return "", false
	}

	var (
		parts    = strings.Split(tag, ",")
		name     = parts[0]
		required = false
	)

	if len(parts) > 1 && parts[1] == "required" {
		required = true
	}

	return name, required
}

type OPM struct {
	Header   OPMHeader
	MetaData OPMMetaData
	Data     OPMData
	Raw      string
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
	// Mandatory
	StateVector StateVector

	// Optional, none or all
	OsculatingKeplerianElements OsculatingKeplerianElements // TODO: add validation for this

	// Optional, mass required if maneuver specified
	SpacecraftParameters SpacecraftParameters // TODO: add validation for this

	// // Optional
	// CovarianceMatrix CovarianceMatrix

	// // Optional, repeats for each maneuver
	// ManeuverParametersList []ManeuverParameters

	// // Optional, defined in an ICD, key must start with "USER_DEFINED_"
	// UserDefinedParameters UserDefinedParameters
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

type SpacecraftParameters struct {
	Comments      []string `odm:"COMMENT"`
	Mass          float64  `odm:"MASS"` // TODO: add validation for this, Conditional, required if maneuver specified, kg
	SolarRadArea  float64  `odm:"SOLAR_RAD_AREA"`
	SolarRadCoeff float64  `odm:"SOLAR_RAD_COEFF"`
	DragArea      float64  `odm:"DRAG_AREA"`
	DragCoeff     float64  `odm:"DRAG_COEFF"`
}

type CovarianceMatrix struct {
	Comments    []string `odm:"COMMENT"`
	CovRefFrame string   `odm:"COV_REF_FRAME"`
	CXX         float64  `odm:"CX_X"`
	CYX         float64  `odm:"CY_X"`
	CYY         float64  `odm:"CY_Y"`
	CZX         float64  `odm:"CZ_X"`
	CZY         float64  `odm:"CZ_Y"`
	CZZ         float64  `odm:"CZ_Z"`
	CXDOTX      float64  `odm:"CX_DOT_X"`
	CXDOTY      float64  `odm:"CX_DOT_Y"`
	CXDOTZ      float64  `odm:"CX_DOT_Z"`
	CXDOTXDOT   float64  `odm:"CX_DOT_X_DOT"`
	CYDOTX      float64  `odm:"CY_DOT_X"`
	CYDOTY      float64  `odm:"CY_DOT_Y"`
	CYDOTZ      float64  `odm:"CY_DOT_Z"`
	CYDOTXDOT   float64  `odm:"CY_DOT_X_DOT"`
	CYDOTYDOT   float64  `odm:"CY_DOT_Y_DOT"`
	CZDOTX      float64  `odm:"CZ_DOT_X"`
	CZDOTY      float64  `odm:"CZ_DOT_Y"`
	CZDOTZ      float64  `odm:"CZ_DOT_Z"`
	CZDOTXDOT   float64  `odm:"CZ_DOT_X_DOT"`
	CZDOTYDOT   float64  `odm:"CZ_DOT_Y_DOT"`
	CZDOTZDOT   float64  `odm:"CZ_DOT_Z_DOT"`
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

type UserDefinedParameters map[string]string

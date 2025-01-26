package godm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseOMM(t *testing.T) {

	exampleBasic := `CCSDS_OMM_VERS = 3.0
CREATION_DATE = 2020-065T16:00:00
ORIGINATOR = NOAA
MESSAGE_ID = OMM 202013719185
OBJECT_NAME = GOES 9
OBJECT_ID = 1995-025A
CENTER_NAME = EARTH
REF_FRAME = TEME
TIME_SYSTEM = UTC
MEAN_ELEMENT_THEORY = SGP/SGP4
EPOCH = 2020-064T10:34:41.4264
MEAN_MOTION = 1.00273272
ECCENTRICITY = 0.0005013
INCLINATION = 3.0539
RA_OF_ASC_NODE = 81.7939
ARG_OF_PERICENTER = 249.2363
MEAN_ANOMALY = 150.1602
GM = 398600.8
EPHEMERIS_TYPE = 0
CLASSIFICATION_TYPE = U
NORAD_CAT_ID = 23581
ELEMENT_SET_NO = 0925
REV_AT_EPOCH = 4316
BSTAR = 0.0001
MEAN_MOTION_DOT = -0.00000113
MEAN_MOTION_DDOT = 0.0`

	expectedBasic := OMM{
		Header: OMMHeader{
			CcsdsOmmVers: "3.0",
			CreationDate: time.Date(2020, 3, 5, 16, 0, 0, 0, time.UTC),
			Originator:   "NOAA",
			MessageId:    "OMM 202013719185",
		},
		MetaData: OMMMetaData{
			ObjectName:        "GOES 9",
			ObjectId:          "1995-025A",
			CenterName:        "EARTH",
			RefFrame:          "TEME",
			TimeSystem:        "UTC",
			MeanElementTheory: "SGP/SGP4",
		},
		Data: OMMData{
			MeanKeplerianElements: MeanKeplerianElements{
				Epoch:           time.Date(2020, 3, 4, 10, 34, 41, 426400000, time.UTC),
				MeanMotion:      1.00273272,
				Eccentricity:    0.0005013,
				Inclination:     3.0539,
				RaOfAscNode:     81.7939,
				ArgOfPericenter: 249.2363,
				MeanAnomaly:     150.1602,
				Gm:              398600.8,
			},
			TLERelatedParameters: TLERelatedParameters{
				EphemerisType:      "0",
				ClassificationType: "U",
				NoradCatId:         23581,
				ElementSetNo:       925,
				RevAtEpoch:         4316,
				BStar:              0.0001,
				MeanMotionDot:      -0.00000113,
				MeanMotionDdot:     0.0,
			},
		},
	}

	t.Run("Basic", func(t *testing.T) {
		got, err := ParseOMM(exampleBasic)
		assert.Nil(t, err)
		assert.Equal(t, expectedBasic, got)
	})

	exampleCov := `CCSDS_OMM_VERS = 3.0
CREATION_DATE = 2020-065T16:00:00
ORIGINATOR = NOAA
OBJECT_NAME = GOES 9
OBJECT_ID = 1995-025A
CENTER_NAME = EARTH
REF_FRAME = TEME
TIME_SYSTEM = UTC
MEAN_ELEMENT_THEORY = SGP/SGP4
EPOCH = 2020-064T10:34:41.4264
MEAN_MOTION = 1.00273272
ECCENTRICITY = 0.0005013
INCLINATION = 3.0539
RA_OF_ASC_NODE = 81.7939
ARG_OF_PERICENTER = 249.2363
MEAN_ANOMALY = 150.1602
GM = 398600.8
EPHEMERIS_TYPE = 0
CLASSIFICATION_TYPE = U
NORAD_CAT_ID = 23581
ELEMENT_SET_NO = 0925
REV_AT_EPOCH = 4316
BSTAR = 0.0001
MEAN_MOTION_DOT = -0.00000113
MEAN_MOTION_DDOT = 0.0
COV_REF_FRAME = TEME
CX_X = 3.331349476038534e-04
CY_X = 4.618927349220216e-04
CY_Y = 6.782421679971363e-04
CZ_X = -3.070007847730449e-04
CZ_Y = -4.221234189514228e-04
CZ_Z = 3.231931992380369e-04
CX_DOT_X = -3.349365033922630e-07
CX_DOT_Y = -4.686084221046758e-07
CX_DOT_Z = 2.484949578400095e-07
CX_DOT_X_DOT = 4.296022805587290e-10
CY_DOT_X = -2.211832501084875e-07
CY_DOT_Y = -2.864186892102733e-07
CY_DOT_Z = 1.798098699846038e-07
CY_DOT_X_DOT = 2.608899201686016e-10
CY_DOT_Y_DOT = 1.767514756338532e-10
CZ_DOT_X = -3.041346050686871e-07
CZ_DOT_Y = -4.989496988610662e-07
CZ_DOT_Z = 3.540310904497689e-07
CZ_DOT_X_DOT = 1.869263192954590e-10
CZ_DOT_Y_DOT = 1.008862586240695e-10
CZ_DOT_Z_DOT = 6.224444338635500e-10`

	expectedCov := OMM{
		Header: OMMHeader{
			CcsdsOmmVers: "3.0",
			CreationDate: time.Date(2020, 3, 5, 16, 0, 0, 0, time.UTC),
			Originator:   "NOAA",
		},
		MetaData: OMMMetaData{
			ObjectName:        "GOES 9",
			ObjectId:          "1995-025A",
			CenterName:        "EARTH",
			RefFrame:          "TEME",
			TimeSystem:        "UTC",
			MeanElementTheory: "SGP/SGP4",
		},
		Data: OMMData{
			MeanKeplerianElements: MeanKeplerianElements{
				Epoch:           time.Date(2020, 3, 4, 10, 34, 41, 426400000, time.UTC),
				MeanMotion:      1.00273272,
				Eccentricity:    0.0005013,
				Inclination:     3.0539,
				RaOfAscNode:     81.7939,
				ArgOfPericenter: 249.2363,
				MeanAnomaly:     150.1602,
				Gm:              398600.8,
			},
			TLERelatedParameters: TLERelatedParameters{
				EphemerisType:      "0",
				ClassificationType: "U",
				NoradCatId:         23581,
				ElementSetNo:       925,
				RevAtEpoch:         4316,
				BStar:              0.0001,
				MeanMotionDot:      -0.00000113,
				MeanMotionDdot:     0.0,
			},
			CovarianceMatrix: CovarianceMatrix{
				CovRefFrame: "TEME",
				CXX:         3.331349476038534e-04,
				CYX:         4.618927349220216e-04,
				CYY:         6.782421679971363e-04,
				CZX:         -3.070007847730449e-04,
				CZY:         -4.221234189514228e-04,
				CZZ:         3.231931992380369e-04,
				CXDOTX:      -3.349365033922630e-07,
				CXDOTY:      -4.686084221046758e-07,
				CXDOTZ:      2.484949578400095e-07,
				CXDOTXDOT:   4.296022805587290e-10,
				CYDOTX:      -2.211832501084875e-07,
				CYDOTY:      -2.864186892102733e-07,
				CYDOTZ:      1.798098699846038e-07,
				CYDOTXDOT:   2.608899201686016e-10,
				CYDOTYDOT:   1.767514756338532e-10,
				CZDOTX:      -3.041346050686871e-07,
				CZDOTY:      -4.989496988610662e-07,
				CZDOTZ:      3.540310904497689e-07,
				CZDOTXDOT:   1.869263192954590e-10,
				CZDOTYDOT:   1.008862586240695e-10,
				CZDOTZDOT:   6.224444338635500e-10,
			},
		},
	}

	t.Run("With Covariance Matrix", func(t *testing.T) {
		got, err := ParseOMM(exampleCov)
		assert.Nil(t, err)
		assert.Equal(t, expectedCov, got)
	})

	exampleUserDefined := `CCSDS_OMM_VERS = 3.0
CREATION_DATE = 2020-065T16:00:00
ORIGINATOR = NOAA
OBJECT_NAME = GOES 9
OBJECT_ID = 1995-025A
CENTER_NAME = EARTH
REF_FRAME = TEME
TIME_SYSTEM = UTC
MEAN_ELEMENT_THEORY = SGP/SGP4
EPOCH = 2020-064T10:34:41.4264
MEAN_MOTION = 1.00273272 [rev/day]
ECCENTRICITY = 0.0005013
INCLINATION = 3.0539 [deg]
RA_OF_ASC_NODE = 81.7939 [deg]
ARG_OF_PERICENTER = 249.2363 [deg]
MEAN_ANOMALY = 150.1602 [deg]
GM = 398600.8 [km**3/s**2]
EPHEMERIS_TYPE = 0
CLASSIFICATION_TYPE = U
NORAD_CAT_ID = 23581
ELEMENT_SET_NO = 0925
REV_AT_EPOCH = 4316
BSTAR = 0.0001 [1/ER]
MEAN_MOTION_DOT = -0.00000113 [rev/day**2]
MEAN_MOTION_DDOT = 0.0 [rev/day**3]
USER_DEFINED_EARTH_MODEL = WGS-84`

	expectedUserDefined := OMM{
		Header: OMMHeader{
			CcsdsOmmVers: "3.0",
			CreationDate: time.Date(2020, 3, 5, 16, 0, 0, 0, time.UTC),
			Originator:   "NOAA",
		},
		MetaData: OMMMetaData{
			ObjectName:        "GOES 9",
			ObjectId:          "1995-025A",
			CenterName:        "EARTH",
			RefFrame:          "TEME",
			TimeSystem:        "UTC",
			MeanElementTheory: "SGP/SGP4",
		},
		Data: OMMData{
			MeanKeplerianElements: MeanKeplerianElements{
				Epoch:           time.Date(2020, 3, 4, 10, 34, 41, 426400000, time.UTC),
				MeanMotion:      1.00273272,
				Eccentricity:    0.0005013,
				Inclination:     3.0539,
				RaOfAscNode:     81.7939,
				ArgOfPericenter: 249.2363,
				MeanAnomaly:     150.1602,
				Gm:              398600.8,
			},
			TLERelatedParameters: TLERelatedParameters{
				EphemerisType:      "0",
				ClassificationType: "U",
				NoradCatId:         23581,
				ElementSetNo:       925,
				RevAtEpoch:         4316,
				BStar:              0.0001,
				MeanMotionDot:      -0.00000113,
				MeanMotionDdot:     0.0,
			},
			UserDefinedParameters: map[string]string{
				"USER_DEFINED_EARTH_MODEL": "WGS-84",
			},
		},
	}

	t.Run("With User Defined Parameters", func(t *testing.T) {
		got, err := ParseOMM(exampleUserDefined)
		assert.Nil(t, err)
		assert.Equal(t, expectedUserDefined, got)
	})
}

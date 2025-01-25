package godm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseOPM(t *testing.T) {
	exampleSimple := `CCSDS_OPM_VERS = 3.0
CREATION_DATE = 2022-11-06T09:23:57
ORIGINATOR = JAXA

COMMENT GEOCENTRIC, CARTESIAN, EARTH FIXED
OBJECT_NAME = OSPREY 5
OBJECT_ID = 1998-999A
CENTER_NAME = EARTH
REF_FRAME = ITRF2000
TIME_SYSTEM = UTC

COMMENT This is the state vector
EPOCH = 2022-12-18T14:28:15.1172
X = 6503.514000
Y = 1239.647000
Z = -717.490000
X_DOT = -0.873160
Y_DOT = 8.740420
Z_DOT = -4.191076

MASS = 3000.000000
SOLAR_RAD_AREA = 18.770000
SOLAR_RAD_COEFF = 1.000000
DRAG_AREA = 18.770000
DRAG_COEFF = 2.500000`

	expectedSimple := OPM{
		Header: OPMHeader{
			CcsdsOpmVers: "3.0",
			CreationDate: time.Date(2022, 11, 6, 9, 23, 57, 0, time.UTC),
			Originator:   "JAXA",
		},
		MetaData: OPMMetaData{
			Comments:   []string{"GEOCENTRIC, CARTESIAN, EARTH FIXED"},
			ObjectName: "OSPREY 5",
			ObjectId:   "1998-999A",
			CenterName: "EARTH",
			RefFrame:   "ITRF2000",
			TimeSystem: "UTC",
		},
		Data: OPMData{
			StateVector: StateVector{
				Comments: []string{"This is the state vector"},
				Epoch:    time.Date(2022, 12, 18, 14, 28, 15, 117200000, time.UTC),
				X:        6503.514000,
				Y:        1239.647000,
				Z:        -717.490000,
				XDOT:     -0.873160,
				YDOT:     8.740420,
				ZDOT:     -4.191076,
			},
			SpacecraftParameters: SpacecraftParameters{
				Mass:          3000.000000,
				SolarRadArea:  18.770000,
				SolarRadCoeff: 1.000000,
				DragArea:      18.770000,
				DragCoeff:     2.500000,
			},
		},
	}

	t.Run("Parses simple", func(t *testing.T) {
		got, err := ParseOPM(exampleSimple)
		assert.Nil(t, err)
		assert.Equal(t, expectedSimple, got)
	})

}

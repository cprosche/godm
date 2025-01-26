package godm

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

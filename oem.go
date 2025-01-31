package godm

import "time"

type OEM struct {
	Header   OEMHeader
	MetaData OEMMetaData
	Data     OEMData

	Raw string
}

type OEMHeader struct {
	CcsdsOemVers   string   `odm:"CCSDS_OEM_VERS,required"`
	Comments       []string `odm:"COMMENT"`
	Classification string   `odm:"CLASSIFICATION"`
	CreationDate   string   `odm:"CREATION_DATE,required"`
	Originator     string   `odm:"ORIGINATOR,required"`
	MessageId      string   `odm:"MESSAGE_ID"`
}

type OEMMetaData struct {
	Comments         []string  `odm:"COMMENT"`
	ObjectName       string    `odm:"OBJECT_NAME,required"`
	ObjectId         string    `odm:"OBJECT_ID,required"`
	CenterName       string    `odm:"CENTER_NAME,required"`
	RefFrame         string    `odm:"REF_FRAME,required"`
	RefFrameEpoch    time.Time `odm:"REF_FRAME_EPOCH"` // TODO: add validation
	TimeSystem       string    `odm:"TIME_SYSTEM,required"`
	StartTime        time.Time `odm:"START_TIME,required"`
	UseableStartTime time.Time `odm:"USEABLE_START_TIME"`
	UseableStopTime  time.Time `odm:"USEABLE_STOP_TIME"`
	StopTime         time.Time `odm:"STOP_TIME,required"`
	Interpolation    string    `odm:"INTERPOLATION"`
	InterpolationDeg int       `odm:"INTERPOLATION_DEGREE"` // TODO: add validation
}

type OEMData struct {
	// data lines
	// covariance
}

package irsdk

import (
	"log"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v2"
)

type session struct {
	WeekendInfo   weekendInfo   `yaml:"WeekendInfo"`
	SessionInfo   sessionInfo   `yaml:"SessionInfo"`
	CameraInfo    cameraInfo    `yaml:"CameraInfo"`
	RadioInfo     radioInfo     `yaml:"RadioInfo"`
	DriverInfo    driverInfo    `yaml:"DriverInfo"`
	SplitTimeInfo splitTimeInfo `yaml:"SplitTimeInfo"`
	CarSetup      carSetup      `yaml:"CarSetup"`
}

type weekendInfo struct {
	TrackName              string           `yaml:"TrackName"`
	TrackID                int              `yaml:"TrackID"`
	TrackLength            string           `yaml:"TrackLength"`
	TrackDisplayName       string           `yaml:"TrackDisplayName"`
	TrackDisplayShortName  string           `yaml:"TrackDisplayShortName"`
	TrackConfigName        string           `yaml:"TrackConfigName"`
	TrackCity              string           `yaml:"TrackCity"`
	TrackCountry           string           `yaml:"TrackCountry"`
	TrackAltitude          string           `yaml:"TrackAltitude"`
	TrackLatitude          string           `yaml:"TrackLatitude"`
	TrackLongitude         string           `yaml:"TrackLongitude"`
	TrackNorthOffset       string           `yaml:"TrackNorthOffset"`
	TrackNumTurns          int              `yaml:"TrackNumTurns"`
	TrackPitSpeedLimit     string           `yaml:"TrackPitSpeedLimit"`
	TrackType              string           `yaml:"TrackType"`
	TrackDirection         string           `yaml:"TrackDirection"`
	TrackWeatherType       string           `yaml:"TrackWeatherType"`
	TrackSkies             string           `yaml:"TrackSkies"`
	TrackSurfaceTemp       string           `yaml:"TrackSurfaceTemp"`
	TrackAirTemp           string           `yaml:"TrackAirTemp"`
	TrackAirPressure       string           `yaml:"TrackAirPressure"`
	TrackWindVel           string           `yaml:"TrackWindVel"`
	TrackWindDir           string           `yaml:"TrackWindDir"`
	TrackRelativeHumidity  string           `yaml:"TrackRelativeHumidity"`
	TrackFogLevel          string           `yaml:"TrackFogLevel"`
	TrackCleanup           int              `yaml:"TrackCleanup"`
	TrackDynamicTrack      int              `yaml:"TrackDynamicTrack"`
	TrackVersion           string           `yaml:"TrackVersion"`
	SeriesID               int              `yaml:"SeriesID"`
	SeasonID               int              `yaml:"SeasonID"`
	SessionID              int              `yaml:"SessionID"`
	SubSessionID           int              `yaml:"SubSessionID"`
	LeagueID               int              `yaml:"LeagueID"`
	Official               int              `yaml:"Official"`
	RaceWeek               int              `yaml:"RaceWeek"`
	EventType              string           `yaml:"EventType"`
	Category               string           `yaml:"Category"`
	SimMode                string           `yaml:"SimMode"`
	TeamRacing             int              `yaml:"TeamRacing"`
	MinDrivers             int              `yaml:"MinDrivers"`
	MaxDrivers             int              `yaml:"MaxDrivers"`
	DCRuleSet              string           `yaml:"DCRuleSet"`
	QualifierMustStartRace int              `yaml:"QualifierMustStartRace"`
	NumCarClasses          int              `yaml:"NumCarClasses"`
	NumCarTypes            int              `yaml:"NumCarTypes"`
	HeatRacing             int              `yaml:"HeatRacing"`
	BuildType              string           `yaml:"BuildType"`
	BuildTarget            string           `yaml:"BuildTarget"`
	BuildVersion           string           `yaml:"BuildVersion"`
	WeekendOptions         weekendOptions   `yaml:"WeekendOptions"`
	TelemetryOptions       telemetryOptions `yaml:"TelemetryOptions"`
}

type sessionInfo struct {
	Sessions []sessionDetails `yaml:"Sessions"`
}

type cameraInfo struct {
	Groups []cameraGroup `yaml:"Groups"`
}

type radioInfo struct {
	SelectedRadioNum int     `yaml:"SelectedRadioNum"`
	Radios           []radio `yaml:"Radios"`
}

type driverInfo struct {
	DriverCarIdx              int      `yaml:"DriverCarIdx"`
	DriverUserID              int      `yaml:"DriverUserID"`
	PaceCarIdx                int      `yaml:"PaceCarIdx"`
	DriverHeadPosX            float64  `yaml:"DriverHeadPosX"`
	DriverHeadPosY            float64  `yaml:"DriverHeadPosY"`
	DriverHeadPosZ            float64  `yaml:"DriverHeadPosZ"`
	DriverCarIdleRPM          float64  `yaml:"DriverCarIdleRPM"`
	DriverCarRedLine          float64  `yaml:"DriverCarRedLine"`
	DriverCarEngCylinderCount int      `yaml:"DriverCarEngCylinderCount"`
	DriverCarFuelKgPerLtr     float64  `yaml:"DriverCarFuelKgPerLtr"`
	DriverCarFuelMaxLtr       float64  `yaml:"DriverCarFuelMaxLtr"`
	DriverCarMaxFuelPct       float64  `yaml:"DriverCarMaxFuelPct"`
	DriverCarGearNumForward   int      `yaml:"DriverCarGearNumForward"`
	DriverCarGearNeutral      int      `yaml:"DriverCarGearNeutral"`
	DriverCarGearReverse      int      `yaml:"DriverCarGearReverse"`
	DriverCarSLFirstRPM       float64  `yaml:"DriverCarSLFirstRPM"`
	DriverCarSLShiftRPM       float64  `yaml:"DriverCarSLShiftRPM"`
	DriverCarSLLastRPM        float64  `yaml:"DriverCarSLLastRPM"`
	DriverCarSLBlinkRPM       float64  `yaml:"DriverCarSLBlinkRPM"`
	DriverCarVersion          string   `yaml:"DriverCarVersion"`
	DriverPitTrkPct           float64  `yaml:"DriverPitTrkPct"`
	DriverCarEstLapTime       float64  `yaml:"DriverCarEstLapTime"`
	DriverSetupName           string   `yaml:"DriverSetupName"`
	DriverSetupIsModified     int      `yaml:"DriverSetupIsModified"`
	DriverSetupLoadTypeName   string   `yaml:"DriverSetupLoadTypeName"`
	DriverSetupPassedTech     int      `yaml:"DriverSetupPassedTech"`
	DriverIncidentCount       int      `yaml:"DriverIncidentCount"`
	Drivers                   []driver `yaml:"Drivers"`
}

type splitTimeInfo struct {
	Sectors []sector `yaml:"Sectors"`
}

type carSetup struct {
	UpdateCount int `yaml:"UpdateCount"`
	TiresAero   struct {
		LeftFront  leftTire  `yaml:"LeftFront"`
		LeftRear   leftTire  `yaml:"LeftRear"`
		RightFront rightTire `yaml:"RightFront"`
		RightRear  rightTire `yaml:"RightRear"`
	} `yaml:"TiresAero"`
	Chassis struct {
		Front      frontChassis     `yaml:"Front"`
		LeftFront  sideFrontChassis `yaml:"LeftFront"`
		LeftRear   sideRearChassis  `yaml:"LeftRear"`
		InCarDials inCarSettings    `yaml:"InCarDials"`
		RightFront sideFrontChassis `yaml:"RightFront"`
		RightRear  sideRearChassis  `yaml:"RightRear"`
		Rear       rearChassis      `yaml:"Rear"`
	} `yaml:"Chassis"`
}

type driver struct {
	CarIdx                  int     `yaml:"CarIdx"`
	UserName                string  `yaml:"UserName"`
	AbbrevName              string  `yaml:"AbbrevName"`
	Initials                string  `yaml:"Initials"`
	UserID                  int     `yaml:"UserID"`
	TeamID                  int     `yaml:"TeamID"`
	TeamName                string  `yaml:"TeamName"`
	CarNumber               string  `yaml:"CarNumber"`
	CarNumberRaw            int     `yaml:"CarNumberRaw"`
	CarPath                 string  `yaml:"CarPath"`
	CarClassID              int     `yaml:"CarClassID"`
	CarID                   int     `yaml:"CarID"`
	CarIsPaceCar            int     `yaml:"CarIsPaceCar"`
	CarIsAI                 int     `yaml:"CarIsAI"`
	CarScreenName           string  `yaml:"CarScreenName"`
	CarScreenNameShort      string  `yaml:"CarScreenNameShort"`
	CarClassShortName       string  `yaml:"CarClassShortName"`
	CarClassRelSpeed        int     `yaml:"CarClassRelSpeed"`
	CarClassLicenseLevel    int     `yaml:"CarClassLicenseLevel"`
	CarClassMaxFuelPct      string  `yaml:"CarClassMaxFuelPct"`
	CarClassWeightPenalty   string  `yaml:"CarClassWeightPenalty"`
	CarClassPowerAdjust     string  `yaml:"CarClassPowerAdjust"`
	CarClassDryTireSetLimit string  `yaml:"CarClassDryTireSetLimit"`
	CarClassColor           int     `yaml:"CarClassColor"`
	CarClassEstLapTime      float64 `yaml:"CarClassEstLapTime"`
	IRating                 int     `yaml:"IRating"`
	LicLevel                int     `yaml:"LicLevel"`
	LicSubLevel             int     `yaml:"LicSubLevel"`
	LicString               string  `yaml:"LicString"`
	LicColor                string  `yaml:"LicColor"`
	IsSpectator             int     `yaml:"IsSpectator"`
	CarDesignStr            string  `yaml:"CarDesignStr"`
	HelmetDesignStr         string  `yaml:"HelmetDesignStr"`
	SuitDesignStr           string  `yaml:"SuitDesignStr"`
	CarNumberDesignStr      string  `yaml:"CarNumberDesignStr"`
	CarSponsor1             int     `yaml:"CarSponsor_1"`
	CarSponsor2             int     `yaml:"CarSponsor_2"`
	CurDriverIncidentCount  int     `yaml:"CurDriverIncidentCount"`
	TeamIncidentCount       int     `yaml:"TeamIncidentCount"`
}

type weekendOptions struct {
	NumStarters                int    `yaml:"NumStarters"`
	StartingGrid               string `yaml:"StartingGrid"`
	QualifyScoring             string `yaml:"QualifyScoring"`
	CourseCautions             string `yaml:"CourseCautions"`
	StandingStart              int    `yaml:"StandingStart"`
	ShortParadeLap             int    `yaml:"ShortParadeLap"`
	Restarts                   string `yaml:"Restarts"`
	WeatherType                string `yaml:"WeatherType"`
	Skies                      string `yaml:"Skies"`
	WindDirection              string `yaml:"WindDirection"`
	WindSpeed                  string `yaml:"WindSpeed"`
	WeatherTemp                string `yaml:"WeatherTemp"`
	RelativeHumidity           string `yaml:"RelativeHumidity"`
	FogLevel                   string `yaml:"FogLevel"`
	TimeOfDay                  string `yaml:"TimeOfDay"`
	Date                       string `yaml:"Date"`
	EarthRotationSpeedupFactor int    `yaml:"EarthRotationSpeedupFactor"`
	Unofficial                 int    `yaml:"Unofficial"`
	CommercialMode             string `yaml:"CommercialMode"`
	NightMode                  string `yaml:"NightMode"`
	IsFixedSetup               int    `yaml:"IsFixedSetup"`
	StrictLapsChecking         string `yaml:"StrictLapsChecking"`
	HasOpenRegistration        int    `yaml:"HasOpenRegistration"`
	HardcoreLevel              int    `yaml:"HardcoreLevel"`
	NumJokerLaps               int    `yaml:"NumJokerLaps"`
	IncidentLimit              string `yaml:"IncidentLimit"`
	FastRepairsLimit           string `yaml:"FastRepairsLimit"`
	GreenWhiteCheckeredLimit   int    `yaml:"GreenWhiteCheckeredLimit"`
}

type telemetryOptions struct {
	TelemetryDiskFile string `yaml:"TelemetryDiskFile"`
}

type sessionDetails struct {
	SessionNum              int                 `yaml:"SessionNum"`
	SessionLaps             string              `yaml:"SessionLaps"`
	SessionTime             string              `yaml:"SessionTime"`
	SessionNumLapsToAvg     int                 `yaml:"SessionNumLapsToAvg"`
	SessionType             string              `yaml:"SessionType"`
	SessionTrackRubberState string              `yaml:"SessionTrackRubberState"`
	SessionName             string              `yaml:"SessionName"`
	SessionSubType          interface{}         `yaml:"SessionSubType"`
	SessionSkipped          int                 `yaml:"SessionSkipped"`
	SessionRunGroupsUsed    int                 `yaml:"SessionRunGroupsUsed"`
	ResultsPositions        []resultsPosition   `yaml:"ResultsPositions"`
	ResultsFastestLap       []resultsFastestLap `yaml:"ResultsFastestLap"`
	ResultsAverageLapTime   int                 `yaml:"ResultsAverageLapTime"`
	ResultsNumCautionFlags  int                 `yaml:"ResultsNumCautionFlags"`
	ResultsNumCautionLaps   int                 `yaml:"ResultsNumCautionLaps"`
	ResultsNumLeadChanges   int                 `yaml:"ResultsNumLeadChanges"`
	ResultsLapsComplete     int                 `yaml:"ResultsLapsComplete"`
	ResultsOfficial         int                 `yaml:"ResultsOfficial"`
}

type resultsFastestLap struct {
	CarIdx      int `yaml:"CarIdx"`
	FastestLap  int `yaml:"FastestLap"`
	FastestTime int `yaml:"FastestTime"`
}

type resultsPosition struct {
	Position          int     `yaml:"Position"`
	ClassPosition     int     `yaml:"ClassPosition"`
	CarIdx            int     `yaml:"CarIdx"`
	Lap               int     `yaml:"Lap"`
	Time              float64 `yaml:"Time"`
	FastestLap        int     `yaml:"FastestLap"`
	FastestTime       float64 `yaml:"FastestTime"`
	LastTime          float64 `yaml:"LastTime"`
	LapsLed           int     `yaml:"LapsLed"`
	LapsComplete      int     `yaml:"LapsComplete"`
	JokerLapsComplete int     `yaml:"JokerLapsComplete"`
	LapsDriven        float64 `yaml:"LapsDriven"`
	Incidents         int     `yaml:"Incidents"`
	ReasonOutId       int     `yaml:"ReasonOutId"`
	ReasonOutStr      string  `yaml:"ReasonOutStr"`
}

type radio struct {
	RadioNum            int         `yaml:"RadioNum"`
	HopCount            int         `yaml:"HopCount"`
	NumFrequencies      int         `yaml:"NumFrequencies"`
	TunedToFrequencyNum int         `yaml:"TunedToFrequencyNum"`
	ScanningIsOn        int         `yaml:"ScanningIsOn"`
	Frequencies         []frequency `yaml:"Frequencies"`
}

type frequency struct {
	FrequencyNum  int    `yaml:"FrequencyNum"`
	FrequencyName string `yaml:"FrequencyName"`
	Priority      int    `yaml:"Priority"`
	CarIdx        int    `yaml:"CarIdx"`
	EntryIdx      int    `yaml:"EntryIdx"`
	ClubID        int    `yaml:"ClubID"`
	CanScan       int    `yaml:"CanScan"`
	CanSquawk     int    `yaml:"CanSquawk"`
	Muted         int    `yaml:"Muted"`
	IsMutable     int    `yaml:"IsMutable"`
	IsDeletable   int    `yaml:"IsDeletable"`
}

type camera struct {
	CameraNum  int    `yaml:"CameraNum"`
	CameraName string `yaml:"CameraName"`
}

type sector struct {
	SectorNum      int     `yaml:"SectorNum"`
	SectorStartPct float64 `yaml:"SectorStartPct"`
}

type cameraGroup struct {
	GroupNum  int      `yaml:"GroupNum"`
	GroupName string   `yaml:"GroupName"`
	Cameras   []camera `yaml:"Cameras"`
	IsScenic  bool     `yaml:"IsScenic,omitempty"`
}

type leftTire struct {
	StartingPressure string `yaml:"StartingPressure"`
	LastHotPressure  string `yaml:"LastHotPressure"`
	LastTempsOMI     string `yaml:"LastTempsOMI"`
	TreadRemaining   string `yaml:"TreadRemaining"`
}

type rightTire struct {
	StartingPressure string `yaml:"StartingPressure"`
	LastHotPressure  string `yaml:"LastHotPressure"`
	LastTempsIMO     string `yaml:"LastTempsIMO"`
	TreadRemaining   string `yaml:"TreadRemaining"`
}

type sideFrontChassis struct {
	CornerWeight      string `yaml:"CornerWeight"`
	RideHeight        string `yaml:"RideHeight"`
	SpringPerchOffset string `yaml:"SpringPerchOffset"`
	Camber            string `yaml:"Camber"`
}

type sideRearChassis struct {
	CornerWeight      string `yaml:"CornerWeight"`
	RideHeight        string `yaml:"RideHeight"`
	SpringPerchOffset string `yaml:"SpringPerchOffset"`
	Camber            string `yaml:"Camber"`
	ToeIn             string `yaml:"ToeIn"`
}

type frontChassis struct {
	ArbSetting  int    `yaml:"ArbSetting"`
	ToeIn       string `yaml:"ToeIn"`
	FuelLevel   string `yaml:"FuelLevel"`
	CrossWeight string `yaml:"CrossWeight"`
}

type rearChassis struct {
	ArbSetting  int `yaml:"ArbSetting"`
	WingSetting int `yaml:"WingSetting"`
}

type inCarSettings struct {
	DisplayPage       string `yaml:"DisplayPage"`
	BrakePressureBias string `yaml:"BrakePressureBias"`
}

func readSessionData(sdk *IRSDK) string {
	dec := charmap.Windows1252.NewDecoder()

	rbuf := make([]byte, sdk.Header.SessionInfoLen)

	_, err := sdk.Reader.ReadAt(rbuf, int64(sdk.Header.SessionInfoOffset))
	if err != nil {
		log.Fatal(err)
	}

	rbuf, err = dec.Bytes(rbuf)
	if err != nil {
		log.Fatal(err)
	}

	yaml := strings.TrimRight(string(rbuf[:sdk.Header.SessionInfoLen]), "\x00")
	return yaml
}

// This function updates the session data in the sdk struct
func updateSessionData(sdk *IRSDK) {
	sRaw := readSessionData(sdk)

	newSession := session{}
	err := yaml.Unmarshal([]byte(sRaw), &newSession)
	if err != nil {
		log.Fatal(err)
	}

	sdk.Session = &newSession
}

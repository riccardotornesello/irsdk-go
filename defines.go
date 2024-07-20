package irsdk

const (
	DataValidEventName = "Local\\IRSDKDataValidEvent"
	MemMapFile         = "Local\\IRSDKMemMapFileName"
	BroadcastMsgName   = "BROADCASTMSG"
	MemMapSize         = 1164 * 1024
	MaxBufs            = 4
	MaxString          = 32
	MaxDesc            = 64
	UnlimitedLaps      = 32767
	UnlimitedTime      = 604800.0
	Ver                = 2
)

type StatusField = int

const (
	stConnected StatusField = 1
)

type TrackLocation int

const (
	LocationNotInWorld TrackLocation = iota - 1
	LocationOffTrack
	LocationInPitStall
	LocationAproachingPits
	LocationOnTrack
)

type TrackSurface int

const (
	SurfaceNotInWorld TrackSurface = iota - 1
	SurfaceUndefinedMaterial
	SurfaceAsphalt1Material
	SurfaceAsphalt2Material
	SurfaceAsphalt3Material
	SurfaceAsphalt4Material
	SurfaceConcrete1Material
	SurfaceConcrete2Material
	SurfaceRacingDirt1Material
	SurfaceRacingDirt2Material
	SurfacePaint1Material
	SurfacePaint2Material
	SurfaceRumble1Material
	SurfaceRumble2Material
	SurfaceRumble3Material
	SurfaceRumble4Material
	SurfaceGrass1Material
	SurfaceGrass2Material
	SurfaceGrass3Material
	SurfaceGrass4Material
	SurfaceDirt1Material
	SurfaceDirt2Material
	SurfaceDirt3Material
	SurfaceDirt4Material
	SurfaceSandMaterial
	SurfaceGravel1Material
	SurfaceGravel2Material
	SurfaceGrasscreteMaterial
	SurfaceAstroturfMaterial
)

type SessionState int

const (
	SessionStateInvalid SessionState = iota
	SessionStateGetInCar
	SessionStateWarmup
	SessionStateParadeLaps
	SessionStateRacing
	SessionStateCheckered
	SessionStateCoolDown
)

type CarLeftRight int

const (
	LROff CarLeftRight = iota
	LRClear
	LRCarLeft
	LRCarRight
	LRCarLeftRight
	LR2CarsLeft
	LR2CarsRight
)

type PitSvStatus int

const (
	PitSvNone PitSvStatus = iota
	PitSvInProgress
	PitSvComplete
	PitSvTooFarLeft    = 100
	PitSvTooFarRight   = 101
	PitSvTooFarForward = 102
	PitSvTooFarBack    = 103
	PitSvBadAngle      = 104
	PitSvCantFixThat   = 105
)

type PaceMode int

const (
	PaceModeSingleFileStart PaceMode = iota
	PaceModeDoubleFileStart
	PaceModeSingleFileRestart
	PaceModeDoubleFileRestart
	PaceModeNotPacing
)

type TrackWetness int

const (
	TrackWetness_UNKNOWN TrackWetness = iota
	TrackWetness_Dry
	TrackWetness_MostlyDry
	TrackWetness_VeryLightlyWet
	TrackWetness_LightlyWet
	TrackWetness_ModeratelyWet
	TrackWetness_VeryWet
	TrackWetness_ExtremelyWet
)

type EngineWarnings int

const (
	waterTempWarning EngineWarnings = 1 << iota
	fuelPressureWarning
	oilPressureWarning
	engineStalled
	pitSpeedLimiter
	revLimiterActive
	oilTempWarning
)

type Flags int

const (
	checkered Flags = 1 << iota
	white
	green
	yellow
	red
	blue
	debris
	crossed
	yellowWaving
	oneLapToGreen
	greenHeld
	tenToGo
	fiveToGo
	randomWaving
	caution
	cautionWaving
	black
	disqualify
	servicible
	furled
	repair
	startHidden
	startReady
	startSet
	startGo
)

type CameraState int

const (
	IsSessionScreen CameraState = 1 << iota
	IsScenicActive
	CamToolActive
	UIHidden
	UseAutoShotSelection
	UseTemporaryEdits
	UseKeyAcceleration
	UseKey10xAcceleration
	UseMouseAimMode
)

type PitSvFlags int

const (
	LFTireChange PitSvFlags = 1 << iota
	RFTireChange
	LRTireChange
	RRTireChange
	FuelFill
	WindshieldTearoff
	FastRepair
)

type PaceFlags int

const (
	PaceFlagsEndOfLine PaceFlags = 1 << iota
	PaceFlagsFreePass
	PaceFlagsWavedAround
)

type DiskSubHeader struct {
	SessionStartDate   int64
	SessionStartTime   float64
	SessionEndTime     float64
	SessionLapCount    int
	SessionRecordCount int
}

type BroadcastMsg int

const (
	BroadcastCamSwitchPos            BroadcastMsg = iota // car position, group, camera
	BroadcastCamSwitchNum                                // driver #, group, camera
	BroadcastCamSetState                                 // CameraState, unused, unused
	BroadcastReplaySetPlaySpeed                          // speed, slowMotion, unused
	BroadcastReplaySetPlayPosition                       // RpyPosMode, Frame Number (high, low)
	BroadcastReplaySearch                                // RpySrchMode, unused, unused
	BroadcastReplaySetState                              // RpyStateMode, unused, unused
	BroadcastReloadTextures                              // ReloadTexturesMode, carIdx, unused
	BroadcastChatComand                                  // ChatCommandMode, subCommand, unused
	BroadcastPitCommand                                  // PitCommandMode, parameter
	BroadcastTelemCommand                                // TelemCommandMode, unused, unused
	BroadcastFFBCommand                                  // FFBCommandMode, value (float, high, low)
	BroadcastReplaySearchSessionTime                     // sessionNum, sessionTimeMS (high, low)
	BroadcastVideoCapture                                // VideoCaptureMode, unused, unused
	BroadcastLast                                        // unused placeholder
)

type ChatCommandMode int

const (
	ChatCommand_Macro     ChatCommandMode = iota // pass in a number from 1-15 representing the chat macro to launch
	ChatCommand_BeginChat                        // Open up a new chat window
	ChatCommand_Reply                            // Reply to last private chat
	ChatCommand_Cancel                           // Close chat window
)

type PitCommandMode int

const (
	PitCommand_Clear      PitCommandMode = iota // Clear all pit checkboxes
	PitCommand_WS                               // Clean the winshield, using one tear off
	PitCommand_Fuel                             // Add fuel, optionally specify the amount to add in liters or pass '0' to use existing amount
	PitCommand_LF                               // Change the left front tire, optionally specifying the pressure in KPa or pass '0' to use existing pressure
	PitCommand_RF                               // right front
	PitCommand_LR                               // left rear
	PitCommand_RR                               // right rear
	PitCommand_ClearTires                       // Clear tire pit checkboxes
	PitCommand_FR                               // Request a fast repair
	PitCommand_ClearWS                          // Uncheck Clean the winshield checkbox
	PitCommand_ClearFR                          // Uncheck request a fast repair
	PitCommand_ClearFuel                        // Uncheck add fuel
	PitCommand_TC                               // Change tire compound
)

type TelemCommandMode int

const (
	TelemCommand_Stop    TelemCommandMode = iota // Turn telemetry recording off
	TelemCommand_Start                           // Turn telemetry recording on
	TelemCommand_Restart                         // Write current file to disk and start a new one
)

type RpyStateMode int

const (
	RpyState_EraseTape RpyStateMode = iota // clear any data in the replay tape
	RpyState_Last                          // unused place holder
)

type ReloadTexturesMode int

const (
	ReloadTextures_All    ReloadTexturesMode = iota // reload all textuers
	ReloadTextures_CarIdx                           // reload only textures for the specific carIdx
)

type RpySrchMode int

const (
	RpySrch_ToStart RpySrchMode = iota
	RpySrch_ToEnd
	RpySrch_PrevSession
	RpySrch_NextSession
	RpySrch_PrevLap
	RpySrch_NextLap
	RpySrch_PrevFrame
	RpySrch_NextFrame
	RpySrch_PrevIncident
	RpySrch_NextIncident
	RpySrch_Last
)

type RpyPosMode int

const (
	RpyPos_Begin RpyPosMode = iota
	RpyPos_Current
	RpyPos_End
	RpyPos_Last
)

type FFBCommandMode int

const (
	FFBCommand_MaxForce FFBCommandMode = iota // Set the maximum force when mapping steering torque force to direct input units (float in Nm)
	FFBCommand_Last                           // unused placeholder
)

type csMode int

const (
	csFocusAtIncident csMode = -3
	csFocusAtLeader   csMode = -2
	csFocusAtExiting  csMode = -1
	csFocusAtDriver   csMode = 0
)

type VideoCaptureMode int

const (
	VideoCapture_TriggerScreenShot  VideoCaptureMode = iota // save a screenshot to disk
	VideoCaptuer_StartVideoCapture                          // start capturing video
	VideoCaptuer_EndVideoCapture                            // stop capturing video
	VideoCaptuer_ToggleVideoCapture                         // toggle video capture on/off
	VideoCaptuer_ShowVideoTimer                             // show video timer in upper left corner of display
	VideoCaptuer_HideVideoTimer                             // hide video timer
)

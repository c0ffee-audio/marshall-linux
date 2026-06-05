package protocol

// ANCMode represents the Active Noise Cancelling state.
type ANCMode byte

const (
	ANCOff          ANCMode = 0x00
	ANCCancelling   ANCMode = 0x01
	ANCTransparency ANCMode = 0x02
)

func (m ANCMode) String() string {
	switch m {
	case ANCOff:
		return "off"
	case ANCCancelling:
		return "anc"
	case ANCTransparency:
		return "transparency"
	default:
		return "unknown"
	}
}

func EncodeANC(mode ANCMode) []byte {
	return []byte{byte(mode)}
}

func DecodeANC(data []byte) ANCMode {
	if len(data) != 1 {
		return ANCOff
	}
	return ANCMode(data[0])
}

// EQPreset represents a built-in EQ profile.
type EQPreset byte

const (
	EQFlat            EQPreset = 0
	EQCustom          EQPreset = 1
	EQRock            EQPreset = 2
	EQMetal           EQPreset = 3
	EQPop             EQPreset = 4
	EQHipHop          EQPreset = 5
	EQElectronic      EQPreset = 6
	EQJazz            EQPreset = 7
	EQBassBoost       EQPreset = 8
	EQMidBoost        EQPreset = 9
	EQTrebleBoost     EQPreset = 10
	EQLoudPushWorkout EQPreset = 11
)

func (p EQPreset) String() string {
	names := []string{
		"flat", "custom", "rock", "metal", "pop",
		"hiphop", "electronic", "jazz",
		"bass-boost", "mid-boost", "treble-boost", "loud-push-workout",
	}
	if int(p) < len(names) {
		return names[p]
	}
	return "unknown"
}

// EncodeEQAssign retourne les bytes pour ASSIGN_STEP_PRESET (step 0, preset).
func EncodeEQAssign(preset EQPreset) []byte {
	return []byte{0x01, 0x00, byte(preset)}
}

// EncodeEQActivate retourne les bytes pour CHANGE_ACTIVE_STEP (activer step 0).
func EncodeEQActivate() []byte {
	return []byte{0x00, 0x00}
}

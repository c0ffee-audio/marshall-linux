package protocol

// UUID suffix constants for Zound Industries BLE protocol.
// Derived from static analysis of com.zoundindustries.marshallbt v3.7.1.
const (
	zoundSuffix  = "-1337-1dea-feed-c0ffee70c0de"
	legacySuffix = "-0000-1000-8000-00805f9b34fb"
	gaiaSuffix   = "-d102-11e1-9b23-00025b00a5a5"
)

func zoundUUID(code string) string  { return "0000" + code + zoundSuffix }
func legacyUUID(code string) string { return "0000" + code + legacySuffix }

// ServiceUUID is the primary Zound GATT service UUID.
const ServiceUUID = "FA302D24-D775-4343-B9ED-8CC68ACE3284"

// Characteristic UUIDs (modern Zound format).
var (
	CharANCConfiguration         = zoundUUID("0013")
	CharANCTransparencyValue     = zoundUUID("0019")
	CharANCNoiseCancellingValue  = zoundUUID("001a")
	CharVolume                   = zoundUUID("0007")
	CharVolumeLimit              = zoundUUID("0008")
	CharEqualizerSettings        = zoundUUID("0017")
	CharEqualizerCustomPreset    = zoundUUID("0018")
	CharGraphicalEqualizer       = zoundUUID("000f")
	CharToneControl              = zoundUUID("0025")
	CharRename                   = zoundUUID("0003")
	CharAutoPowerOff             = zoundUUID("0004")
	CharAutoOffTimer             = zoundUUID("0032")
	CharBTPairingMode            = zoundUUID("0001")
	CharActionButtonEvent        = zoundUUID("000c")
	CharActionButtonConfig       = zoundUUID("000d")
	CharAudioControl             = zoundUUID("0009")
	CharAudioNowPlaying          = zoundUUID("000a")
	CharAudioSource              = zoundUUID("001b")
	CharUISOunds                 = zoundUUID("000b")
	CharTouchLock                = zoundUUID("0014")
	CharPartyMode                = zoundUUID("001c")
	CharEcoCharging              = zoundUUID("001d")
	CharRoomPlacement            = zoundUUID("001e")
	CharNightMode                = zoundUUID("001f")
	CharWearSensorStatus         = zoundUUID("0027")
	CharWearSensorAction         = zoundUUID("0028")
	CharBatteryPreservation      = zoundUUID("002f")
	CharDynamicAudio             = zoundUUID("0030")
	CharSoundstage               = zoundUUID("0033")
	CharBTConnectionControl      = zoundUUID("0034")
	CharLEDIntensity             = zoundUUID("003a")
	CharLEAudioConfig            = zoundUUID("003d")
	CharAudioFeatureConfig       = zoundUUID("0044")
	CharAudioInputFeatureConfig  = zoundUUID("0045")
	CharUSBConfig                = zoundUUID("0048")
)

// Standard Bluetooth characteristic UUIDs.
var (
	CharManufacturerName  = legacyUUID("2a29")
	CharModelName         = legacyUUID("2a24")
	CharSerialNumber      = legacyUUID("2a25")
	CharFirmwareRevision  = legacyUUID("2a26")
	CharHardwareRevision  = legacyUUID("2a27")
	CharBatteryLevel      = legacyUUID("2a19")
	CharBatteryLevelStatus = legacyUUID("2bed")
)

// Earbuds battery UUIDs (Tymphany SDK).
const (
	CharRightEarbudBattery       = "7a573e5d-9330-4d9b-8660-63c33fc50001"
	CharLeftEarbudBattery        = "7a573e5d-9330-4d9b-8660-63c33fc50002"
	CharCaseBattery              = "7a573e5d-9330-4d9b-8660-63c33fc50003"
	CharRightEarbudBatteryStatus = "7a573e5d-9330-4d9b-8660-63c33fc50101"
	CharLeftEarbudBatteryStatus  = "7a573e5d-9330-4d9b-8660-63c33fc50102"
	CharCaseBatteryStatus        = "7a573e5d-9330-4d9b-8660-63c33fc50103"
)

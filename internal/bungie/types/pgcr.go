package bungie

import (
	"encoding/json"
	"time"
)

type PGCR struct {
	Archived                        time.Time       `json:"archived"`
	Period                          time.Time       `json:"period"`
	StartingPhaseIndex              json.Number     `json:"startingPhaseIndex"`
	ActivityWasStartedFromBeginning bool            `json:"activityWasStartedFromBeginning"`
	ActivityDetails                 ActivityDetails `json:"activityDetails"`
	Entries                         []Entry         `json:"entries"`
	Teams                           []Team          `json:"teams"`
}

type ActivityDetails struct {
	ReferenceID          json.Number   `json:"referenceId"`
	DirectorActivityHash json.Number   `json:"directorActivityHash"`
	InstanceID           json.Number   `json:"instanceId"`
	Mode                 json.Number   `json:"mode"`
	Modes                []json.Number `json:"modes"`
	IsPrivate            bool          `json:"isPrivate"`
	MembershipType       json.Number   `json:"membershipType"`
}

type Entry struct {
	Standing    int      `json:"standing"`
	Score       int      `json:"score"`
	Player      Player   `json:"player"`
	CharacterID string   `json:"characterId"`
	Values      Values   `json:"values"`
	Extended    Extended `json:"extended"`
}

type Player struct {
	DestinyUserInfo DestinyUserInfo `json:"destinyUserInfo"`
	CharacterClass  string          `json:"characterClass"`
	ClassHash       uint32          `json:"classHash"`
	RaceHash        uint32          `json:"raceHash"`
	GenderHash      uint32          `json:"genderHash"`
	CharacterLevel  int             `json:"characterLevel"`
	LightLevel      int             `json:"lightLevel"`
	EmblemHash      uint32          `json:"emblemHash"`
}

type DestinyUserInfo struct {
	IconPath                    string      `json:"iconPath"`
	CrossSaveOverride           int         `json:"crossSaveOverride"`
	ApplicableMembershipTypes   []int       `json:"applicableMembershipTypes"`
	IsPublic                    bool        `json:"isPublic"`
	MembershipType              int         `json:"membershipType"`
	MembershipID                json.Number `json:"membershipId,omitzero"`
	DisplayName                 string      `json:"displayName"`
	BungieGlobalDisplayName     string      `json:"bungieGlobalDisplayName"`
	BungieGlobalDisplayNameCode int         `json:"bungieGlobalDisplayNameCode"`
}

type Values struct {
	Assists           float64 `json:"assists"`
	Completed         float64 `json:"completed"`
	Deaths            float64 `json:"deaths"`
	Kills             float64 `json:"kills"`
	OpponentsDefeated float64 `json:"opponentsDefeated"`
	Efficiency        float64 `json:"efficiency"`
	KD                float64 `json:"killsDeathsRatio"`
	KDA               float64 `json:"killsDeathsAssists"`
	Score             float64 `json:"score"`
	ActivityDuration  float64 `json:"activityDurationSeconds"`
	CompletionReason  float64 `json:"completionReason"`
	FireteamID        float64 `json:"fireteamId"`
	StartSeconds      float64 `json:"startSeconds"`
	TimePlayedSeconds float64 `json:"timePlayedSeconds"`
	PlayerCount       float64 `json:"playerCount"`
	TeamScore         float64 `json:"teamScore"`
}

type Extended struct {
	Values ExtendedValues `json:"values"`
}

type ExtendedValues struct {
	PrecisionKills     float64 `json:"precisionKills"`
	WeaponKillsGrenade float64 `json:"weaponKillsGrenade"`
	WeaponKillsMelee   float64 `json:"weaponKillsMelee"`
	WeaponKillsSuper   float64 `json:"weaponKillsSuper"`
	WeaponKillsAbility float64 `json:"weaponKillsAbility"`
}

type Team struct {
	// empty array in sample JSON; add fields if needed later
}

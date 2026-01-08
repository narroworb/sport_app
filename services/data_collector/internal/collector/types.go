package collector

type UnactualTournamentsAndTours struct {
	LeagueName string
	Season     string
	Tours      []uint16
}

type ShortTypeMatch struct {
	HomeTeamName  string
	AwayTeamName  string
	HomeScore     int16
	AwayScore     int16
	Status        string
	DateTimestamp int64
}

type ManagersOfMatch struct {
	HomeTeamManagerID uint32
	AwayTeamManagerID uint32
}

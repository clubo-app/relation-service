package datastruct

type FavoritePartyCount struct {
	PartyId            string `db:"party_id"`
	FavoritePartyCount int64  `db:"favorite_party_count"`
}

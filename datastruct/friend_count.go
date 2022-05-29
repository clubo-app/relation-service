package datastruct

type FriendCount struct {
	UserId      string `db:"user_id"`
	FriendCount int64  `db:"friend_count"`
}

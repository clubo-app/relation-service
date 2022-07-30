CREATE TABLE IF NOT EXISTS party_participants (
    user_id text,
    party_id text,
    joined_at timestamp,
    PRIMARY KEY (party_id, user_id)
) WITH cdc = {'enabled':true};

CREATE INDEX ON party_participants(user_id);

CREATE TABLE IF NOT EXISTS party_invites (
  user_id text,
  inviter_id text,
  party_id text,
  PRIMARY KEY (user_id, party_id)
) WITH cdc = {'enabled':true};

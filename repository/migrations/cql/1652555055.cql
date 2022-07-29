CREATE TABLE IF NOT EXISTS favorite_parties (
    user_id text,
    party_id text,
    favorited_at timestamp,
    PRIMARY KEY (party_id, user_id)
)WITH cdc = {'enabled':true};

CREATE MATERIALIZED VIEW IF NOT EXISTS favorite_parties_by_user AS
    SELECT * FROM favorite_parties 
    WHERE party_id IS NOT NULL AND user_id IS NOT NULL AND favorited_at IS NOT NULL
    PRIMARY KEY (user_id, favorited_at, party_id)
    WITH CLUSTERING ORDER BY (favorited_at DESC, party_id ASC);

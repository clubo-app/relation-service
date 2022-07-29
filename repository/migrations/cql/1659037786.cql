CREATE TABLE IF NOT EXISTS party_participants (
    user_id text,
    party_id text,
    joined boolean,
    joined_at timestamp,
    PRIMARY KEY (party_id, user_id)
) WITH cdc = {'enabled':true};

CREATE MATERIALIZED VIEW IF NOT EXISTS party_participants_by_user AS
    SELECT * FROM party_participants 
    WHERE party_id IS NOT NULL AND user_id IS NOT NULL AND joined IS NOT NULL AND joined_at IS NOT NULL
    PRIMARY KEY ((user_id, joined), joined_at, party_id)
    WITH CLUSTERING ORDER BY (joined_at DESC, party_id ASC);

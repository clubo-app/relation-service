CREATE TABLE IF NOT EXISTS stories_by_user (
    story_id text,
    user_id text,
    party_id text,
    url text,
    tagged_friends list<text>,
    PRIMARY KEY (user_id, story_id) 
)WITH CLUSTERING ORDER BY (story_id DESC);

CREATE MATERIALIZED VIEW IF NOT EXISTS stories_by_party AS
    SELECT * FROM stories_by_user 
    WHERE party_id IS NOT NULL AND user_id IS NOT NULL AND story_id IS NOT NULL
    PRIMARY KEY (party_id, story_id, user_id)
    WITH CLUSTERING ORDER BY (story_id DESC, user_id ASC);

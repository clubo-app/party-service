CREATE TABLE parties (
    id varchar(27) PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    location geometry(POINT, 4326) NOT NULL,
    street_address TEXT,
    postal_code TEXT,
    state TEXT,
    country TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP
);

CREATE INDEX parties_by_user_id_idx ON parties (user_id, id);

CREATE INDEX party_location_idx
  ON parties
  USING GIST (location);

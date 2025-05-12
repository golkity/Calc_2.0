CREATE TABLE IF NOT EXISTS expressions (
    id          SERIAL           PRIMARY KEY,
    user_id     INTEGER          NOT NULL
                                 REFERENCES users(id)
                                 ON DELETE CASCADE,
    raw         TEXT             NOT NULL, 
    result      DOUBLE PRECISION,
    status      TEXT             NOT NULL,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ      NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tasks (
    id          SERIAL           PRIMARY KEY,
    user_id     INTEGER          NOT NULL
                                 REFERENCES users(id)
                                 ON DELETE CASCADE,
    title       TEXT             NOT NULL,
    status      TEXT             NOT NULL,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ      NOT NULL DEFAULT now()
);
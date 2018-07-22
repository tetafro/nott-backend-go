BEGIN;

CREATE TABLE "user" (
    id         SERIAL,
    email      VARCHAR NOT NULL,
    password   VARCHAR NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE "token" (
    id      SERIAL,
    user_id INTEGER NOT NULL,
    string  VARCHAR NOT NULL,
    ttl     INTEGER NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE TABLE "folder" (
    id        SERIAL,
    user_id   INTEGER NOT NULL,
    parent_id INTEGER,
    title     VARCHAR NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" (id),
    FOREIGN KEY (parent_id) REFERENCES "folder" (id) ON DELETE CASCADE
);

CREATE TABLE "notepad" (
    id        SERIAL,
    user_id   INTEGER NOT NULL,
    folder_id INTEGER NOT NULL,
    title     VARCHAR NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" (id),
    FOREIGN KEY (folder_id) REFERENCES "folder" (id) ON DELETE CASCADE
);

CREATE TABLE "note" (
    id         SERIAL,
    user_id    INTEGER NOT NULL,
    notepad_id INTEGER NOT NULL,
    title      VARCHAR NOT NULL,
    text       VARCHAR NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" (id),
    FOREIGN KEY (notepad_id) REFERENCES "notepad" (id) ON DELETE CASCADE
);

COMMIT;

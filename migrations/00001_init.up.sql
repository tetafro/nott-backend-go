BEGIN;

CREATE TABLE "user" (
    id         SERIAL,
    email      VARCHAR NOT NULL,
    password   VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE "folder" (
    id         SERIAL,
    user_id    INTEGER NOT NULL,
    parent_id  INTEGER,
    title      VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" (id),
    FOREIGN KEY (parent_id) REFERENCES "folder" (id) ON DELETE CASCADE
);

CREATE TABLE "notepad" (
    id         SERIAL,
    user_id    INTEGER NOT NULL,
    folder_id  INTEGER NOT NULL,
    title      VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
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
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" (id),
    FOREIGN KEY (notepad_id) REFERENCES "notepad" (id) ON DELETE CASCADE
);

INSERT INTO "user" (id, email, password, created_at)
VALUES (1, 'bob@example.com', '$2a$14$u5zeH6lOZmOg64iZwpUYc.pyBi9LlGBOs5FgTf9lpi7NM4.jTW7OS', NOW());

COMMIT;

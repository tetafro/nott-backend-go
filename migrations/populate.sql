BEGIN;

INSERT INTO "user" (id, email, password)
VALUES (1, 'bob@example.com', 'qwerty');
INSERT INTO "user" (id, email, password)
VALUES (2, 'alice@example.com', 'ytrewq');

INSERT INTO folder (id, user_id, parent_id, title)
VALUES (1, 1, NULL, 'Folder 1');
INSERT INTO folder (id, user_id, parent_id, title)
VALUES (2, 1, NULL, 'Folder 2');
INSERT INTO folder (id, user_id, parent_id, title)
VALUES (3, 1, NULL, 'Folder 3');

INSERT INTO notepad (id, user_id, folder_id, title)
VALUES (1, 1, 1, 'Notepad 1');
INSERT INTO notepad (id, user_id, folder_id, title)
VALUES (2, 1, 1, 'Notepad 2');
INSERT INTO notepad (id, user_id, folder_id, title)
VALUES (3, 1, 1, 'Notepad 3');

INSERT INTO note (id, user_id, notepad_id, title, text)
VALUES (1, 1, 1, 'Note 1', 'Text');
INSERT INTO note (id, user_id, notepad_id, title, text)
VALUES (2, 1, 1, 'Note 2', 'Text');
INSERT INTO note (id, user_id, notepad_id, title, text)
VALUES (3, 1, 1, 'Note 3', 'Text');

COMMIT;

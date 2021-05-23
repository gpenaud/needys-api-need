/* TABLE NEED */
CREATE TABLE need (id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, name VARCHAR(100), priority VARCHAR(100)) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO need (name, priority) VALUES ('sécurité', 'very high');
INSERT INTO need (name, priority) VALUES ('partage', 'haut');
INSERT INTO need (name, priority) VALUES ('accomplissement personnel', 'bas');

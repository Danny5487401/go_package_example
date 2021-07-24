rm -f example.db
sqlite3 example.db 'CREATE TABLE people(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);'
sqlite3 example.db 'INSERT INTO people (name, age) VALUES ("drew", 35);'
sqlite3 example.db 'INSERT INTO people (name, age) VALUES ("jane", 29);'

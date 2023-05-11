CREATE TABLE IF NOT EXISTS users
(
    id integer not null primary key autoincrement,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE TABLE IF NOT EXISTS post
(
    id integer primary key autoincrement,
    title varchar(100),
    userid integer not null,
    username VARCHAR(50),
    description TEXT not null,
    created_at datetime,
    Foreign key(userid) REFERENCES users (id),
    Foreign key (username) references users (username)
);

CREATE TABLE IF NOT EXISTS category(
    id integer primary key autoincrement,
    type varchar(100)
);

INSERT INTO category (type) VALUES ('General topic'), ('Food'), ('Life style'), ('Sport'), ('Fashion');

CREATE TABLE IF NOT EXISTS relation (
    id integer primary key autoincrement,
    post_id integer not null,
    cat_id integer not null,
    Foreign key(post_id) REFERENCES post (id),
    Foreign key(cat_id) references category (id)
);



CREATE TABLE IF NOT EXISTS session
(
    id integer primary key autoincrement,
    user_id integer not null,
    user_name text not null,
    token text not null
);
CREATE TABLE IF NOT EXISTS comments
(
    id integer primary key autoincrement,
    author text,
    userid integer not null,
    postid integer not null,
    content text,
    FOREIGN KEY (userid) REFERENCES users(id),
    FOREIGN KEY (postid) REFERENCES post(id) ON DELETE CASCADE
);
CREATE TABLE if NOT EXISTS reaction
(
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
    userid integer not null,
    postid integer not null,
    like integer,
    dislike integer
);
CREATE TABLE if NOT EXISTS commentreaction
(
    id integer primary key autoincrement,
    userid integer not null,
    commentid integer not null,
    like integer not null,
    dislike integer not null
);
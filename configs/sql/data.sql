insert into users (nickname , fullname, about , email)
values
('username', 'name surname', 'motherfucker', 'user@mail.ru'),
('username2', 'name surname', 'motherfucker2', 'user2@mail.ru');

insert into forums (id, title, nickname, slug, posts, threads)
values
(1, 'forum-title', 'username', 'superslug', 0, 0);

insert into threads (id, title, author, forum_title, forum_id, message)
values
(1, 'thread-title', 'username2', 'forum-title', 1, 'message');

insert into posts (id, path, parent_id, author, message, forum_id, thread_id)
values
(1, '{1}', 0, 'username', 'my_message', 1, 1);
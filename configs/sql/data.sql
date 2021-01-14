insert into users (nickname , fullname, about , email)
values
('username', 'name surname', 'motherfucker', 'user@mail.ru'),
('username2', 'name surname', 'motherfucker2', 'user2@mail.ru');



insert into forums (id, title, nickname, slug, posts, threads)
values
(1, 'forum-title', 'username', 'superslug', 0, 0);



insert into threads (id, title, author, forum_slug, forum_id, message, slug)
values
(1, 'thread-title', 'username2', 'superslug', 1, 'message', 'slug1231243234');
insert into threads (id, title, author, forum_slug, forum_id, message, slug)
values
(2, 'thread-title2', 'username', 'superslug', 1, 'message2', 'slugurhouweiur');


insert into posts (id, path, parent_id, author, message, forum_id, thread_id)
values
(1, '{1}', 0, 'username', 'my_message', 1, 1);

insert into posts (id, path, parent_id, author, message, forum_id, thread_id)
values
(2, '{1, 2}', 0, 'username', 'my_message1.2', 1, 1);

insert into posts (id, path, parent_id, author, message, forum_id, thread_id)
values
(3, '{1, 3}', 0, 'username', 'my_message1.3', 1, 1);

insert into posts (id, path, parent_id, author, message, forum_id, thread_id)
values
(4, '{1, 2, 4}', 0, 'username', 'my_message1.2.4', 1, 1);

insert into posts (id, path, parent_id, author, message, forum_id, thread_id)
values
(5, '{1, 2, 5}', 0, 'username', 'my_message1.2.5', 1, 1);

insert into posts (id, path, parent_id, author, message, forum_id, thread_id)
values
(6, '{6}', 0, 'username', 'my_message6', 1, 1);

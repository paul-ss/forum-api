
ANALYZE;

explain ANALYZE
WITH f AS (SELECT title, id FROM forums WHERE slug = 'YQ0yPFo6H9Kf8')
        SELECT u.nickname, u.fullname, u.about, u.email
        FROM
        (SELECT DISTINCT author FROM threads WHERE forum_id = (SELECT id from f)
        UNION
        SELECT DISTINCT author FROM posts WHERE forum_id = (SELECT id from f)) AS a
        JOIN users u ON a.author = u.nickname
        WHERE nickname > 'asd'
        ORDER BY nickname desc
        LIMIT 3000;

explain ANALYZE
WITH f AS (SELECT title, id FROM forums WHERE slug = 'YQ0yPFo6H9Kf8')
SELECT u.nickname, u.fullname, u.about, u.email
FROM
    (SELECT author FROM threads WHERE forum_id = (SELECT id from f) AND author > 'asd'
     UNION
     SELECT author FROM posts WHERE forum_id = (SELECT id from f) AND author > 'asd') AS a
        JOIN users u ON a.author = u.nickname
    GROUP BY nickname
ORDER BY nickname desc
LIMIT 3000;
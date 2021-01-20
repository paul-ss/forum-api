
ANALYZE ;

WITH t AS (SELECT id FROM threads WHERE id = 6953)
		SELECT id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created,
			CASE WHEN (SELECT id FROM t) IS NOT NULL THEN 1 ELSE 0 END as exists
			FROM posts p
			WHERE thread_id = (SELECT id FROM t)
			    AND id  >  5000
			ORDER BY id
			;

WITH t AS (SELECT id FROM threads WHERE id = 450053)
SELECT id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created,
       CASE WHEN (SELECT id FROM t) IS NOT NULL THEN 1 ELSE 0 END as exists
FROM posts p
WHERE thread_id = (SELECT id FROM t)
    AND p.path > COALESCE((SELECT pp.path FROM posts pp WHERE id = 450053), ARRAY[0])
ORDER BY p.path
LIMIT 2000000;



WITH t AS (SELECT id FROM threads WHERE id = 6953),
     ls AS
    (WITH edge AS (SELECT COALESCE((SELECT pip.path FROM posts pip WHERE id = 450053), ARRAY[0]) AS p)
    SELECT min(id) AS min, max(id) AS max
    FROM (SELECT id FROM posts pi
    WHERE thread_id = (SELECT id FROM t) AND array_length(path, 1) = 1
    AND pi.path  > (SELECT p FROM edge)
    AND pi.path[1] > (SELECT p[1] FROM edge)
    ORDER BY pi.path[1] , pi.path
    LIMIT 10000) as unused_alias)

SELECT id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created,
       CASE WHEN (SELECT id FROM t) IS NOT NULL THEN 1 ELSE 0 END as exists
FROM posts p
WHERE thread_id = (SELECT id FROM t)
    AND path[1] >= (SELECT min FROM ls)  AND path[1] <= (SELECT max FROM ls)
ORDER BY p.path[1], p.path

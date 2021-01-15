DROP TRIGGER IF EXISTS vote_threads ON votes;
DROP FUNCTION IF EXISTS on_vote_threads();

DROP FUNCTION IF EXISTS on_threads_ins_upd();
DROP TRIGGER IF EXISTS threads_ins_upd ON threads;

DROP FUNCTION IF EXISTS on_posts_ins_upd();
DROP TRIGGER IF EXISTS posts_ins_upd ON posts;

DROP SEQUENCE pidseq;
ALTER TABLE likes
ADD CONSTRAINT unique_blog_user UNIQUE (blog_id, user_id);
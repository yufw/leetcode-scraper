drop table if exists problems cascade;
create table problems (
       leetcode_id integer primary key,
       title text not null,
       title_slug text not null,
       content text not null,
       difficulty text not null,
       likes integer not null,
       dislikes integer not null,
       total_accepted bigint not null,
       total_submission bigint not null,
       similar_questions text,
       hints text[]
);

drop table if exists topics cascade;
create table topics (
       slug text primary key,
       name text not null
);

drop table if exists problem_topic;
create table problem_topic (
       problem_id integer not null references problems (leetcode_id) on delete cascade,
       topic_slug text not null references topics (slug) on delete cascade,
       constraint problem_topic_key primary key (problem_id, topic_slug)
);

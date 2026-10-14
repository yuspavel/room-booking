-- +goose Up
WITH 
user_inserted AS(
INSERT INTO users (name,pass_hash,email,phone)
VALUES ('Artur','f1c3a7f2360d86d5543176893de6a389','artur@gmail.com','+78901256742'),
       ('Gena','2b47ddb6c5c943d956820159b7888d1a','gena@yahho.com','89001730947'),
       ('Ostin','de5e1b5ec6032ffd033dfce170bf8901','ostin.pauers@mail.ru','007812765974'),
       ('Andrew','f1707436e2eeb1fc97da8b07a6468a73','andre@gmail.com','+71209002828')
RETURNING id,name
),
rooms_inserted AS(
INSERT INTO room (name,capacity, hour_price,currency)
VALUES('Green meeting office',40,100.5,'RUB')
RETURNING id,name
)

INSERT INTO booking(room_id,user_id,title,description,participants_count,from_time,to_time,status,paid)
SELECT r.id,u.id,'Обсуждение проекта строительства','Составление расписания, определение команд и зон ответственности',(random()*30+1)::int,'2026-01-10T10:00:00+05:00'::TIMESTAMPTZ,'2026-01-10T15:00:00+05:00' TIMESTAMPTZ,'pending',false
FROM user_inserted u, rooms_inserted r;

-- +goose Down
DELETE FROM users WHERE name IN ('Artur','Gena','Ostin','Andrew');
DELETE FROM room WHERE name = 'Green meeting office';
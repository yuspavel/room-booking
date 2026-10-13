-- +goose Up
ALTER TABLE booking ADD CONSTRAINT fk_booking_room FOREIGN KEY (room_id) REFERENCES room(id) ON DELETE CASCADE;

ALTER TABLE booking ADD CONSTRAINT fk_booking_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE booking DROP CONSTRAINT fk_booking_user;

ALTER TABLE booking DROP CONSTRAINT fk_booking_room;
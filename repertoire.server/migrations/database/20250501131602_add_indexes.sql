-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_albums_user_id ON albums(user_id);
CREATE INDEX idx_artists_user_id ON artists(user_id);
CREATE INDEX idx_band_member_roles_user_id ON band_member_roles(user_id);
CREATE INDEX idx_guitar_tunings_user_id ON guitar_tunings(user_id);
CREATE INDEX idx_instruments_user_id ON instruments(user_id);
CREATE INDEX idx_playlists_user_id ON playlists(user_id);
CREATE INDEX idx_song_section_types_user_id ON song_section_types(user_id);
CREATE INDEX idx_songs_user_id ON songs(user_id);

CREATE INDEX idx_band_members_artist_id ON band_members(artist_id);

CREATE INDEX idx_song_sections_song_id ON song_sections(song_id);

CREATE INDEX idx_song_section_histories_section_id ON song_section_histories(song_section_id);

CREATE INDEX idx_users_email ON users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_albums_user_id;
DROP INDEX idx_artists_user_id;
DROP INDEX idx_band_member_roles_user_id;
DROP INDEX idx_guitar_tunings_user_id;
DROP INDEX idx_instruments_user_id;
DROP INDEX idx_playlists_user_id;
DROP INDEX idx_song_section_types_user_id;
DROP INDEX idx_songs_user_id;

DROP INDEX idx_band_members_artist_id;

DROP INDEX idx_song_sections_song_id;

DROP INDEX idx_song_section_histories_section_id;

DROP INDEX idx_users_email;
-- +goose StatementEnd

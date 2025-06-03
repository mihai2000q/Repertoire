-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.playlist_songs
    ADD COLUMN id uuid not null default gen_random_uuid();

ALTER TABLE public.playlist_songs
    DROP CONSTRAINT playlist_songs_pkey;

ALTER TABLE public.playlist_songs
    ADD PRIMARY KEY (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.playlist_songs
    DROP CONSTRAINT playlist_songs_pkey;

ALTER TABLE public.playlist_songs
    ADD PRIMARY KEY (song_id, playlist_id);

ALTER TABLE public.playlist_songs
    DROP COLUMN id;
-- +goose StatementEnd

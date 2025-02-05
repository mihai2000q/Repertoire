-- +goose Up
-- +goose StatementBegin

ALTER TABLE public.song_sections
    ADD COLUMN occurrences bigint not null default 0;

ALTER TABLE public.song_sections
    ADD COLUMN band_member_id uuid;

ALTER TABLE public.song_sections
    ADD CONSTRAINT fk_band_members_song_sections
    FOREIGN KEY (band_member_id) REFERENCES band_members(id)
    ON DELETE SET NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE public.song_sections
    DROP COLUMN occurrences;

ALTER TABLE public.song_sections
    DROP COLUMN band_member_id;

ALTER TABLE public.song_sections
    DROP CONSTRAINT fk_band_members_sections;

-- +goose StatementEnd

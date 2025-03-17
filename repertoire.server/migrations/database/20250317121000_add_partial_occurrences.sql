-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.song_sections
    ADD COLUMN partial_occurrences bigint not null default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.song_sections
    DROP COLUMN partial_occurrences;
-- +goose StatementEnd

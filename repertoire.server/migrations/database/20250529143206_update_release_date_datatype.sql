-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.albums
    ALTER COLUMN release_date TYPE date
        USING (release_date AT TIME ZONE 'CEST')::date;

ALTER TABLE public.songs
    ALTER COLUMN release_date TYPE date
        USING (release_date AT TIME ZONE 'CEST')::date;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.albums
    ALTER COLUMN release_date TYPE timestamp with time zone;

ALTER TABLE public.songs
    ALTER COLUMN release_date TYPE timestamp with time zone;
-- +goose StatementEnd

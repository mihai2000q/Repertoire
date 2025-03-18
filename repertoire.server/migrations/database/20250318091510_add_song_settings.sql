-- +goose Up
-- +goose StatementBegin
-- FIX: Wrong naming for constraint
ALTER TABLE public.band_member_roles
    RENAME CONSTRAINT fk_users_song_section_types to fk_users_band_member_roles;

CREATE TABLE public.song_settings
(
    id             uuid not null primary key,
    instrument_id  uuid constraint fk_instruments_song_settings references public.instruments,
    band_member_id uuid constraint fk_band_members_song_settings references public.band_members,
    song_id        uuid not null constraint fk_songs_song_settings references public.songs
);

DO
$$
    DECLARE
        current_song_id uuid;
    BEGIN
        -- Loop through users
        FOR current_song_id IN (SELECT id FROM public.songs)
            LOOP
                INSERT INTO public.song_settings
                VALUES (gen_random_uuid(), null, null, current_song_id);
            END LOOP;
    END
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.band_member_roles
    RENAME CONSTRAINT fk_users_band_member_roles to fk_users_song_section_types;

DROP TABLE public.song_settings;
-- +goose StatementEnd

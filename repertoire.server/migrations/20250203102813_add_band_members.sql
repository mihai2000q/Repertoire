-- +goose Up
-- +goose StatementBegin

ALTER TABLE public.artists
    ADD is_band bool not null default false;

CREATE TABLE public.band_member_roles
(
    id      uuid   not null primary key,
    name    varchar(24),
    "order" bigint not null,
    user_id uuid   not null constraint fk_users_song_section_types references public.users
);

CREATE TABLE public.band_members
(
    id          uuid   not null primary key,
    name        varchar(100),
    "order"     bigint not null,
    color       varchar(7),
    image_url   text,
    artist_id   uuid not null constraint fk_artists_band_members references public.artists on delete cascade,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP not null
);

CREATE TABLE public.band_member_has_roles(
      band_member_id      uuid not null constraint fk_band_members_band_member_has_roles references public.band_members on delete cascade,
      band_member_role_id uuid not null constraint fk_band_member_roles_band_member_has_roles references public.band_member_roles on delete cascade,
      primary key(band_member_id, band_member_role_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE public.artists
    DROP COLUMN is_band;

DROP TABLE public.band_member_has_roles;
DROP TABLE public.band_members;
DROP TABLE public.band_member_roles;

-- +goose StatementEnd

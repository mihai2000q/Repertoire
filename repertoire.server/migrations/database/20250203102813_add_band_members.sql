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

-- define default Band Member Roles
WITH band_member_roles AS (
    SELECT unnest(ARRAY['Vocalist', 'Lead Guitarist', 'Rhythm Guitarist', 'Bassist', 'Drummer', 'Pianist']) AS name,
    generate_series(0, 5) AS name_order
),

-- generate a sequence of order indices for each user band member role together with the name
user_band_member_roles AS (
    SELECT
        users.id AS user_id,
        band_member_roles.name,
        band_member_roles.name_order AS order_index
    FROM
        users
    CROSS JOIN
        band_member_roles
)

-- Insert new user band member roles
INSERT INTO public.band_member_roles (id, name, "order", user_id)
SELECT
    gen_random_uuid() AS id,
    ub.name,
    ub.order_index,
    ub.user_id
FROM
    user_band_member_roles ub;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE public.artists
    DROP COLUMN is_band;

DROP TABLE public.band_member_has_roles;
DROP TABLE public.band_members;
DROP TABLE public.band_member_roles;

-- +goose StatementEnd

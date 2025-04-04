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

CREATE TABLE public.instruments(
   id uuid not null primary key,
   name varchar(30),
   "order" bigint not null,
   user_id uuid not null constraint fk_users_instruments references public.users
);

ALTER TABLE public.song_sections
    ADD COLUMN instrument_id uuid;

ALTER TABLE public.song_sections
    ADD CONSTRAINT fk_instruments_song_sections
    FOREIGN KEY (instrument_id) REFERENCES instruments(id)
    ON DELETE SET NULL;

-- Insert Instruments per user

-- define default Instruments' names
WITH instrument_names AS (
    SELECT unnest(ARRAY['Voice', 'Piano', 'Keyboard', 'Drums', 'Electric Guitar',
        'Acoustic Guitar', 'Bass', 'Ukulele', 'Violin', 'Saxophone', 'Flute', 'Harp'
    ]) AS name,
    generate_series(0, 11) AS name_order
),

-- generate a sequence of order indices for each user instrument together with the name
user_instruments AS (
    SELECT
        users.id AS user_id,
        instrument_names.name,
        instrument_names.name_order AS order_index
    FROM
        users
    CROSS JOIN
        instrument_names
)

-- Insert new user instruments
INSERT INTO public.instruments (id, name, "order", user_id)
SELECT
    gen_random_uuid() AS id,
    ui.name,
    ui.order_index,
    ui.user_id
FROM
    user_instruments ui;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE public.song_sections
    DROP COLUMN occurrences;

ALTER TABLE public.song_sections
    DROP COLUMN band_member_id;

ALTER TABLE public.song_sections
    DROP CONSTRAINT fk_band_members_song_sections;

ALTER TABLE public.song_sections
    DROP COLUMN instrument_id;

ALTER TABLE public.song_sections
    DROP CONSTRAINT fk_instruments_song_sections;

DROP TABLE public.instruments;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

-- to preserve the order of types, we decided to remove them all and add them again
DELETE FROM public.song_section_types WHERE true;

-- define new Song Section Types' names
WITH section_types AS (
    SELECT unnest(ARRAY[
        'Intro', 'Verse', 'Pre-Chorus', 'Chorus', 'Interlude',
        'Bridge', 'Breakdown', 'Solo', 'Riff', 'Outro'
    ]) AS name,
    generate_series(0, 9) AS name_order
),

-- generate a sequence of order indices for each user instrument together with the name
user_section_types AS (
    SELECT
        users.id AS user_id,
        section_types.name,
        section_types.name_order AS order_index
    FROM
        users
    CROSS JOIN
        section_types
)

-- Insert new user instruments
INSERT INTO public.song_section_types (id, name, "order", user_id)
SELECT
    gen_random_uuid() AS id,
    us.name,
    us.order_index,
    us.user_id
FROM
    user_section_types us;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- to preserve the order of types, remove them all and add them again
DELETE FROM public.song_section_types WHERE true;

-- define old Song Section Types' names
WITH section_types AS (
    SELECT unnest(ARRAY[
        'Intro', 'Verse', 'Chorus', 'Interlude', 'Breakdown', 'Solo', 'Riff', 'Outro'
    ]) AS name,
    generate_series(0, 9) AS name_order
),

-- generate a sequence of order indices for each user song section type together with the name
     user_section_types AS (
         SELECT
             users.id AS user_id,
             section_types.name,
             section_types.name_order AS order_index
         FROM
             users
                 CROSS JOIN
             section_types
     )

-- Insert new user song section types
INSERT INTO public.song_section_types (id, name, "order", user_id)
SELECT
    gen_random_uuid() AS id,
    us.name,
    us.order_index,
    us.user_id
FROM
    user_section_types us;

-- +goose StatementEnd

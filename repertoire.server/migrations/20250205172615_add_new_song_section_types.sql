-- +goose Up
-- +goose StatementBegin

-- insert new types and reorder the old ones
DO
$$
    DECLARE
        current_user_id uuid;
    BEGIN
        -- Loop through users
        FOR current_user_id IN (SELECT id FROM public.users)
            LOOP
                INSERT INTO public.song_section_types
                VALUES (gen_random_uuid(), 'Pre-Chorus', 2, current_user_id),
                       (gen_random_uuid(), 'Bridge', 6, current_user_id);
            END LOOP;
    END
$$;

UPDATE public.song_section_types
SET "order" = 3
WHERE name = 'Chorus';

UPDATE public.song_section_types
SET "order" = 4
WHERE name = 'Interlude';

UPDATE public.song_section_types
SET "order" = 5
WHERE name = 'Breakdown';

UPDATE public.song_section_types
SET "order" = 7
WHERE name = 'Solo';

UPDATE public.song_section_types
SET "order" = 8
WHERE name = 'Riff';

UPDATE public.song_section_types
SET "order" = 9
WHERE name = 'Outro';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- remove new types and reorder
DELETE
FROM public.song_section_types
WHERE name = 'Pre-Chorus'
   OR name = 'Bridge';

UPDATE public.song_section_types
SET "order" = 2
WHERE name = 'Chorus';

UPDATE public.song_section_types
SET "order" = 3
WHERE name = 'Interlude';

UPDATE public.song_section_types
SET "order" = 4
WHERE name = 'Breakdown';

UPDATE public.song_section_types
SET "order" = 5
WHERE name = 'Solo';

UPDATE public.song_section_types
SET "order" = 6
WHERE name = 'Riff';

UPDATE public.song_section_types
SET "order" = 7
WHERE name = 'Outro';

-- +goose StatementEnd

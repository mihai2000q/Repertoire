-- +goose Up
-- +goose StatementBegin
DO
$$
    DECLARE
        current_user_id uuid;
    BEGIN
        -- Loop through users
        FOR current_user_id IN (SELECT id FROM public.users)
            LOOP
                INSERT INTO public.song_section_types
                VALUES (gen_random_uuid(), 'Keyboardist', 6, current_user_id),
                       (gen_random_uuid(), 'Backing Vocalist', 7, current_user_id);
            END LOOP;
    END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM public.song_section_types
WHERE name = 'Keyboardist'
   OR name = 'Backing Vocalist';
-- +goose StatementEnd

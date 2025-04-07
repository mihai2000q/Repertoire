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
                INSERT INTO public.band_member_roles
                VALUES (gen_random_uuid(), 'Keyboardist', 6, current_user_id),
                       (gen_random_uuid(), 'Backing Vocalist', 7, current_user_id);
            END LOOP;
    END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM public.band_member_roles
WHERE name = 'Keyboardist'
   OR name = 'Backing Vocalist';
-- +goose StatementEnd

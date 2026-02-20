-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.song_arrangements
(
    id      uuid        not null primary key,
    name    varchar(30) not null,
    "order" bigint      not null,
    song_id uuid        not null
        constraint fk_songs_song_arrangements references public.songs on delete cascade
);
CREATE TABLE public.song_section_occurrences
(
    occurrences    bigint not null,
    arrangement_id uuid   not null
        constraint fk_song_arrangements_song_section_occurrences references public.song_arrangements on delete cascade,
    section_id     uuid   not null
        constraint fk_song_sections_song_section_occurrences references public.song_sections on delete cascade,
    primary key (arrangement_id, section_id)
);

ALTER TABLE public.songs
    ADD COLUMN default_arrangement_id uuid,
    ADD CONSTRAINT fk_song_arrangements_songs
        foreign key (default_arrangement_id) references song_arrangements (id) on delete set null;

-- Migrate song sections occurrences and partial occurrences to song arrangements and set default arrangement
DO
$$
    DECLARE
        perfect_arrangement_id      uuid;
        partial_arrangement_id      uuid;
        current_occurrences         bigint;
        current_partial_occurrences bigint;
        current_song_id             uuid;
        current_song_section_id     uuid;
        sections_count              bigint;
        partial_sections_count      bigint;
        skip_partial_occurrences    bool;
    BEGIN
        FOR current_song_id IN (SELECT id FROM public.songs)
            LOOP
            -- if all song sections' have 0 partial occurrences,
            -- then the user must have not set them,
            -- so there is no point to migrate them to a partial rehearsal arrangement
                SELECT COUNT(*)
                INTO sections_count
                FROM public.song_sections
                WHERE song_id = current_song_id;
                SELECT COUNT(*)
                INTO partial_sections_count
                FROM public.song_sections
                WHERE song_id = current_song_id
                  AND partial_occurrences = 0;

                -- Insert perfect rehearsal arrangements
                perfect_arrangement_id = gen_random_uuid();
                INSERT INTO public.song_arrangements
                VALUES (perfect_arrangement_id, 'Perfect Rehearsal', 0, current_song_id);
                -- set it as default arrangement
                UPDATE public.songs SET default_arrangement_id = perfect_arrangement_id WHERE id = current_song_id;

                -- Insert partial rehearsal arrangements, if there is any partial occurrence within the sections
                IF sections_count = partial_sections_count THEN
                    skip_partial_occurrences = true;
                ELSE
                    skip_partial_occurrences = false;
                    partial_arrangement_id = gen_random_uuid();
                    INSERT INTO public.song_arrangements
                    VALUES (partial_arrangement_id, 'Partial Rehearsal', 1, current_song_id);
                END IF;

                -- Insert the occurrences
                FOR current_occurrences, current_partial_occurrences, current_song_section_id IN
                    (SELECT occurrences, partial_occurrences, id
                     FROM public.song_sections
                     WHERE song_id = current_song_id)
                    LOOP
                        INSERT INTO public.song_section_occurrences
                        VALUES (current_occurrences,
                                perfect_arrangement_id,
                                current_song_section_id);

                        IF NOT skip_partial_occurrences THEN
                            INSERT INTO public.song_section_occurrences
                            VALUES (current_partial_occurrences,
                                    partial_arrangement_id,
                                    current_song_section_id);
                        END IF;
                    END LOOP;
            END LOOP;
    END
$$;

ALTER TABLE public.song_sections
    DROP COLUMN occurrences;
ALTER TABLE public.song_sections
    DROP COLUMN partial_occurrences;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.song_sections
    ADD COLUMN occurrences bigint not null default 0;
ALTER TABLE public.song_sections
    ADD COLUMN partial_occurrences bigint not null default 0;

-- Migrate song arrangements' occurrences to occurrences and partial occurrences of song sections
DO
$$
    DECLARE
        current_occurrences     bigint;
        current_arrangement_id  uuid;
        current_song_section_id uuid;
    BEGIN
        -- Take the occurrences from the arrangements titled "perfect rehearsal"
        FOR current_arrangement_id IN (SELECT id
                                       FROM public.song_arrangements
                                       WHERE name = 'Perfect Rehearsal')
            LOOP
                -- Alter song sections' occurrences
                FOR current_occurrences, current_song_section_id IN (SELECT occurrences, section_id
                                                                     FROM public.song_section_occurrences
                                                                     WHERE arrangement_id = current_arrangement_id)
                    LOOP
                        UPDATE public.song_sections
                        SET occurrences = current_occurrences
                        WHERE id = current_song_section_id;
                    END LOOP;
            END LOOP;

        -- Take the partial occurrences from the arrangements titled "partial rehearsal"
        FOR current_arrangement_id IN (SELECT id
                                       FROM public.song_arrangements
                                       WHERE name = 'Partial Rehearsal')
            LOOP
                FOR current_occurrences, current_song_section_id IN (SELECT occurrences, section_id
                                                                     FROM public.song_section_occurrences
                                                                     WHERE arrangement_id = current_arrangement_id)
                    LOOP
                        -- Alter song sections' partial occurrences
                        UPDATE public.song_sections
                        SET partial_occurrences = current_occurrences
                        WHERE id = current_song_section_id;
                    END LOOP;
            END LOOP;
    END
$$;

ALTER TABLE public.songs
    DROP CONSTRAINT fk_song_arrangements_songs;
ALTER TABLE public.songs
    DROP COLUMN default_arrangement_id;

DROP TABLE public.song_section_occurrences;
DROP TABLE public.song_arrangements;
-- +goose StatementEnd

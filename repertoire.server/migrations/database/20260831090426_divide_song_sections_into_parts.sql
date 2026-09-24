-- +goose Up
-- +goose StatementBegin
-- Song Parts --
CREATE TABLE public.song_parts
(
    id               uuid                                               not null primary key,
    name             varchar(30)                                        not null,
    song_order       bigint                                             not null,
    rehearsals       bigint                                             not null,
    confidence       bigint                                             not null,
    rehearsals_score bigint                                             not null,
    confidence_score bigint                                             not null,
    progress         bigint                                             not null,
    song_id          uuid                                               not null
        constraint fk_songs_song_parts references public.songs on delete cascade,
    instrument_id    uuid
        constraint fk_instruments_song_parts references public.instruments on delete set null,
    created_at       timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at       timestamp with time zone default CURRENT_TIMESTAMP not null
);
CREATE INDEX idx_song_parts_song_id ON song_parts (song_id);

CREATE TABLE public.song_section_parts
(
    section_id uuid                                               not null
        constraint fk_song_sections_song_section_parts references public.song_sections on delete cascade,
    part_id    uuid                                               not null
        constraint fk_song_parts_song_section_parts references public.song_parts on delete cascade,
    "order"    bigint                                             not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    PRIMARY KEY (part_id, section_id)
);

-- Band members that play a part inside a section (many per section part)
CREATE TABLE public.song_section_part_band_members
(
    part_id        uuid not null,
    section_id     uuid not null,
    band_member_id uuid not null
        constraint fk_band_members_song_section_part_band_members references public.band_members on delete cascade,
    PRIMARY KEY (part_id, section_id, band_member_id),
    CONSTRAINT fk_song_section_parts_song_section_part_band_members
        FOREIGN KEY (part_id, section_id) REFERENCES public.song_section_parts (part_id, section_id)
            ON DELETE CASCADE
);
CREATE INDEX idx_song_section_part_band_members_band_member_id
    ON public.song_section_part_band_members (band_member_id);

-- Rename Section History and Occurrences to Part History and Occurrences
ALTER TABLE public.song_section_histories
    RENAME TO song_part_histories;
ALTER TABLE public.song_part_histories
    RENAME CONSTRAINT song_section_histories_pkey TO song_part_histories_pkey;
ALTER TABLE public.song_part_histories
    ADD COLUMN part_id uuid,
    ADD CONSTRAINT fk_song_parts_song_part_histories
        FOREIGN KEY (part_id) REFERENCES public.song_parts (id)
            ON DELETE CASCADE;
CREATE INDEX idx_song_part_histories_part_id ON public.song_part_histories (part_id);

ALTER TABLE public.song_section_occurrences
    RENAME TO song_part_occurrences;
ALTER TABLE public.song_part_occurrences
    RENAME CONSTRAINT fk_song_arrangements_song_section_occurrences TO fk_song_arrangements_song_part_occurrences;
ALTER TABLE public.song_part_occurrences
    ADD COLUMN part_id uuid,
    ADD CONSTRAINT fk_song_parts_song_part_occurrences
        FOREIGN KEY (part_id) REFERENCES public.song_parts (id)
            ON DELETE CASCADE;
CREATE INDEX idx_song_part_occurrences_part_id ON public.song_part_occurrences (part_id);

-- Migrate Song Sections into Song Parts (with histories and occurrences) --
DO
$$
    DECLARE
        sec             RECORD;
        current_part_id uuid;
        part_order      bigint;
        riff_type_id    uuid;
    BEGIN
        -- Loop through each section
        FOR sec IN
            SELECT ss.*,
                   s.user_id
            FROM public.song_sections ss
                     JOIN public.songs s ON ss.song_id = s.id
            ORDER BY ss.song_id, ss."order"
            LOOP
                -- Get the 'Riff' type ID for this user
                SELECT id
                INTO riff_type_id
                FROM public.song_section_types
                WHERE name = 'Riff'
                  AND user_id = sec.user_id;

                -- Compute next order for this song
                SELECT COALESCE(MAX(song_order), -1) + 1
                INTO part_order
                FROM public.song_parts
                WHERE song_id = sec.song_id;

                -- Insert new part
                INSERT INTO public.song_parts (id, name, song_order, rehearsals, confidence,
                                               rehearsals_score, confidence_score, progress,
                                               song_id, instrument_id,
                                               created_at, updated_at)
                VALUES (gen_random_uuid(), sec.name, part_order,
                        sec.rehearsals, sec.confidence,
                        sec.rehearsals_score, sec.confidence_score, sec.progress,
                        sec.song_id, sec.instrument_id,
                        sec.created_at, sec.updated_at)
                RETURNING id INTO current_part_id;

                -- If the section type is NOT 'Riff', create a SongSectionPart
                -- (and its band member, if the section had one)
                IF sec.song_section_type_id IS DISTINCT FROM riff_type_id THEN
                    INSERT INTO public.song_section_parts (section_id, part_id, "order", created_at)
                    VALUES (sec.id, current_part_id, 0, sec.created_at);

                    IF sec.band_member_id IS NOT NULL THEN
                        INSERT INTO public.song_section_part_band_members (part_id, section_id, band_member_id)
                        VALUES (current_part_id, sec.id, sec.band_member_id);
                    END IF;
                END IF;

                -- Migrate histories for this section to the new part
                UPDATE public.song_part_histories
                SET part_id = current_part_id
                WHERE song_section_id = sec.id;

                -- Migrate occurrences for this section to the new part
                UPDATE public.song_part_occurrences
                SET part_id = current_part_id
                WHERE section_id = sec.id;
            END LOOP;
    END
$$;

ALTER TABLE public.song_part_histories
    DROP COLUMN song_section_id CASCADE,
    ALTER COLUMN part_id SET NOT NULL;
ALTER TABLE public.song_part_occurrences
    DROP COLUMN section_id CASCADE,
    ADD PRIMARY KEY (arrangement_id, part_id);

DELETE
FROM public.song_sections ss
    USING public.song_section_types t
WHERE ss.song_section_type_id = t.id
  AND t.name = 'Riff';

-- Song Section Types --
-- 1.add Post Chorus
DO
$$
    DECLARE
        current_user_id uuid;
    BEGIN
        -- Loop through users
        FOR current_user_id IN (SELECT id FROM public.users)
            LOOP
                INSERT INTO public.song_section_types
                VALUES (gen_random_uuid(), 'Post-Chorus', 4, current_user_id);
            END LOOP;
    END
$$;

-- 2.remove riff
DELETE
FROM public.song_section_types
WHERE name = 'Riff';

-- 3.reorder remaining types
UPDATE public.song_section_types
SET "order" = 5
WHERE name = 'Interlude';
UPDATE public.song_section_types
SET "order" = 6
WHERE name = 'Bridge';
UPDATE public.song_section_types
SET "order" = 7
WHERE name = 'Breakdown';
UPDATE public.song_section_types
SET "order" = 8
WHERE name = 'Solo';

-- Drop old columns
ALTER TABLE public.song_sections
    DROP COLUMN band_member_id,
    DROP COLUMN instrument_id,
    DROP COLUMN rehearsals,
    DROP COLUMN rehearsals_score,
    DROP COLUMN confidence,
    DROP COLUMN confidence_score,
    DROP COLUMN progress;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Restore columns on Song Section
ALTER TABLE public.song_sections
    ADD COLUMN band_member_id   uuid,
    ADD COLUMN instrument_id    uuid,
    ADD COLUMN rehearsals       bigint not null default 0,
    ADD COLUMN rehearsals_score bigint not null default 0,
    ADD COLUMN confidence       bigint not null default 0,
    ADD COLUMN confidence_score bigint not null default 0,
    ADD COLUMN progress         bigint not null default 0,
    ADD CONSTRAINT fk_band_members_song_sections
        FOREIGN KEY (band_member_id) REFERENCES public.band_members (id)
            ON DELETE SET NULL,
    ADD CONSTRAINT fk_instruments_song_sections
        FOREIGN KEY (instrument_id) REFERENCES public.instruments (id)
            ON DELETE SET NULL;

-- Restore names and section fks on History and Occurrences
ALTER TABLE public.song_part_histories
    RENAME TO song_section_histories;
ALTER TABLE public.song_section_histories
    RENAME CONSTRAINT song_part_histories_pkey TO song_section_histories_pkey;

ALTER TABLE public.song_part_occurrences
    RENAME TO song_section_occurrences;
ALTER TABLE public.song_section_occurrences
    RENAME CONSTRAINT fk_song_arrangements_song_part_occurrences TO fk_song_arrangements_song_section_occurrences;

ALTER TABLE public.song_section_histories
    ADD COLUMN song_section_id uuid,
    ADD CONSTRAINT fk_song_sections_history
        FOREIGN KEY (song_section_id) REFERENCES public.song_sections (id)
            ON DELETE CASCADE;
CREATE INDEX idx_song_section_histories_section_id ON public.song_section_histories (song_section_id);
ALTER TABLE public.song_section_occurrences
    ADD COLUMN section_id uuid,
    ADD CONSTRAINT fk_song_sections_song_section_occurrences
        FOREIGN KEY (section_id) REFERENCES public.song_sections (id)
            ON DELETE CASCADE;

-- Song Section Types --
-- 1. add Riff
DO
$$
    DECLARE
        current_user_id uuid;
    BEGIN
        FOR current_user_id IN (SELECT id FROM public.users)
            LOOP
                INSERT INTO public.song_section_types (id, name, "order", user_id)
                VALUES (gen_random_uuid(), 'Riff', 8, current_user_id);
            END LOOP;
    END
$$;

-- 2. before removing post-chorus, reassign to Chorus
DO
$$
    DECLARE
        user_record    RECORD;
        post_chorus_id uuid;
        chorus_id      uuid;
    BEGIN
        FOR user_record IN (SELECT id FROM public.users)
            LOOP
                -- Get the type IDs for this user
                SELECT id
                INTO post_chorus_id
                FROM public.song_section_types
                WHERE user_id = user_record.id
                  AND name = 'Post-Chorus';
                SELECT id
                INTO chorus_id
                FROM public.song_section_types
                WHERE user_id = user_record.id
                  AND name = 'Chorus';

                UPDATE public.song_sections
                SET song_section_type_id = chorus_id
                WHERE song_section_type_id = post_chorus_id;
            END LOOP;
    END
$$;

-- 3. remove post-chorus
DELETE
FROM public.song_section_types
WHERE name = 'Post-Chorus';

-- 4. reorder remaining types
UPDATE public.song_section_types
SET "order" = 4
WHERE name = 'Interlude';
UPDATE public.song_section_types
SET "order" = 5
WHERE name = 'Bridge';
UPDATE public.song_section_types
SET "order" = 6
WHERE name = 'Breakdown';
UPDATE public.song_section_types
SET "order" = 7
WHERE name = 'Solo';

-- 5. Reconstruct sections from parts (using song_section_parts)
-- We'll group by section_id and compute average stats from all parts in that section.
-- Reconstruct sections from parts and migrate histories/occurrences
DO
$$
    DECLARE
        part_record   RECORD;
        sec_id        uuid;
        next_order    bigint;
        found_user_id uuid;
        riff_type_id  uuid;
    BEGIN
        -- Loop over all parts
        FOR part_record IN
            SELECT p.*
            FROM public.song_parts p
            ORDER BY p.song_id, p.song_order
            LOOP
                -- Step 1: Find section via song_section_parts (linked part)
                SELECT ssp.section_id
                INTO sec_id
                FROM public.song_section_parts ssp
                WHERE ssp.part_id = part_record.id
                ORDER BY ssp."order"
                LIMIT 1;

                -- Step 2: If not linked, try to find a section with same name and type 'Riff'
                -- (For estranged parts that were already converted back to sections earlier)
                IF sec_id IS NULL THEN
                    SELECT s.id
                    INTO sec_id
                    FROM public.song_sections s
                             JOIN public.song_section_types t ON t.id = s.song_section_type_id
                    WHERE s.song_id = part_record.song_id
                      AND s.name = part_record.name
                      AND t.name = 'Riff'
                    LIMIT 1;
                END IF;

                -- Step 3: Still NULL? Create a new section for this part (estrangement)
                -- An estranged part has no song_section_parts row, so it has no band member either.
                IF sec_id IS NULL THEN
                    -- Get user_id and riff_type_id for this song
                    SELECT s.user_id
                    INTO found_user_id
                    FROM public.songs s
                    WHERE s.id = part_record.song_id;

                    SELECT id
                    INTO riff_type_id
                    FROM public.song_section_types
                    WHERE name = 'Riff'
                      AND user_id = found_user_id;

                    -- Compute next order for this song
                    SELECT COALESCE(MAX("order"), -1) + 1
                    INTO next_order
                    FROM public.song_sections
                    WHERE song_id = part_record.song_id;

                    -- Insert new section
                    INSERT INTO public.song_sections (id, name, "order", song_id, song_section_type_id,
                                                      band_member_id, instrument_id,
                                                      rehearsals, rehearsals_score, confidence, confidence_score,
                                                      progress,
                                                      created_at, updated_at)
                    VALUES (gen_random_uuid(),
                            part_record.name,
                            next_order,
                            part_record.song_id,
                            riff_type_id,
                            NULL,
                            part_record.instrument_id,
                            part_record.rehearsals,
                            part_record.rehearsals_score,
                            part_record.confidence,
                            part_record.confidence_score,
                            part_record.progress,
                            part_record.created_at,
                            part_record.updated_at)
                    RETURNING id INTO sec_id;
                ELSE
                    -- Section already exists: update its stats and band member from the parts
                    UPDATE public.song_sections s
                    SET rehearsals       = agg.avg_rehearsals,
                        rehearsals_score = agg.avg_rehearsals_score,
                        confidence       = agg.avg_confidence,
                        confidence_score = agg.avg_confidence_score,
                        progress         = agg.avg_progress,
                        band_member_id   = agg.band_member_id
                    FROM (SELECT ssp.section_id,
                                 AVG(p.rehearsals)::bigint       AS avg_rehearsals,
                                 AVG(p.rehearsals_score)::bigint AS avg_rehearsals_score,
                                 AVG(p.confidence)::bigint       AS avg_confidence,
                                 AVG(p.confidence_score)::bigint AS avg_confidence_score,
                                 AVG(p.progress)::bigint         AS avg_progress,
                                 -- sections only hold one band member: pick one from the first part
                                 -- (by "order") that has any
                                 (SELECT bm.band_member_id
                                  FROM public.song_section_part_band_members bm
                                           JOIN public.song_section_parts ssp2
                                                ON ssp2.part_id = bm.part_id
                                                    AND ssp2.section_id = bm.section_id
                                  WHERE bm.section_id = ssp.section_id
                                  ORDER BY ssp2."order", bm.band_member_id
                                  LIMIT 1)                       AS band_member_id
                          FROM public.song_section_parts ssp
                                   JOIN public.song_parts p ON p.id = ssp.part_id
                          WHERE ssp.section_id = sec_id
                          GROUP BY ssp.section_id) agg
                    WHERE s.id = agg.section_id;
                END IF;

                -- Now sec_id is guaranteed to be non-NULL.
                -- Update histories and occurrences for this part
                UPDATE public.song_section_histories
                SET song_section_id = sec_id
                WHERE part_id = part_record.id;

                UPDATE public.song_section_occurrences
                SET section_id = sec_id
                WHERE part_id = part_record.id;
            END LOOP;
    END
$$;

-- Keep a single occurrence per (section_id, arrangement_id): the one from the part with the lowest "order".
-- LEFT JOIN so that sections without a song_section_parts row (estranged parts) are deduplicated too,
-- otherwise the primary key below fails when two parts end up in the same section.
WITH ranked AS (SELECT o.ctid,
                       ROW_NUMBER() OVER (
                           PARTITION BY o.section_id, o.arrangement_id
                           ORDER BY COALESCE(ssp."order", 0), o.part_id
                           ) AS rn
                FROM public.song_section_occurrences o
                         LEFT JOIN public.song_section_parts ssp
                                   ON ssp.part_id = o.part_id
                                       AND ssp.section_id = o.section_id)
DELETE
FROM public.song_section_occurrences
WHERE ctid IN (SELECT ctid FROM ranked WHERE rn > 1);

ALTER TABLE public.song_section_histories
    DROP COLUMN part_id CASCADE,
    ALTER COLUMN song_section_id SET NOT NULL;
ALTER TABLE public.song_section_occurrences
    DROP COLUMN part_id CASCADE,
    ADD PRIMARY KEY (arrangement_id, section_id);

-- The join table references song_section_parts, so it has to go first
DROP TABLE public.song_section_part_band_members;
DROP TABLE public.song_section_parts;
DROP INDEX idx_song_parts_song_id;
DROP TABLE public.song_parts;
-- +goose StatementEnd
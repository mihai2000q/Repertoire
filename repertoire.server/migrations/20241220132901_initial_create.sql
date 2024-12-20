-- +goose Up
-- +goose StatementBegin

create table public.users(
    id uuid not null primary key,
    email varchar(256) not null constraint uni_users_email unique,
    password text not null,
    name varchar(100) not null,
    profile_picture_url text,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);

create table public.artists(
    id uuid not null primary key,
    name varchar(100) not null,
    image_url text,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    user_id uuid not null constraint fk_users_artists references public.users
);

create table public.playlists(
    id uuid not null primary key,
    title varchar(100) not null,
    description text not null,
    image_url text,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    user_id uuid not null constraint fk_users_playlists references public.users
);

create table public.guitar_tunings(
    id uuid not null primary key,
    name varchar(16) not null,
    "order" bigint not null,
    user_id uuid not null constraint fk_users_guitar_tunings references public.users
);

create table public.albums(
    id uuid not null primary key,
    title varchar(100) not null,
    release_date timestamp with time zone,
    image_url text,
    artist_id uuid constraint fk_artists_albums references public.artists on delete set null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    user_id uuid not null constraint fk_users_albums references public.users
);

create table public.songs(
    id uuid not null primary key,
    title varchar(100) not null,
    description text not null,
    release_date timestamp with time zone,
    image_url text,
    is_recorded boolean,
    bpm bigint,
    difficulty text,
    songsterr_link text,
    youtube_link text,
    album_track_no bigint,
    last_time_played timestamp with time zone,
    rehearsals numeric not null,
    confidence numeric not null,
    progress numeric not null,
    album_id uuid constraint fk_albums_songs references public.albums on delete set null,
    artist_id uuid constraint fk_artists_songs references public.artists on delete set null,
    guitar_tuning_id uuid constraint fk_guitar_tunings_songs references public.guitar_tunings on delete set null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    user_id uuid not null constraint fk_users_songs references public.users
);

create table public.playlist_songs(
    playlist_id uuid not null constraint fk_playlists_playlist_songs references public.playlists on delete cascade,
    song_id uuid not null constraint fk_songs_playlist_songs references public.songs on delete cascade,
    song_track_no bigint not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    primary key(playlist_id, song_id)
);

create table public.song_section_types(
    id uuid not null primary key,
    name varchar(16),
    "order" bigint not null,
    user_id uuid not null constraint fk_users_song_section_types references public.users
);

create table public.song_sections(
    id uuid not null primary key,
    name varchar(30),
    "order" bigint not null,
    rehearsals bigint not null,
    confidence bigint not null,
    rehearsals_score bigint not null,
    confidence_score bigint not null,
    progress bigint not null,
    song_id uuid not null constraint fk_songs_sections references public.songs on delete cascade,
    song_section_type_id uuid not null constraint fk_song_section_types_sections references public.song_section_types on delete cascade,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);

create table public.song_section_histories(
    id uuid not null primary key,
    property varchar(255) not null,
    "from" bigint not null,
    "to" bigint not null,
    song_section_id uuid not null constraint fk_song_sections_history references public.song_sections on delete cascade,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table public.song_section_histories;

drop table public.song_sections;

drop table public.song_section_types;

drop table public.playlist_songs;

drop table public.songs;

drop table public.albums;

drop table public.guitar_tunings;

drop table public.playlists;

drop table public.artists;

drop table public.users;

-- +goose StatementEnd

import { IconSearch } from '@tabler/icons-react'
import {
  alpha,
  Avatar,
  Center,
  Chip,
  Combobox,
  ComboboxProps,
  Group,
  Highlight,
  Indicator,
  LoadingOverlay,
  MantineStyleProps,
  MantineTheme,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  TextInputProps,
  useCombobox
} from '@mantine/core'
import { useDebouncedValue } from '@mantine/hooks'
import { useGetSearchQuery } from '../../state/api/searchApi.ts'
import artistPlaceholder from '../../assets/user-placeholder.jpg'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import playlistPlaceholder from '../../assets/image-placeholder-1.jpg'
import { AlbumSearch, ArtistSearch, PlaylistSearch, SongSearch } from '../../types/models/Search.ts'
import { useNavigate } from 'react-router-dom'
import SearchType from '../../utils/enums/SearchType.ts'
import { MouseEvent, ReactNode, useRef, useState } from 'react'
import CustomIconAlbumVinyl from '../@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconMusicNoteEighth from '../@ui/icons/CustomIconMusicNoteEighth.tsx'
import CustomIconPlaylist2 from '../@ui/icons/CustomIconPlaylist2.tsx'
import useSearchQueryCacheInvalidation from "../../hooks/useSearchQueryCacheInvalidation.ts";

const optionStyle = (theme: MantineTheme) => ({
  borderRadius: '12px',
  transition: '0.15s',
  '&:hover': { backgroundColor: alpha(theme.colors.gray[1], 0.7) }
})

const optionProps: MantineStyleProps = {
  pl: 'sm',
  pr: 0,
  mx: 'xs'
}

const AvatarIndicator = ({ src, alt, icon }: { src: string; alt: string; icon: ReactNode }) => (
  <Indicator
    position={'bottom-end'}
    color={'transparent'}
    label={
      <Center
        bg={'gray.0'}
        c={'primary.5'}
        style={(theme) => ({ borderRadius: '50%', boxShadow: theme.shadows.md })}
        p={3}
      >
        {icon}
      </Center>
    }
  >
    <Avatar
      radius={'md'}
      src={src}
      alt={alt}
      style={(theme) => ({ boxShadow: theme.shadows.sm })}
    />
  </Indicator>
)

const TypeChip = ({
  type,
  children,
  onClick
}: {
  type: SearchType
  children: string
  onClick: (e: MouseEvent) => void
}) => (
  <Chip variant={'light'} size={'xs'} value={type} onClick={onClick}>
    {children}
  </Chip>
)

interface TopbarSearchProps extends TextInputProps {
  comboboxProps?: ComboboxProps
  dropdownMinHeight?: number | string
}

function TopbarSearch({ comboboxProps, dropdownMinHeight = 200, ...others }: TopbarSearchProps) {
  useSearchQueryCacheInvalidation()

  const textInputRef = useRef(null)

  const navigate = useNavigate()
  const combobox = useCombobox()

  const [value, setValue] = useState('')
  const [search] = useDebouncedValue(value, 200)
  const [type, setType] = useState<SearchType | null>(null)

  const { data: searchResults, isFetching } = useGetSearchQuery({
    query: search,
    currentPage: 1,
    pageSize: 20,
    type: type === null ? undefined : type,
    order: search.trim() !== '' ? [] : ['createdAt:desc']
  })

  const ArtistOption = ({ artist }: { artist: ArtistSearch }) => (
    <Combobox.Option
      value={artist.name}
      aria-label={artist.name}
      onClick={() => navigate(`/artist/${artist.id}`)}
      sx={optionStyle}
      {...optionProps}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <Avatar
          src={artist.imageUrl ?? artistPlaceholder}
          alt={artist.name}
          style={(theme) => ({ boxShadow: theme.shadows.sm })}
        />
        <Highlight
          highlight={search}
          highlightStyles={{ fontWeight: 800 }}
          fw={500}
          lineClamp={2}
          lh={'xxs'}
        >
          {artist.name}
        </Highlight>
      </Group>
    </Combobox.Option>
  )

  const AlbumOption = ({ album }: { album: AlbumSearch }) => (
    <Combobox.Option
      value={album.title}
      aria-label={album.title}
      onClick={() => navigate(`/album/${album.id}`)}
      sx={optionStyle}
      {...optionProps}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <AvatarIndicator
          src={album.imageUrl ?? albumPlaceholder}
          alt={album.title}
          icon={<CustomIconAlbumVinyl size={12} />}
        />
        <Stack gap={0}>
          <Highlight
            highlight={search}
            highlightStyles={{ fontWeight: 800 }}
            lh={'xxs'}
            fw={500}
            lineClamp={1}
          >
            {album.title}
          </Highlight>
          {album.artist && (
            <Highlight
              highlight={search}
              highlightStyles={{ fontWeight: 800 }}
              lh={'xxs'}
              c={'dimmed'}
              fz={'xs'}
              fw={500}
              lineClamp={1}
            >
              {album.artist.name}
            </Highlight>
          )}
        </Stack>
      </Group>
    </Combobox.Option>
  )

  const SongOption = ({ song }: { song: SongSearch }) => (
    <Combobox.Option
      value={song.title}
      aria-label={song.title}
      onClick={() => navigate(`/song/${song.id}`)}
      sx={optionStyle}
      {...optionProps}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <AvatarIndicator
          src={song.imageUrl ?? song.album?.imageUrl ?? songPlaceholder}
          alt={song.title}
          icon={<CustomIconMusicNoteEighth size={12} />}
        />
        <Stack gap={0}>
          <Highlight
            highlight={search}
            highlightStyles={{ fontWeight: 800 }}
            lh={'xxs'}
            fw={500}
            lineClamp={1}
          >
            {song.title}
          </Highlight>
          {song.artist && (
            <Highlight
              highlight={search}
              highlightStyles={{ fontWeight: 800 }}
              lh={'xxs'}
              c={'dimmed'}
              fz={'xs'}
              fw={500}
              lineClamp={1}
            >
              {song.artist.name}
            </Highlight>
          )}
        </Stack>
      </Group>
    </Combobox.Option>
  )

  const PlaylistOption = ({ playlist }: { playlist: PlaylistSearch }) => (
    <Combobox.Option
      value={playlist.title}
      aria-label={playlist.title}
      onClick={() => navigate(`/playlist/${playlist.id}`)}
      sx={optionStyle}
      {...optionProps}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <AvatarIndicator
          src={playlist.imageUrl ?? playlistPlaceholder}
          alt={playlist.title}
          icon={<CustomIconPlaylist2 size={12} />}
        />
        <Highlight
          highlight={search}
          highlightStyles={{ fontWeight: 800 }}
          lh={'xxs'}
          fw={500}
          lineClamp={2}
        >
          {playlist.title}
        </Highlight>
      </Group>
    </Combobox.Option>
  )

  function handleSubmit() {
    setValue('')
    combobox.closeDropdown()
    textInputRef.current.blur()
    setType(null)
  }

  function handleChipClick(event: MouseEvent<HTMLInputElement>) {
    if (event.currentTarget.value === type) {
      setType(null)
    }
  }

  return (
    <Combobox onOptionSubmit={handleSubmit} store={combobox} {...comboboxProps}>
      <Combobox.Target>
        <TextInput
          ref={textInputRef}
          role={'searchbox'}
          aria-label={'search'}
          placeholder={'Search'}
          leftSection={<IconSearch size={16} stroke={2} />}
          value={value}
          fw={500}
          radius={'lg'}
          styles={(theme) => ({
            input: {
              transition: '0.3s',
              backgroundColor: alpha(theme.colors.gray[0], 0.1),
              borderWidth: 0,
              '&:focus, &:hover': {
                boxShadow: theme.shadows.sm,
                backgroundColor: alpha(theme.colors.gray[0], 0.2)
              },
              ...(combobox.dropdownOpened && {
                boxShadow: theme.shadows.sm,
                backgroundColor: alpha(theme.colors.gray[0], 0.2)
              })
            }
          })}
          {...others}
          onChange={(event) => {
            setValue(event.currentTarget.value)
            combobox.openDropdown()
            combobox.updateSelectedOptionIndex()
          }}
          onFocus={() => combobox.openDropdown()}
        />
      </Combobox.Target>

      <Combobox.Dropdown p={0}>
        <LoadingOverlay visible={isFetching} />

        <Stack gap={'xs'}>
          <Chip.Group multiple={false} value={type} onChange={(e) => setType(e as SearchType)}>
            <Group px={'xs'} pt={'xs'} gap={6} wrap={'nowrap'} style={{ alignSelf: 'center' }}>
              <TypeChip type={SearchType.Artist} onClick={handleChipClick}>
                Artists
              </TypeChip>
              <TypeChip type={SearchType.Album} onClick={handleChipClick}>
                Albums
              </TypeChip>
              <TypeChip type={SearchType.Song} onClick={handleChipClick}>
                Songs
              </TypeChip>
              <TypeChip type={SearchType.Playlist} onClick={handleChipClick}>
                Playlists
              </TypeChip>
            </Group>
          </Chip.Group>

          <Stack gap={0}>
            {(searchResults?.totalCount !== 0 ||
              (searchResults?.totalCount !== 0 && search.trim() !== '')) && (
              <Text fw={500} px={'lg'} c={'dimmed'} fz={'xs'} pt={'xxs'}>
                {search.trim() === '' ? 'Recently added' : 'Search results'}
              </Text>
            )}

            <Combobox.Options pt={'xxs'}>
              <ScrollArea.Autosize mah={dropdownMinHeight} scrollbarSize={5}>
                {searchResults?.totalCount === 0 && search.trim() === '' ? (
                  <Combobox.Empty pb={'md'} fw={500}>
                    There is nothing in your library
                  </Combobox.Empty>
                ) : searchResults?.totalCount === 0 ? (
                  <Combobox.Empty pb={'md'} fw={500}>
                    No results found
                  </Combobox.Empty>
                ) : (
                  <Stack gap={0} py={'xxs'}>
                    {searchResults?.models?.map((result) =>
                      result.type === SearchType.Artist ? (
                        <ArtistOption key={result.id} artist={result as ArtistSearch} />
                      ) : result.type === SearchType.Album ? (
                        <AlbumOption key={result.id} album={result as AlbumSearch} />
                      ) : result.type === SearchType.Song ? (
                        <SongOption key={result.id} song={result as SongSearch} />
                      ) : result.type === SearchType.Playlist ? (
                        <PlaylistOption key={result.id} playlist={result as PlaylistSearch} />
                      ) : (
                        <></>
                      )
                    )}
                  </Stack>
                )}
              </ScrollArea.Autosize>
            </Combobox.Options>
          </Stack>
        </Stack>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default TopbarSearch

import {
  alpha,
  Avatar,
  Box,
  Button,
  Center,
  Checkbox,
  Chip,
  Group,
  Highlight,
  LoadingOverlay,
  Modal,
  ScrollArea,
  Skeleton,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import {
  useDebouncedValue,
  useDidUpdate,
  useFocusTrap,
  useInputState,
  useListState
} from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useAddSongsToPlaylistMutation } from '../../../state/api/playlistsApi.ts'
import { useGetSongsQuery } from '../../../state/api/songsApi.ts'
import { IconSearch } from '@tabler/icons-react'
import { MouseEvent } from 'react'
import CustomIconMusicNoteEighth from '../../@ui/icons/CustomIconMusicNoteEighth.tsx'
import OrderType from '../../../types/enums/OrderType.ts'
import SongProperty from '../../../types/enums/SongProperty.ts'
import useOrderBy from '../../../hooks/api/useOrderBy.ts'
import useSearchBy from '../../../hooks/api/useSearchBy.ts'
import useFilters from '../../../hooks/filter/useFilters.ts'
import FilterOperator from '../../../types/enums/FilterOperator.ts'
import useFiltersHandlers from '../../../hooks/filter/useFiltersHandlers.ts'
import Song from '../../../types/models/Song.ts'

const SongsLoader = () => (
  <Box data-testid={'songs-loader'}>
    {Array.from(Array(5)).map((_, i) => (
      <Group key={i} w={'100%'} px={'xl'} py={'xs'}>
        <Skeleton mr={'sm'} radius={'md'} width={22} height={22} />
        <Skeleton width={37} height={37} radius={'md'} />
        <Skeleton width={160} height={18} />
      </Group>
    ))}
  </Box>
)

function SongOption({
  song,
  selectedSongs,
  checkSong,
  searchValue
}: {
  song: Song
  selectedSongs: Song[]
  checkSong: (song: Song, checked: boolean) => void
  searchValue: string
}) {
  const checked = selectedSongs.some((s) => s.id === song.id)

  return (
    <Group
      aria-label={`song-${song.title}`}
      aria-selected={checked}
      w={'100%'}
      wrap={'nowrap'}
      px={'xl'}
      py={'xs'}
      sx={(theme) => ({
        cursor: 'pointer',
        transition: '0.3s',
        '&:hover': {
          boxShadow: theme.shadows.xl,
          backgroundColor: alpha(theme.colors.primary[0], 0.15)
        }
      })}
      onClick={() => checkSong(song, !checked)}
    >
      <Checkbox
        aria-label={song.title}
        checked={checked}
        onChange={(e) => checkSong(song, e.currentTarget.checked)}
        onClick={(e) => e.stopPropagation()}
        pr={'sm'}
      />
      <Avatar
        radius={'md'}
        src={song.imageUrl ?? song.album?.imageUrl}
        alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
        bg={'gray.5'}
      >
        <Center c={'white'}>
          <CustomIconMusicNoteEighth aria-label={`default-icon-${song.title}`} size={18} />
        </Center>
      </Avatar>
      <Stack gap={0} style={{ overflow: 'hidden' }}>
        <Group gap={'xxs'} wrap={'nowrap'}>
          <Highlight
            highlight={searchValue}
            highlightStyles={{ fontWeight: 800 }}
            fw={500}
            truncate={'end'}
          >
            {song.title}
          </Highlight>
          {song.album && (
            <Group gap={'xxs'} wrap={'nowrap'}>
              <Text fz={'sm'} c={'dimmed'}>
                -
              </Text>
              <Highlight highlight={searchValue} fz={'sm'} c={'dimmed'} lineClamp={1}>
                {song.album.title}
              </Highlight>
            </Group>
          )}
        </Group>
        {song.artist && (
          <Highlight highlight={searchValue} fz={'sm'} c={'dimmed'} truncate={'end'}>
            {song.artist.name}
          </Highlight>
        )}
      </Stack>
    </Group>
  )
}

interface AddPlaylistSongsModalProps {
  opened: boolean
  onClose: () => void
  playlistId: string
}

function AddPlaylistSongsModal({ opened, onClose, playlistId }: AddPlaylistSongsModalProps) {
  const searchRef = useFocusTrap()
  const [search, setSearch] = useInputState('')
  const [searchValue] = useDebouncedValue(search, 200)

  const [addSongsMutation, { isLoading: addSongsIsLoading }] = useAddSongsToPlaylistMutation()

  const orderBy = useOrderBy([{ property: SongProperty.LastModified, type: OrderType.Descending }])
  const [filters, setFilters] = useFilters([
    {
      property: SongProperty.PlaylistId,
      operator: FilterOperator.NotEqualVariant,
      value: playlistId,
      isSet: true
    },
    { property: SongProperty.Title, operator: FilterOperator.PatternMatching, isSet: false }
  ])
  const { handleValueChange, handleIsSetChange } = useFiltersHandlers(filters, setFilters)
  const searchBy = useSearchBy(filters)

  useDidUpdate(
    () =>
      handleValueChange(SongProperty.Title + FilterOperator.PatternMatching, searchValue.trim()),
    [searchValue]
  )

  const {
    data,
    isLoading: songsIsLoading,
    isFetching: songsIsFetching
  } = useGetSongsQuery({
    currentPage: 1,
    pageSize: 20,
    orderBy: orderBy,
    searchBy: searchBy
  })
  const totalCount = data?.totalCount
  const songs = data?.models ?? []
  const [selectedSongs, selectedSongsHandlers] = useListState<Song>([])
  const filteredSongs = songs.filter((s) => !selectedSongs.some((ss) => s.id === ss.id))
  const totalSongs = selectedSongs.concat(filteredSongs)

  const areAllSongsChecked =
    filteredSongs.length === 0 && totalSongs.length === selectedSongs.length

  function checkAllSongs(check: boolean) {
    if (check) {
      filteredSongs.forEach((song) => selectedSongsHandlers.append(song))
    } else {
      selectedSongsHandlers.setState([])
    }
  }

  function checkSong(song: Song, check: boolean) {
    if (check) {
      selectedSongsHandlers.append(song)
    } else {
      selectedSongsHandlers.filter((s) => s.id !== song.id)
    }
  }

  async function addSongs(e: MouseEvent) {
    if (selectedSongs.length === 0) {
      e.preventDefault()
      return
    }

    await addSongsMutation({
      id: playlistId,
      songIds: selectedSongs.map((s) => s.id),
      forceAdd: true
    }).unwrap()

    toast.success(`Songs added to playlist!`)
    onClose()
    selectedSongsHandlers.setState([])
    setSearch('')
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={'Add Playlist Songs'}
      styles={{ body: { padding: 0 } }}
      trapFocus={false}
    >
      <ScrollArea.Autosize offsetScrollbars={'y'} scrollbars={'y'} scrollbarSize={7} mah={'77vh'}>
        <LoadingOverlay visible={addSongsIsLoading} loaderProps={{ type: 'bars' }} />

        <Stack align={'center'} w={'100%'}>
          <Text fw={500} fz={'lg'}>
            Choose songs
          </Text>

          <TextInput
            ref={searchRef}
            w={250}
            role={'searchbox'}
            aria-label={'search'}
            leftSection={<IconSearch size={15} />}
            placeholder={'Search by title'}
            disabled={songsIsLoading}
            value={search}
            onChange={setSearch}
          />

          {totalCount === 0 && selectedSongs.length === 0 && searchValue.trim() === '' && (
            <Text>There are no songs without playlist</Text>
          )}
          {totalCount === 0 && selectedSongs.length === 0 && searchValue.trim() !== '' && (
            <Text>There are no songs with that title</Text>
          )}

          {(totalCount > 0 || selectedSongs.length > 0) && (
            <Group w={'100%'} pl={'xl'} pr={'md'} justify={'space-between'}>
              <Checkbox
                label={areAllSongsChecked ? 'Deselect all' : 'Select all'}
                checked={areAllSongsChecked}
                onChange={(e) => checkAllSongs(e.currentTarget.checked)}
              />
              <Chip
                checked={
                  !filters.get(SongProperty.PlaylistId + FilterOperator.NotEqualVariant).isSet
                }
                onChange={(checked) =>
                  handleIsSetChange(
                    SongProperty.PlaylistId + FilterOperator.NotEqualVariant,
                    !checked
                  )
                }
                variant={'filled'}
                size={'xs'}
              >
                Show All
              </Chip>
            </Group>
          )}

          <ScrollArea.Autosize
            mah={'50vh'}
            w={'100%'}
            scrollbars={'y'}
            scrollbarSize={7}
            styles={{
              viewport: {
                '> div': {
                  width: 0,
                  minWidth: '100%'
                }
              }
            }}
          >
            <Stack gap={0}>
              <LoadingOverlay
                data-testid={'loading-overlay-fetching'}
                visible={!songsIsLoading && songsIsFetching}
              />
              {songsIsLoading ? (
                <SongsLoader />
              ) : (
                totalSongs.map((song) => (
                  <SongOption
                    key={song.id}
                    song={song}
                    selectedSongs={selectedSongs}
                    checkSong={checkSong}
                    searchValue={searchValue}
                  />
                ))
              )}
            </Stack>
          </ScrollArea.Autosize>

          <Box p={'md'} style={{ alignSelf: 'end' }}>
            <Tooltip disabled={selectedSongs.length > 0} label="Select songs">
              <Button data-disabled={selectedSongs.length === 0} onClick={addSongs}>
                Add
              </Button>
            </Tooltip>
          </Box>
        </Stack>
      </ScrollArea.Autosize>
    </Modal>
  )
}

export default AddPlaylistSongsModal

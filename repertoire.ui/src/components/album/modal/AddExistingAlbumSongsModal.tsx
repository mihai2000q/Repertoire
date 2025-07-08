import {
  alpha,
  Avatar,
  Box,
  Button,
  Center,
  Checkbox,
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
import { useDebouncedValue, useFocusTrap, useInputState, useListState } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { IconInfoCircleFilled, IconSearch } from '@tabler/icons-react'
import { MouseEvent } from 'react'
import { useAddSongsToAlbumMutation } from '../../../state/api/albumsApi.ts'
import { SongSearch } from '../../../types/models/Search.ts'
import { useGetSearchQuery } from '../../../state/api/searchApi.ts'
import SearchType from '../../../types/enums/SearchType.ts'
import CustomIconMusicNoteEighth from '../../@ui/icons/CustomIconMusicNoteEighth.tsx'

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
  song: SongSearch
  selectedSongs: SongSearch[]
  checkSong: (song: SongSearch, checked: boolean) => void
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
      <Stack gap={0}>
        <Highlight
          highlight={searchValue}
          highlightStyles={{ fontWeight: 800 }}
          fw={500}
          lineClamp={2}
        >
          {song.title}
        </Highlight>
        {song.artist && (
          <Highlight highlight={searchValue} fz={'sm'} c={'dimmed'} lineClamp={1}>
            {song.artist.name}
          </Highlight>
        )}
      </Stack>
    </Group>
  )
}

interface AddExistingAlbumSongsModalProps {
  opened: boolean
  onClose: () => void
  albumId: string
  artistId: string
}

function AddExistingAlbumSongsModal({
  opened,
  onClose,
  albumId,
  artistId
}: AddExistingAlbumSongsModalProps) {
  const searchRef = useFocusTrap()
  const [search, setSearch] = useInputState('')
  const [searchValue] = useDebouncedValue(search, 200)

  const [addSongsMutation, { isLoading: addSongsIsLoading }] = useAddSongsToAlbumMutation()

  const {
    data,
    isLoading: songsIsLoading,
    isFetching: songsIsFetching
  } = useGetSearchQuery({
    query: searchValue,
    currentPage: 1,
    pageSize: 20,
    type: SearchType.Song,
    filter: ['album IS NULL', `(artist IS NULL${artistId ? ` OR artist.id = ${artistId})` : ')'}`],
    order: ['updatedAt:desc']
  })
  const totalCount = data?.totalCount
  const songs = (data?.models as SongSearch[]) ?? []
  const [selectedSongs, selectedSongsHandlers] = useListState<SongSearch>([])
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

  function checkSong(song: SongSearch, check: boolean) {
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

    await addSongsMutation({ id: albumId, songIds: selectedSongs.map((s) => s.id) }).unwrap()

    toast.success(`Songs added to album!`)
    onClose()
    selectedSongsHandlers.setState([])
    setSearch('')
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={'Add Existing Songs'}
      styles={{ body: { padding: 0 } }}
      trapFocus={false}
    >
      <ScrollArea.Autosize offsetScrollbars={'y'} scrollbars={'y'} scrollbarSize={7} mah={'77vh'}>
        <LoadingOverlay visible={addSongsIsLoading} loaderProps={{ type: 'bars' }} />

        <Stack align={'center'} w={'100%'}>
          <Group gap={'xxs'} align={'start'}>
            <Text fw={500} ta={'center'} fz={'lg'}>
              Choose songs
            </Text>
            <Tooltip label={'All songs will inherit the artist of the album'}>
              <Center c={'primary.8'}>
                <IconInfoCircleFilled aria-label={'info-icon'} size={15} />
              </Center>
            </Tooltip>
          </Group>

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
            <Text>There are no songs without album</Text>
          )}
          {totalCount === 0 && selectedSongs.length === 0 && searchValue.trim() !== '' && (
            <Text>There are no songs with that title</Text>
          )}

          {(totalCount > 0 || selectedSongs.length > 0) && (
            <Checkbox
              label={areAllSongsChecked ? 'Deselect all' : 'Select all'}
              checked={areAllSongsChecked}
              onChange={(e) => checkAllSongs(e.currentTarget.checked)}
              px={'xl'}
              style={{ alignSelf: 'flex-start' }}
            />
          )}

          <ScrollArea.Autosize mah={'50vh'} w={'100%'} scrollbars={'y'} scrollbarSize={7}>
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

export default AddExistingAlbumSongsModal

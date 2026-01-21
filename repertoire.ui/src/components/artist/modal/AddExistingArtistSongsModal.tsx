import {
  alpha,
  Avatar,
  Box,
  Button,
  Center,
  Checkbox,
  Group,
  Highlight,
  Loader,
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
  useIntersection,
  useListState
} from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useAddSongsToArtistMutation } from '../../../state/api/artistsApi.ts'
import { IconSearch } from '@tabler/icons-react'
import { MouseEvent, useEffect, useRef } from 'react'
import { useGetInfiniteSearchInfiniteQuery } from '../../../state/api/searchApi.ts'
import SearchType from '../../../types/enums/SearchType.ts'
import { SongSearch } from '../../../types/models/Search.ts'
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
        {song.album && (
          <Highlight highlight={searchValue} fz={'sm'} c={'dimmed'} lineClamp={1}>
            {song.album.title}
          </Highlight>
        )}
      </Stack>
    </Group>
  )
}

interface AddExistingArtistSongsModalProps {
  opened: boolean
  onClose: () => void
  artistId: string
}

function AddExistingArtistSongsModal({
  opened,
  onClose,
  artistId
}: AddExistingArtistSongsModalProps) {
  const [search, setSearch] = useInputState('')
  const [searchValue] = useDebouncedValue(search, 200)

  const [addSongsMutation, { isLoading: addSongsIsLoading }] = useAddSongsToArtistMutation()

  const {
    data,
    isLoading: songsIsLoading,
    isFetching: songsIsFetching,
    isFetchingNextPage,
    fetchNextPage
  } = useGetInfiniteSearchInfiniteQuery({
    query: searchValue,
    type: SearchType.Song,
    filter: ['artist IS NULL'],
    order: ['updatedAt:desc']
  })
  const totalCount = data?.pages[0].totalCount
  const songs = data?.pages.flatMap((s) => (s.models as SongSearch[]) ?? []) ?? []
  const [selectedSongs, selectedSongsHandlers] = useListState<SongSearch>([])
  const filteredSongs = songs.filter((s) => !selectedSongs.some((ss) => s.id === ss.id))
  const totalSongs = selectedSongs.concat(filteredSongs)

  const areAllSongsChecked =
    filteredSongs.length === 0 && totalSongs.length === selectedSongs.length

  const searchRef = useFocusTrap(!songsIsLoading)

  const scrollRef = useRef<HTMLDivElement>()
  const { ref: lastSongRef, entry } = useIntersection({
    root: scrollRef.current,
    threshold: 0.1
  })
  useEffect(() => {
    if (entry?.isIntersecting === true) fetchNextPage()
  }, [entry?.isIntersecting])
  useDidUpdate(() => scrollRef.current.scrollTo({ top: 0, behavior: 'instant' }), [searchValue])

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

    await addSongsMutation({ id: artistId, songIds: selectedSongs.map((s) => s.id) }).unwrap()

    toast.success(`Songs added to artist!`)
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
            <Text>There are no songs without artist</Text>
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

          <ScrollArea.Autosize
            mah={'50vh'}
            w={'100%'}
            scrollbars={'y'}
            scrollbarSize={7}
            viewportRef={scrollRef}
          >
            <Stack gap={0}>
              <LoadingOverlay
                data-testid={'loading-overlay-fetching'}
                visible={!songsIsLoading && songsIsFetching && !isFetchingNextPage}
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

              <Stack gap={0} align={'center'}>
                <div ref={lastSongRef} />
                {isFetchingNextPage && <Loader size={30} mt={'sm'} />}
              </Stack>
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

export default AddExistingArtistSongsModal

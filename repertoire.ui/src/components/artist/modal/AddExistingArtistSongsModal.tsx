import {
  alpha,
  Avatar,
  Box,
  Button,
  Checkbox,
  Group,
  LoadingOverlay,
  Modal,
  Skeleton,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import { useDebouncedValue, useInputState, useListState } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useAddSongsToArtistMutation } from '../../../state/artistsApi.ts'
import { useGetSongsQuery } from '../../../state/songsApi.ts'
import { IconSearch } from '@tabler/icons-react'
import songPlaceholder from '../../../assets/image-placeholder-1.jpg'
import { MouseEvent, useEffect } from 'react'

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

  const {
    data: songs,
    isLoading: songsIsLoading,
    isFetching: songsIsFetching
  } = useGetSongsQuery({
    currentPage: 1,
    pageSize: 20,
    orderBy: ['title asc'],
    searchBy: [
      'songs.artist_id IS NULL',
      ...(searchValue.trim() === '' ? [] : [`title ~* '${searchValue}'`])
    ]
  })

  const [addSongMutation, { isLoading: addSongIsLoading }] = useAddSongsToArtistMutation()

  const [songIds, songIdsHandlers] = useListState<string>([])

  useEffect(() => {
    songIdsHandlers.filter((songId) => songs.models.some((song) => song.id === songId))
  }, [searchValue, songs])

  function checkAllSongs(check: boolean) {
    songIdsHandlers.setState([])
    if (check) {
      songs.models.forEach((song) => songIdsHandlers.append(song.id))
    }
  }

  function checkSong(songId: string, check: boolean) {
    if (check) {
      songIdsHandlers.append(songId)
    } else {
      songIdsHandlers.filter((s) => s !== songId)
    }
  }

  async function addSongs(e: MouseEvent) {
    if (songIds.length === 0) {
      e.preventDefault()
      return
    }

    await addSongMutation({ id: artistId, songIds }).unwrap()

    toast.success(`Songs added to artist!`)
    onClose()
    songIdsHandlers.setState([])
    setSearch('')
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={'Add Existing Songs'}
      styles={{ body: { padding: 0 } }}
    >
      <Modal.Body p={0} pos={'relative'}>
        <LoadingOverlay visible={addSongIsLoading} loaderProps={{ type: 'bars' }} />

        <Stack align={'center'} w={'100%'}>
          <Text fw={500} fz={'lg'}>
            Choose songs
          </Text>

          <TextInput
            w={250}
            role={'searchbox'}
            aria-label={'search'}
            leftSection={<IconSearch size={15} />}
            placeholder={'Search by title'}
            disabled={songsIsLoading}
            value={search}
            onChange={setSearch}
          />

          {songs?.totalCount === 0 && <Text>There are no songs without artist</Text>}
          {songs?.totalCount > 0 && (
            <Group w={'100%'} px={'xl'}>
              <Checkbox
                aria-label={songIds.length === songs.models.length ? 'deselect-all' : 'select-all'}
                checked={songIds.length === songs.models.length}
                onChange={(e) => checkAllSongs(e.currentTarget.checked)}
              />
              <Text>{songIds.length === songs.models.length ? 'Deselect' : 'Select'} All</Text>
            </Group>
          )}

          <Stack w={'100%'} gap={0} style={{ overflow: 'auto', maxHeight: '50vh' }}>
            <LoadingOverlay
              data-testid={'loading-overlay-fetching'}
              visible={!songsIsLoading && songsIsFetching}
            />
            {songsIsLoading ? (
              <Box data-testid={'songs-loader'}>
                {Array.from(Array(5)).map((_, i) => (
                  <Group key={i} w={'100%'} px={'xl'} py={'xs'}>
                    <Skeleton mr={'sm'} radius={'md'} width={22} height={22} />
                    <Skeleton width={37} height={37} radius={'md'} />
                    <Skeleton width={160} height={18} />
                  </Group>
                ))}
              </Box>
            ) : (
              songs.models.map((song) => (
                <Group
                  key={song.id}
                  aria-label={`song-${song.title}`}
                  aria-selected={songIds.some((id) => id === song.id)}
                  sx={(theme) => ({
                    transition: '0.3s',
                    '&:hover': {
                      boxShadow: theme.shadows.xl,
                      backgroundColor: alpha(theme.colors.cyan[0], 0.15)
                    }
                  })}
                  w={'100%'}
                  px={'xl'}
                  py={'xs'}
                >
                  <Checkbox
                    aria-label={song.title}
                    checked={songIds.includes(song.id)}
                    onChange={(e) => checkSong(song.id, e.currentTarget.checked)}
                    pr={'sm'}
                  />
                  <Avatar
                    radius={'md'}
                    src={song.imageUrl ?? song.album?.imageUrl ?? songPlaceholder}
                    alt={song.title}
                  />
                  <Stack gap={0}>
                    <Text fw={500}>{song.title}</Text>
                    {song.album && (
                      <Text fz={'sm'} c={'dimmed'}>
                        {song.album.title}
                      </Text>
                    )}
                  </Stack>
                </Group>
              ))
            )}
          </Stack>

          <Box p={'md'} style={{ alignSelf: 'end' }}>
            <Tooltip disabled={songIds.length > 0} label="Select songs">
              <Button data-disabled={songIds.length === 0} onClick={addSongs}>
                Add
              </Button>
            </Tooltip>
          </Box>
        </Stack>
      </Modal.Body>
    </Modal>
  )
}

export default AddExistingArtistSongsModal

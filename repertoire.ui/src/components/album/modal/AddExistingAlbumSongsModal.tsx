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
import { useGetSongsQuery } from '../../../state/api/songsApi.ts'
import { IconInfoCircleFilled, IconSearch } from '@tabler/icons-react'
import songPlaceholder from '../../../assets/image-placeholder-1.jpg'
import { MouseEvent, useEffect } from 'react'
import { useAddSongsToAlbumMutation } from '../../../state/api/albumsApi.ts'

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
      'album_id IS NULL',
      `songs.artist_id IS NULL${artistId ? ` OR songs.artist_id = '${artistId}'` : ''}`,
      ...(searchValue.trim() === '' ? [] : [`songs.title ~* '${searchValue}'`])
    ]
  })

  const [addSongsMutation, { isLoading: addSongIsLoading }] = useAddSongsToAlbumMutation()

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

    await addSongsMutation({ id: albumId, songIds }).unwrap()

    toast.success(`Songs added to album!`)
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
          <Group gap={4} align={'start'}>
            <Text fw={500} ta={'center'} fz={'lg'}>
              Choose songs
            </Text>
            <Tooltip label={'All songs will inherit the artist of the album'}>
              <Box c={'primary.8'}>
                <IconInfoCircleFilled aria-label={'info-icon'} size={15} />
              </Box>
            </Tooltip>
          </Group>

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

          {songs?.totalCount === 0 && <Text>There are no songs without album</Text>}
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
                      backgroundColor: alpha(theme.colors.primary[0], 0.15)
                    }
                  })}
                  w={'100%'}
                  wrap={'nowrap'}
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
                    <Text fw={500} lineClamp={2}>{song.title}</Text>
                    {song.artist && (
                      <Text fz={'sm'} c={'dimmed'} lineClamp={1}>
                        {song.artist.name}
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

export default AddExistingAlbumSongsModal

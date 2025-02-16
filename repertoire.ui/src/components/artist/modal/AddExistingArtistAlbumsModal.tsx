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
import { useAddAlbumsToArtistMutation } from '../../../state/api/artistsApi.ts'
import { useGetAlbumsQuery } from '../../../state/api/albumsApi.ts'
import { IconInfoCircleFilled, IconSearch } from '@tabler/icons-react'
import albumPlaceholder from '../../../assets/image-placeholder-1.jpg'
import { MouseEvent, useEffect } from 'react'

interface AddExistingArtistAlbumsModalProps {
  opened: boolean
  onClose: () => void
  artistId: string
}

function AddExistingArtistAlbumsModal({
  opened,
  onClose,
  artistId
}: AddExistingArtistAlbumsModalProps) {
  const [search, setSearch] = useInputState('')
  const [searchValue] = useDebouncedValue(search, 200)

  const {
    data: albums,
    isLoading: albumsIsLoading,
    isFetching: albumsIsFetching
  } = useGetAlbumsQuery({
    currentPage: 1,
    pageSize: 20,
    orderBy: ['title asc'],
    searchBy: [
      'artist_id IS NULL',
      ...(searchValue.trim() === '' ? [] : [`title ~* '${searchValue}'`])
    ]
  })

  const [addAlbumMutation, { isLoading: addAlbumIsLoading }] = useAddAlbumsToArtistMutation()

  const [albumIds, albumIdsHandlers] = useListState<string>([])

  useEffect(() => {
    albumIdsHandlers.filter((albumId) => albums.models.some((album) => album.id === albumId))
  }, [searchValue, albums])

  function checkAllAlbums(check: boolean) {
    albumIdsHandlers.setState([])
    if (check) {
      albums.models.forEach((album) => albumIdsHandlers.append(album.id))
    }
  }

  function checkAlbum(albumId: string, check: boolean) {
    if (check) {
      albumIdsHandlers.append(albumId)
    } else {
      albumIdsHandlers.filter((s) => s !== albumId)
    }
  }

  async function addAlbums(e: MouseEvent) {
    if (albumIds.length === 0) {
      e.preventDefault()
      return
    }

    await addAlbumMutation({ id: artistId, albumIds }).unwrap()

    toast.success(`Albums added to artist!`)
    onClose()
    albumIdsHandlers.setState([])
    setSearch('')
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={'Add Existing Albums'}
      styles={{ body: { padding: 0 } }}
    >
      <Modal.Body p={0} pos={'relative'}>
        <LoadingOverlay visible={addAlbumIsLoading} loaderProps={{ type: 'bars' }} />

        <Stack align={'center'} w={'100%'}>
          <Group gap={'xxs'} align={'start'}>
            <Text fw={500} ta={'center'} fz={'lg'}>
              Choose albums
            </Text>
            <Tooltip
              w={230}
              multiline
              ta={'center'}
              label={'All songs related to the added album will be added to the artist too'}
            >
              <Box c={'primary.8'}>
                <IconInfoCircleFilled size={15} aria-label={'info-icon'} />
              </Box>
            </Tooltip>
          </Group>

          <TextInput
            w={250}
            role={'searchbox'}
            aria-label={'search'}
            leftSection={<IconSearch size={15} />}
            placeholder={'Search by title'}
            disabled={albumsIsLoading}
            value={search}
            onChange={setSearch}
          />

          {albums?.totalCount === 0 && <Text>There are no albums without artist</Text>}
          {albums?.totalCount > 0 && (
            <Group w={'100%'} px={'xl'}>
              <Checkbox
                aria-label={
                  albumIds.length === albums.models.length ? 'deselect-all' : 'select-all'
                }
                checked={albumIds.length === albums.models.length}
                onChange={(e) => checkAllAlbums(e.currentTarget.checked)}
              />
              <Text>{albumIds.length === albums.models.length ? 'Deselect' : 'Select'} All</Text>
            </Group>
          )}

          <Stack w={'100%'} gap={0} style={{ overflow: 'auto', maxHeight: '50vh' }}>
            <LoadingOverlay
              data-testid={'loading-overlay-fetching'}
              visible={!albumsIsLoading && albumsIsFetching}
            />
            {albumsIsLoading ? (
              <Box data-testid={'albums-loader'}>
                {Array.from(Array(5)).map((_, i) => (
                  <Group key={i} w={'100%'} px={'xl'} py={'xs'}>
                    <Skeleton mr={'sm'} radius={'md'} width={22} height={22} />
                    <Skeleton width={37} height={37} radius={'md'} />
                    <Skeleton width={160} height={18} />
                  </Group>
                ))}
              </Box>
            ) : (
              albums.models.map((album) => (
                <Group
                  key={album.id}
                  aria-label={`album-${album.title}`}
                  aria-selected={albumIds.some((id) => id === album.id)}
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
                    aria-label={album.title}
                    checked={albumIds.includes(album.id)}
                    onChange={(e) => checkAlbum(album.id, e.currentTarget.checked)}
                    pr={'sm'}
                  />
                  <Avatar
                    radius={'md'}
                    src={album.imageUrl ?? albumPlaceholder}
                    alt={album.title}
                  />
                  <Text fw={500} lineClamp={2}>{album.title}</Text>
                </Group>
              ))
            )}
          </Stack>

          <Box p={'md'} style={{ alignSelf: 'end' }}>
            <Tooltip disabled={albumIds.length > 0} label="Select albums">
              <Button data-disabled={albumIds.length === 0} onClick={addAlbums}>
                Add
              </Button>
            </Tooltip>
          </Box>
        </Stack>
      </Modal.Body>
    </Modal>
  )
}

export default AddExistingArtistAlbumsModal

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
import { useAddAlbumsToArtistMutation } from '../../../state/api/artistsApi.ts'
import { IconInfoCircleFilled, IconSearch } from '@tabler/icons-react'
import { MouseEvent } from 'react'
import { useGetSearchQuery } from '../../../state/api/searchApi.ts'
import SearchType from '../../../types/enums/SearchType.ts'
import { AlbumSearch } from '../../../types/models/Search.ts'
import CustomIconAlbumVinyl from '../../@ui/icons/CustomIconAlbumVinyl.tsx'

const AlbumsLoader = () => (
  <Box data-testid={'albums-loader'}>
    {Array.from(Array(5)).map((_, i) => (
      <Group key={i} w={'100%'} px={'xl'} py={'xs'}>
        <Skeleton mr={'sm'} radius={'md'} width={22} height={22} />
        <Skeleton width={37} height={37} radius={'md'} />
        <Skeleton width={160} height={18} />
      </Group>
    ))}
  </Box>
)

function AlbumOption({
  album,
  selectedAlbums,
  checkAlbum,
  searchValue
}: {
  album: AlbumSearch
  selectedAlbums: AlbumSearch[]
  checkAlbum: (album: AlbumSearch, checked: boolean) => void
  searchValue: string
}) {
  const checked = selectedAlbums.some((a) => a.id === album.id)

  return (
    <Group
      aria-label={`album-${album.title}`}
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
      onClick={() => checkAlbum(album, !checked)}
    >
      <Checkbox
        aria-label={album.title}
        checked={checked}
        onChange={(e) => checkAlbum(album, e.currentTarget.checked)}
        onClick={(e) => e.stopPropagation()}
        pr={'sm'}
      />
      <Avatar radius={'md'} src={album.imageUrl} alt={album.imageUrl && album.title} bg={'gray.5'}>
        <Center c={'white'}>
          <CustomIconAlbumVinyl aria-label={`default-icon-${album.title}`} size={14} />
        </Center>
      </Avatar>
      <Highlight
        highlight={searchValue}
        highlightStyles={{ fontWeight: 800 }}
        fw={500}
        lineClamp={2}
      >
        {album.title}
      </Highlight>
    </Group>
  )
}

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
  const searchRef = useFocusTrap()
  const [search, setSearch] = useInputState('')
  const [searchValue] = useDebouncedValue(search, 200)

  const [addAlbumsMutation, { isLoading: addAlbumsIsLoading }] = useAddAlbumsToArtistMutation()

  const {
    data,
    isLoading: albumsIsLoading,
    isFetching: albumsIsFetching
  } = useGetSearchQuery({
    query: searchValue,
    currentPage: 1,
    pageSize: 20,
    type: SearchType.Album,
    filter: ['artist IS NULL'],
    order: ['updatedAt:desc']
  })
  const totalCount = data?.totalCount
  const albums = data?.models as AlbumSearch[] ?? []
  const [selectedAlbums, selectedAlbumsHandlers] = useListState<AlbumSearch>([])
  const filteredAlbums = albums.filter((a) => !selectedAlbums.some((aa) => a.id === aa.id))
  const totalAlbums = selectedAlbums.concat(filteredAlbums)

  const areAllAlbumsChecked =
    filteredAlbums.length === 0 && totalAlbums.length === selectedAlbums.length

  function checkAllAlbums(check: boolean) {
    if (check) {
      filteredAlbums.forEach((album) => selectedAlbumsHandlers.append(album))
    } else {
      selectedAlbumsHandlers.setState([])
    }
  }

  function checkAlbum(album: AlbumSearch, check: boolean) {
    if (check) {
      selectedAlbumsHandlers.append(album)
    } else {
      selectedAlbumsHandlers.filter((a) => a.id !== album.id)
    }
  }

  async function addAlbums(e: MouseEvent) {
    if (selectedAlbums.length === 0) {
      e.preventDefault()
      return
    }

    await addAlbumsMutation({ id: artistId, albumIds: selectedAlbums.map((a) => a.id) }).unwrap()

    toast.success(`Albums added to artist!`)
    onClose()
    selectedAlbumsHandlers.setState([])
    setSearch('')
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={'Add Existing Albums'}
      styles={{ body: { padding: 0 } }}
      trapFocus={false}
    >
      <ScrollArea.Autosize offsetScrollbars={'y'} scrollbars={'y'} scrollbarSize={7} mah={'77vh'}>
        <LoadingOverlay visible={addAlbumsIsLoading} loaderProps={{ type: 'bars' }} />

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
              <Center c={'primary.8'}>
                <IconInfoCircleFilled size={15} aria-label={'info-icon'} />
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
            disabled={albumsIsLoading}
            value={search}
            onChange={setSearch}
          />

          {totalCount === 0 && selectedAlbums.length === 0 && searchValue.trim() === '' && (
            <Text>There are no albums without artist</Text>
          )}
          {totalCount === 0 && selectedAlbums.length === 0 && searchValue.trim() !== '' && (
            <Text>There are no albums with that title</Text>
          )}

          {(totalCount > 0 || selectedAlbums.length > 0) && (
            <Checkbox
              label={areAllAlbumsChecked ? 'Deselect all' : 'Select all'}
              checked={areAllAlbumsChecked}
              onChange={(e) => checkAllAlbums(e.currentTarget.checked)}
              px={'xl'}
              style={{ alignSelf: 'flex-start' }}
            />
          )}

          <ScrollArea.Autosize mah={'50vh'} w={'100%'} scrollbars={'y'} scrollbarSize={7}>
            <Stack gap={0}>
              <LoadingOverlay
                data-testid={'loading-overlay-fetching'}
                visible={!albumsIsLoading && albumsIsFetching}
              />
              {albumsIsLoading ? (
                <AlbumsLoader />
              ) : (
                totalAlbums.map((album) => (
                  <AlbumOption
                    key={album.id}
                    album={album}
                    selectedAlbums={selectedAlbums}
                    checkAlbum={checkAlbum}
                    searchValue={searchValue}
                  />
                ))
              )}
            </Stack>
          </ScrollArea.Autosize>

          <Box p={'md'} style={{ alignSelf: 'end' }}>
            <Tooltip disabled={selectedAlbums.length > 0} label="Select albums">
              <Button data-disabled={selectedAlbums.length === 0} onClick={addAlbums}>
                Add
              </Button>
            </Tooltip>
          </Box>
        </Stack>
      </ScrollArea.Autosize>
    </Modal>
  )
}

export default AddExistingArtistAlbumsModal

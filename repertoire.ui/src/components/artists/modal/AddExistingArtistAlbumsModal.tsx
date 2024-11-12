import {
  alpha,
  Avatar,
  Box,
  Button,
  Checkbox,
  Group,
  LoadingOverlay,
  Modal, Skeleton,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import { toast } from 'react-toastify'
import { useAddAlbumsToArtistMutation } from '../../../state/artistsApi.ts'
import { useListState } from '@mantine/hooks'
import { useGetAlbumsQuery } from '../../../state/albumsApi.ts'
import { IconSearch } from '@tabler/icons-react'
import albumPlaceholder from '../../../assets/image-placeholder-1.jpg'
import {useState} from "react";

interface AddNewArtistAlbumModalProps {
  opened: boolean
  onClose: () => void
  artistId: string
}

function AddExistingArtistAlbumsModal({ opened, onClose, artistId }: AddNewArtistAlbumModalProps) {
  const [searchValue, setSearchValue] = useState('')

  const { data: albums, isLoading: albumsIsLoading } = useGetAlbumsQuery({})

  const [addAlbumMutation, { isLoading: addAlbumIsLoading }] = useAddAlbumsToArtistMutation()

  const [albumIds, albumIdsHandlers] = useListState<string>([])

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

  async function addAlbums() {
    await addAlbumMutation({ id: artistId, albumIds }).unwrap()

    toast.success(`Albums added to artist!`)
    onClose()
    albumIdsHandlers.setState([])
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={'Add Existing Album'}
      styles={{ body: { padding: 0 } }}
    >
      <Modal.Body p={0} pos={'relative'}>
        <LoadingOverlay
          visible={addAlbumIsLoading}
          zIndex={1000}
          overlayProps={{ radius: 'sm', blur: 2 }}
        />
        <Stack align={'center'} w={'100%'}>
          <Text fw={500} fz={'lg'}>
            Choose albums
          </Text>
          <TextInput
            leftSection={<IconSearch size={15} />}
            placeholder={'Search by title'}
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
          />
          {albums?.totalCount === 0 && <Text>There are no albums without artist</Text>}
          {albums?.totalCount > 0 && (
            <Group w={'100%'} px={'xl'}>
              <Checkbox
                checked={albumIds.length === albums.models.length}
                onChange={(e) => checkAllAlbums(e.currentTarget.checked)}
                disabled={albumsIsLoading}
              />
              <Text>{albumIds.length === albums.models.length ? 'Deselect' : 'Select'} All</Text>
            </Group>
          )}
          <Stack w={'100%'} gap={0} style={{ overflow: 'auto', maxHeight: '50vh' }}>
            {albumsIsLoading ? (
              Array.from(Array(5)).map((_, i) => (
                <Group key={i} w={'100%'} px={'xl'} py={'xs'}>
                  <Skeleton mr={'sm'} radius={'md'} width={22} height={22} />
                  <Skeleton width={37} height={37} radius={'md'} />
                  <Skeleton width={160} height={18} />
                </Group>
              ))
            ) : (
              albums.models.map((album) => (
                <Group
                  key={album.id}
                  align={'center'}
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
                    checked={albumIds.includes(album.id)}
                    onChange={(e) => checkAlbum(album.id, e.currentTarget.checked)}
                    pr={'sm'}
                  />
                  <Avatar radius={'md'} src={album.imageUrl ? album.imageUrl : albumPlaceholder} />
                  <Text fw={500}>{album.title}</Text>
                </Group>
              ))
            )}
          </Stack>
          <Box p={'md'} style={{ alignSelf: 'end' }}>
            <Tooltip label="Select albums">
              <Button
                data-disabled={albumIds.length === 0}
                onClick={(e) => {
                  e.preventDefault()
                  addAlbums().then()
                }}
              >
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

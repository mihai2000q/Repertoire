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
import { useAddSongsToArtistMutation } from '../../../state/artistsApi.ts'
import { useListState } from '@mantine/hooks'
import { useGetSongsQuery } from '../../../state/songsApi.ts'
import { IconSearch } from '@tabler/icons-react'
import songPlaceholder from '../../../assets/image-placeholder-1.jpg'
import {useState} from "react";

interface AddNewArtistSongModalProps {
  opened: boolean
  onClose: () => void
  artistId: string
}

function AddExistingArtistSongsModal({ opened, onClose, artistId }: AddNewArtistSongModalProps) {
  const [searchValue, setSearchValue] = useState('')

  const { data: songs, isLoading: songsIsLoading } = useGetSongsQuery({})

  const [addSongMutation, { isLoading: addSongIsLoading }] = useAddSongsToArtistMutation()

  const [songIds, songIdsHandlers] = useListState<string>([])

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

  async function addSongs() {
    await addSongMutation({ id: artistId, songIds }).unwrap()

    toast.success(`Songs added to artist!`)
    onClose()
    songIdsHandlers.setState([])
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={'Add Existing Song'}
      styles={{ body: { padding: 0 } }}
    >
      <Modal.Body p={0} pos={'relative'}>
        <LoadingOverlay
          visible={addSongIsLoading}
          zIndex={1000}
          overlayProps={{ radius: 'sm', blur: 2 }}
        />
        <Stack align={'center'} w={'100%'}>
          <Text fw={500} fz={'lg'}>
            Choose songs
          </Text>
          <TextInput
            leftSection={<IconSearch size={15} />}
            placeholder={'Search by title'}
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
          />
          {songs?.totalCount === 0 && <Text>There are no songs without artist</Text>}
          {songs?.totalCount > 0 && (
            <Group w={'100%'} px={'xl'}>
              <Checkbox
                checked={songIds.length === songs.models.length}
                onChange={(e) => checkAllSongs(e.currentTarget.checked)}
                disabled={songsIsLoading}
              />
              <Text>{songIds.length === songs.models.length ? 'Deselect' : 'Select'} All</Text>
            </Group>
          )}
          <Stack w={'100%'} gap={0} style={{ overflow: 'auto', maxHeight: '50vh' }}>
            {songsIsLoading ? (
              Array.from(Array(5)).map((_, i) => (
                <Group key={i} w={'100%'} px={'xl'} py={'xs'}>
                  <Skeleton mr={'sm'} radius={'md'} width={22} height={22} />
                  <Skeleton width={37} height={37} radius={'md'} />
                  <Skeleton width={160} height={18} />
                </Group>
              ))
            ) : (
              songs.models.map((song) => (
                <Group
                  key={song.id}
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
                    checked={songIds.includes(song.id)}
                    onChange={(e) => checkSong(song.id, e.currentTarget.checked)}
                    pr={'sm'}
                  />
                  <Avatar radius={'md'} src={song.imageUrl ? song.imageUrl : songPlaceholder} />
                  <Text fw={500}>{song.title}</Text>
                </Group>
              ))
            )}
          </Stack>
          <Box p={'md'} style={{ alignSelf: 'end' }}>
            <Tooltip label="Select songs">
              <Button
                data-disabled={songIds.length === 0}
                onClick={(e) => {
                  e.preventDefault()
                  addSongs().then()
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

export default AddExistingArtistSongsModal

import Song from '../../../types/models/Song.ts'
import { Button, LoadingOverlay, Modal, Stack, Textarea, Tooltip } from '@mantine/core'
import { useUpdateSongMutation } from '../../../state/api/songsApi.ts'
import { useInputState } from '@mantine/hooks'
import { MouseEvent } from 'react'
import { toast } from 'react-toastify'

interface EditSongDescriptionModalProps {
  song: Song
  opened: boolean
  onClose: () => void
}

function EditSongDescriptionModal({ song, opened, onClose }: EditSongDescriptionModalProps) {
  const [updateSongMutation, { isLoading }] = useUpdateSongMutation()

  const [description, setDescription] = useInputState(song.description)

  const hasChanged = description !== song.description

  async function updateSong(e: MouseEvent) {
    if (!hasChanged) {
      e.preventDefault()
      return
    }

    await updateSongMutation({
      ...song,
      guitarTuningId: song.guitarTuning?.id,
      albumId: song.album?.id,
      artistId: song.artist?.id,
      id: song.id,
      description: description
    }).unwrap()

    onClose()

    toast.info('Song description updated!')
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={'Edit Song Description'}
      styles={{ body: { padding: 0 } }}
    >
      <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

      <Stack px={26} pb={'md'}>
        <Textarea
          label={'Description'}
          placeholder={'Enter a description'}
          value={description}
          onChange={setDescription}
          autosize
          minRows={4}
          styles={{
            input: {
              overflow: 'auto',
              maxHeight: '55vh'
            }
          }}
        />

        <Tooltip
          disabled={hasChanged}
          label={'You need to make a change before saving'}
          position="bottom"
        >
          <Button data-disabled={!hasChanged} onClick={updateSong}>
            Save
          </Button>
        </Tooltip>
      </Stack>
    </Modal>
  )
}

export default EditSongDescriptionModal

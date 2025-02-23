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
      artistId: song.album?.id,
      id: song.id,
      description: description
    }).unwrap()

    onClose()

    toast.info('Song description updated!')
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Description'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

        <Stack>
          <Textarea
            label={'Description'}
            placeholder={'Enter a description'}
            value={description}
            onChange={setDescription}
            minRows={4}
            maxRows={10}
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
      </Modal.Body>
    </Modal>
  )
}

export default EditSongDescriptionModal

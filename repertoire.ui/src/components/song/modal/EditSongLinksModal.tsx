import Song from '../../../types/models/Song.ts'
import { Button, LoadingOverlay, Modal, Stack, TextInput, Tooltip } from '@mantine/core'
import { useForm, zodResolver } from '@mantine/form'
import { EditSongLinksForm, editSongLinksValidation } from '../../../validation/songsForm.ts'
import { IconBrandYoutubeFilled, IconGuitarPickFilled } from '@tabler/icons-react'
import { useUpdateSongMutation } from '../../../state/songsApi.ts'
import { useState } from 'react'
import { toast } from 'react-toastify'

interface EditSongLinksModalProps {
  song: Song
  opened: boolean
  onClose: () => void
}

function EditSongLinksModal({ song, opened, onClose }: EditSongLinksModalProps) {
  const [updateSongMutation, { isLoading }] = useUpdateSongMutation()

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      songsterrLink: song.songsterrLink,
      youtubeLink: song.youtubeLink
    } as EditSongLinksForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editSongLinksValidation),
    onValuesChange: (values) => {
      setHasChanged(
        song.songsterrLink !== values.songsterrLink || song.youtubeLink !== values.youtubeLink
      )
    }
  })

  async function updateSong({ songsterrLink, youtubeLink }: EditSongLinksForm) {
    songsterrLink = songsterrLink?.trim() === '' ? null : songsterrLink?.trim()
    youtubeLink = youtubeLink?.trim() === '' ? null : youtubeLink?.trim()

    await updateSongMutation({
      ...song,
      guitarTuningId: song.guitarTuning?.id,
      id: song.id,
      songsterrLink: songsterrLink,
      youtubeLink: youtubeLink
    }).unwrap()

    onClose()
    setHasChanged(false)
    toast.info('Song links updated!')
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Links'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

        <form onSubmit={form.onSubmit(updateSong)}>
          <Stack>
            <TextInput
              leftSection={<IconGuitarPickFilled size={20} />}
              label="Songsterr"
              placeholder="Songsterr link"
              key={form.key('songsterrLink')}
              {...form.getInputProps('songsterrLink')}
            />
            <TextInput
              leftSection={<IconBrandYoutubeFilled size={20} />}
              label="Youtube"
              placeholder="Youtube link"
              key={form.key('youtubeLink')}
              {...form.getInputProps('youtubeLink')}
            />

            <Tooltip
              disabled={hasChanged}
              label={'You need to make a change before saving'}
              position="bottom"
            >
              <Button
                type={'submit'}
                data-disabled={!hasChanged}
                onClick={(e) => (!hasChanged ? e.preventDefault() : {})}
              >
                Save
              </Button>
            </Tooltip>
          </Stack>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default EditSongLinksModal

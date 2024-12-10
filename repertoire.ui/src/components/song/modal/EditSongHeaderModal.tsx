import Song from '../../../types/models/Song.ts'
import { Button, LoadingOverlay, Modal, Stack, TextInput, Tooltip } from '@mantine/core'
import { useUpdateSongMutation } from '../../../state/songsApi.ts'
import { useState } from 'react'
import { useForm, zodResolver } from '@mantine/form'
import { EditSongHeaderForm, editSongHeaderValidation } from '../../../validation/songsForm.ts'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarFilled } from '@tabler/icons-react'

interface EditSongHeaderModalProps {
  song: Song
  opened: boolean
  onClose: () => void
}

function EditSongHeaderModal({ song, opened, onClose }: EditSongHeaderModalProps) {
  const [updateSongMutation, { isLoading }] = useUpdateSongMutation()

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: song.title,
      releaseDate: new Date(song.releaseDate)
    } as EditSongHeaderForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editSongHeaderValidation),
    onValuesChange: (values) => {
      setHasChanged(
        values.title !== song.title ||
          values.releaseDate.toISOString() !== new Date(song.releaseDate).toISOString()
      )
    }
  })

  async function updateSong({ title, releaseDate }: EditSongHeaderForm) {
    await updateSongMutation({
      ...song,
      guitarTuningId: song.guitarTuning?.id,
      id: song.id,
      title: title,
      releaseDate: releaseDate
    }).unwrap()

    onClose()
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Header'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} />

        <form onSubmit={form.onSubmit(updateSong)}>
          <Stack>
            <TextInput
              withAsterisk={true}
              maxLength={100}
              label="Title"
              placeholder="The title of the song"
              key={form.key('title')}
              {...form.getInputProps('title')}
            />

            <DatePickerInput
              label={'Release Date'}
              leftSection={<IconCalendarFilled size={20} />}
              placeholder={'Choose the release date'}
              key={form.key('releaseDate')}
              {...form.getInputProps('releaseDate')}
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

export default EditSongHeaderModal

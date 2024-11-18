import Album from '../../../types/models/Album.ts'
import { Button, LoadingOverlay, Modal, Stack, TextInput, Tooltip } from '@mantine/core'
import { useUpdateAlbumMutation } from '../../../state/albumsApi.ts'
import { useState } from 'react'
import { useForm, zodResolver } from '@mantine/form'
import { EditAlbumHeaderForm, editAlbumHeaderValidation } from '../../../validation/albumsForm.ts'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarFilled } from '@tabler/icons-react'

interface EditAlbumHeaderModalProps {
  album: Album
  opened: boolean
  onClose: () => void
}

function EditAlbumHeaderModal({ album, opened, onClose }: EditAlbumHeaderModalProps) {
  const [updateAlbumMutation, { isLoading }] = useUpdateAlbumMutation()

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: album.title,
      releaseDate: new Date(album.releaseDate)
    } as EditAlbumHeaderForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editAlbumHeaderValidation),
    onValuesChange: (values) => {
      setHasChanged(
        values.title !== album.title ||
          values.releaseDate.toISOString() !== new Date(album.releaseDate).toISOString()
      )
    }
  })

  async function updateAlbum({ title, releaseDate }) {
    await updateAlbumMutation({
      ...album,
      id: album.id,
      title: title,
      releaseDate: releaseDate
    }).unwrap()

    onClose()
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Album Header'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} />

        <form onSubmit={form.onSubmit(updateAlbum)}>
          <Stack>
            <TextInput
              withAsterisk={true}
              maxLength={100}
              label="Title"
              placeholder="The title of the album"
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

export default EditAlbumHeaderModal

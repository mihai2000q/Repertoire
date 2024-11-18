import Artist from '../../../types/models/Artist.ts'
import { Button, LoadingOverlay, Modal, Stack, TextInput, Tooltip } from '@mantine/core'
import { useUpdateArtistMutation } from '../../../state/artistsApi.ts'
import { useState } from 'react'
import { useForm, zodResolver } from '@mantine/form'
import {
  EditArtistHeaderForm,
  editArtistHeaderValidation
} from '../../../validation/artistsForm.ts'

interface EditArtistHeaderModalProps {
  artist: Artist
  opened: boolean
  onClose: () => void
}

function EditArtistHeaderModal({ artist, opened, onClose }: EditArtistHeaderModalProps) {
  const [updateArtistMutation, { isLoading }] = useUpdateArtistMutation()

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      name: artist.name
    } as EditArtistHeaderForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editArtistHeaderValidation),
    onValuesChange: (values) => {
      setHasChanged(values.name !== artist.name)
    }
  })

  async function updateArtist({ name }) {
    await updateArtistMutation({
      ...artist,
      id: artist.id,
      name: name
    }).unwrap()

    onClose()
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Artist Header'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} />

        <form onSubmit={form.onSubmit(updateArtist)}>
          <Stack>
            <TextInput
              withAsterisk={true}
              maxLength={100}
              label="Name"
              placeholder="The name of the artist"
              key={form.key('name')}
              {...form.getInputProps('name')}
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

export default EditArtistHeaderModal

import Artist from '../../../types/models/Artist.ts'
import { Button, LoadingOverlay, Modal, Stack, TextInput, Tooltip } from '@mantine/core'
import {
  useDeleteImageFromArtistMutation,
  useSaveImageToArtistMutation,
  useUpdateArtistMutation
} from '../../../state/artistsApi.ts'
import { useEffect, useState } from 'react'
import { useForm, zodResolver } from '@mantine/form'
import {
  EditArtistHeaderForm,
  editArtistHeaderValidation
} from '../../../validation/artistsForm.ts'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { toast } from 'react-toastify'
import { useDidUpdate } from '@mantine/hooks'
import { FileWithPath } from '@mantine/dropzone'

interface EditArtistHeaderModalProps {
  artist: Artist
  opened: boolean
  onClose: () => void
}

function EditArtistHeaderModal({ artist, opened, onClose }: EditArtistHeaderModalProps) {
  const [updateArtistMutation, { isLoading: isUpdateLoading }] = useUpdateArtistMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToArtistMutation()
  const [deleteImageMutation, { isLoading: isDeleteImageLoading }] =
    useDeleteImageFromArtistMutation()
  const isLoading = isUpdateLoading || isSaveImageLoading || isDeleteImageLoading

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      name: artist.name,
      image: artist.imageUrl
    } as EditArtistHeaderForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editArtistHeaderValidation),
    onValuesChange: (values) => {
      setHasChanged(values.name !== artist.name || values.image !== artist.imageUrl)
    }
  })

  const [image, setImage] = useState<string | FileWithPath>(artist.imageUrl)
  useEffect(() => form.setFieldValue('image', image), [image])
  useDidUpdate(() => setImage(artist.imageUrl), [artist])

  async function updateArtist({ name, image }: EditArtistHeaderForm) {
    name = name.trim()

    await updateArtistMutation({
      id: artist.id,
      name: name
    }).unwrap()

    if (image !== null && typeof image !== 'string') {
      await saveImageMutation({
        id: artist.id,
        image: image
      })
    } else if (image === null && artist.imageUrl) {
      await deleteImageMutation(artist.id)
    }

    toast.info('Artist updated!')
    onClose()
    setHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Artist Header'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

        <form onSubmit={form.onSubmit(updateArtist)}>
          <Stack>
            <LargeImageDropzoneWithPreview
              image={image}
              setImage={setImage}
              defaultValue={artist.imageUrl}
            />

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

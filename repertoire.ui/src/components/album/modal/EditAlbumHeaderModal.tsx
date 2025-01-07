import Album from '../../../types/models/Album.ts'
import { Button, LoadingOverlay, Modal, Stack, TextInput, Tooltip } from '@mantine/core'
import {
  useDeleteImageFromAlbumMutation,
  useSaveImageToAlbumMutation,
  useUpdateAlbumMutation
} from '../../../state/albumsApi.ts'
import { useEffect, useState } from 'react'
import { useForm, zodResolver } from '@mantine/form'
import { EditAlbumHeaderForm, editAlbumHeaderValidation } from '../../../validation/albumsForm.ts'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarFilled } from '@tabler/icons-react'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'

interface EditAlbumHeaderModalProps {
  album: Album
  opened: boolean
  onClose: () => void
}

function EditAlbumHeaderModal({ album, opened, onClose }: EditAlbumHeaderModalProps) {
  const [updateAlbumMutation, { isLoading: isUpdateLoading }] = useUpdateAlbumMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToAlbumMutation()
  const [deleteImageMutation, { isLoading: isDeleteImageLoading }] =
    useDeleteImageFromAlbumMutation()
  const isLoading = isUpdateLoading || isSaveImageLoading || isDeleteImageLoading

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: album.title,
      releaseDate: album.releaseDate && new Date(album.releaseDate),
      image: album.imageUrl ?? null
    } as EditAlbumHeaderForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editAlbumHeaderValidation),
    onValuesChange: (values) => {
      setHasChanged(
        values.title !== album.title ||
          values.releaseDate?.toISOString() !== new Date(album.releaseDate).toISOString() ||
          values.image !== album.imageUrl
      )
    }
  })

  const [image, setImage] = useState(album.imageUrl ?? null)
  useEffect(() => form.setFieldValue('image', image), [image])
  useEffect(() => setImage(album.imageUrl), [album])

  async function updateAlbum({ title, releaseDate, image }: EditAlbumHeaderForm) {
    title = title.trim()

    await updateAlbumMutation({
      ...album,
      id: album.id,
      title: title,
      releaseDate: releaseDate
    }).unwrap()

    if (image !== null && typeof image !== 'string') {
      await saveImageMutation({
        id: album.id,
        image: image
      })
    } else if (image === null && album.imageUrl) {
      await deleteImageMutation(album.id)
    }

    onClose()
    setHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Album Header'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} />

        <form onSubmit={form.onSubmit(updateAlbum)}>
          <Stack>
            <LargeImageDropzoneWithPreview
              image={image}
              setImage={setImage}
              defaultValue={album.imageUrl}
            />

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

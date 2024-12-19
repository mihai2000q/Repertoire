import Playlist from '../../../types/models/Playlist.ts'
import { Button, LoadingOverlay, Modal, Stack, Textarea, TextInput, Tooltip } from '@mantine/core'
import {
  useDeleteImageFromPlaylistMutation,
  useSaveImageToPlaylistMutation,
  useUpdatePlaylistMutation
} from '../../../state/playlistsApi.ts'
import { useEffect, useState } from 'react'
import { useForm, zodResolver } from '@mantine/form'
import {
  EditPlaylistHeaderForm,
  editPlaylistHeaderValidation
} from '../../../validation/playlistsForm.ts'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'

interface EditPlaylistHeaderModalProps {
  playlist: Playlist
  opened: boolean
  onClose: () => void
}

function EditPlaylistHeaderModal({ playlist, opened, onClose }: EditPlaylistHeaderModalProps) {
  const [updatePlaylistMutation, { isLoading: isUpdateLoading }] = useUpdatePlaylistMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToPlaylistMutation()
  const [deleteImageMutation, { isLoading: isDeleteImageLoading }] =
    useDeleteImageFromPlaylistMutation()
  const isLoading = isUpdateLoading || isSaveImageLoading || isDeleteImageLoading

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: playlist.title,
      description: playlist.description,
      image: playlist?.imageUrl ?? null
    } as EditPlaylistHeaderForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editPlaylistHeaderValidation),
    onValuesChange: (values) => {
      setHasChanged(
        values.title !== playlist.title ||
          values.description !== playlist.description ||
          values.image !== playlist.imageUrl
      )
    }
  })

  const [image, setImage] = useState(playlist?.imageUrl ?? null)
  useEffect(() => form.setFieldValue('image', image), [image])

  async function updatePlaylist({ title, description, image }: EditPlaylistHeaderForm) {
    title = title.trim()

    await updatePlaylistMutation({
      ...playlist,
      id: playlist.id,
      title: title,
      description: description
    }).unwrap()

    if (image !== null && typeof image !== 'string') {
      await saveImageMutation({
        id: playlist.id,
        image: image
      })
    } else if (image === null && playlist.imageUrl) {
      await deleteImageMutation(playlist.id)
    }

    onClose()
    setHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Playlist Header'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} />

        <form onSubmit={form.onSubmit(updatePlaylist)}>
          <Stack>
            <LargeImageDropzoneWithPreview
              image={image}
              setImage={setImage}
              defaultValue={playlist?.imageUrl}
            />

            <TextInput
              withAsterisk={true}
              maxLength={100}
              label="Title"
              placeholder="The title of the playlist"
              key={form.key('title')}
              {...form.getInputProps('title')}
            />

            <Textarea
              label="Description"
              placeholder="The description of the playlist"
              key={form.key('description')}
              {...form.getInputProps('description')}
              rows={4}
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

export default EditPlaylistHeaderModal

import Playlist from '../../../types/models/Playlist.ts'
import { Button, LoadingOverlay, Modal, Stack, Textarea, TextInput, Tooltip } from '@mantine/core'
import {
  useDeleteImageFromPlaylistMutation,
  useSaveImageToPlaylistMutation,
  useUpdatePlaylistMutation
} from '../../../state/api/playlistsApi.ts'
import { useEffect, useState } from 'react'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import {
  EditPlaylistHeaderForm,
  editPlaylistHeaderSchema
} from '../../../validation/playlistsForm.ts'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { useDidUpdate } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { FileWithPath } from '@mantine/dropzone'

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

  const [playlistHasChanged, setPlaylistHasChanged] = useState(false)
  const [imageHasChanged, setImageHasChanged] = useState(false)
  const hasChanged = playlistHasChanged || imageHasChanged

  const form = useForm<EditPlaylistHeaderForm>({
    mode: 'uncontrolled',
    initialValues: {
      title: playlist.title,
      description: playlist.description,
      image: playlist.imageUrl
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(editPlaylistHeaderSchema),
    onValuesChange: (values) => {
      setPlaylistHasChanged(
        values.title !== playlist.title || values.description !== playlist.description
      )
      setImageHasChanged(values.image !== playlist.imageUrl)
    }
  })

  const [image, setImage] = useState<string | FileWithPath>(playlist.imageUrl)
  useEffect(() => form.setFieldValue('image', image), [image])
  useDidUpdate(() => setImage(playlist.imageUrl), [playlist])

  async function updatePlaylist({ title, description, image }: EditPlaylistHeaderForm) {
    if (playlistHasChanged)
      await updatePlaylistMutation({
        id: playlist.id,
        title: title.trim(),
        description: description
      }).unwrap()

    if (image !== null && typeof image !== 'string')
      await saveImageMutation({
        id: playlist.id,
        image: image as FileWithPath
      })
    else if (image === null && playlist.imageUrl) await deleteImageMutation(playlist.id)

    toast.info('Playlist updated!')
    onClose()
    setPlaylistHasChanged(false)
    setImageHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Playlist Header'}>
      <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

      <form onSubmit={form.onSubmit(updatePlaylist)}>
        <Stack px={'xs'} py={0}>
          <LargeImageDropzoneWithPreview
            image={image}
            setImage={setImage}
            defaultValue={playlist.imageUrl}
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
    </Modal>
  )
}

export default EditPlaylistHeaderModal

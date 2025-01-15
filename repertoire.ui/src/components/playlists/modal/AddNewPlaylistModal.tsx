import { useState } from 'react'
import {
  useCreatePlaylistMutation,
  useSaveImageToPlaylistMutation
} from '../../../state/playlistsApi.ts'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import { toast } from 'react-toastify'
import { Button, Group, Modal, Stack, Textarea, TextInput } from '@mantine/core'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import { AddNewPlaylistForm, addNewPlaylistValidation } from '../../../validation/playlistsForm.ts'
import { IconPlaylist } from '@tabler/icons-react'

interface AddNewPlaylistModalProps {
  opened: boolean
  onClose: () => void
}

function AddNewPlaylistModal({ opened, onClose }: AddNewPlaylistModalProps) {
  const [createPlaylistMutation, { isLoading: isCreatePlaylistLoading }] =
    useCreatePlaylistMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToPlaylistMutation()
  const isLoading = isCreatePlaylistLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: '',
      description: ''
    } as AddNewPlaylistForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewPlaylistValidation)
  })

  async function addPlaylist({ title, description }: AddNewPlaylistForm) {
    title = title.trim()

    const res = await createPlaylistMutation({ title, description }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)

    onCloseWithImage()
    form.reset()
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Playlist'} size={475}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addPlaylist)}>
          <Stack>
            <Group>
              <ImageDropzoneWithPreview
                image={image}
                setImage={setImage}
                w={170}
                h={200}
                iconSizes={55}
                icon={<IconPlaylist size={55} />}
              />

              <Stack flex={1}>
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
                  rows={6}
                />
              </Stack>
            </Group>

            <Button style={{ alignSelf: 'center' }} type={'submit'} loading={isLoading}>
              Submit
            </Button>
          </Stack>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default AddNewPlaylistModal

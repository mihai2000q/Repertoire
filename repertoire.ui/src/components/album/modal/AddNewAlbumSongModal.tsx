import { Box, Button, Group, Modal, Stack, Text, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import { AddNewArtistSongForm } from '../../../validation/artistsForm.ts'
import { toast } from 'react-toastify'
import { useCreateSongMutation, useSaveImageToSongMutation } from '../../../state/songsApi.ts'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import { AddNewAlbumSongForm, addNewAlbumSongValidation } from '../../../validation/albumsForm.ts'
import Album from '../../../types/models/Album.ts'
import { IconInfoCircleFilled } from '@tabler/icons-react'

interface AddNewAlbumSongModalProps {
  opened: boolean
  onClose: () => void
  album: Album | undefined
}

function AddNewAlbumSongModal({ opened, onClose, album }: AddNewAlbumSongModalProps) {
  const [createSongMutation, { isLoading: isCreateSongLoading }] = useCreateSongMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToSongMutation()
  const isLoading = isCreateSongLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const inheritedValues = [
    ...(album?.releaseDate ? ['release date'] : []),
    ...(album?.artist ? ['artist'] : [])
  ]

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: ''
    } as AddNewArtistSongForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewAlbumSongValidation)
  })

  async function addSong({ title }: AddNewAlbumSongForm) {
    title = title.trim()

    const res = await createSongMutation({ title, description: '', albumId: album?.id }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)

    onCloseWithImage()
    form.reset()
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Song'}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addSong)}>
          <Stack>
            <Group align={'center'}>
              <ImageDropzoneWithPreview image={image} setImage={setImage} />

              <Stack flex={1} gap={0}>
                <TextInput
                  withAsterisk={true}
                  maxLength={100}
                  label="Title"
                  placeholder="The title of the song"
                  key={form.key('title')}
                  {...form.getInputProps('title')}
                />

                <Stack gap={0} mt={3}>
                  {!image && album?.imageUrl && (
                    <Group gap={4}>
                      <Box c={'cyan.8'}>
                        <IconInfoCircleFilled size={13} />
                      </Box>

                      <Text inline fw={500} c={'dimmed'} fz={'xs'}>
                        If no image is uploaded, it will be inherited.
                      </Text>
                    </Group>
                  )}
                  {inheritedValues.length > 0 && (
                    <Group gap={4} wrap={'nowrap'}>
                      <Box c={'cyan.8'}>
                        <IconInfoCircleFilled size={13} />
                      </Box>

                      <Text inline fw={500} c={'dimmed'} fz={'xs'}>
                        The new song will inherit the <b>{inheritedValues.join(', ')}</b>.
                      </Text>
                    </Group>
                  )}
                </Stack>
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

export default AddNewAlbumSongModal

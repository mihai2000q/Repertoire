import { Button, Group, Modal, Stack, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import {
  AddNewArtistAlbumForm,
  addNewArtistAlbumValidation
} from '../../../validation/artistsForm.ts'
import { toast } from 'react-toastify'
import { useCreateAlbumMutation, useSaveImageToAlbumMutation } from '../../../state/albumsApi.ts'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import { IconAlbum } from '@tabler/icons-react'

interface AddNewArtistAlbumModalProps {
  opened: boolean
  onClose: () => void
  artistId: string | undefined
}

function AddNewArtistAlbumModal({ opened, onClose, artistId }: AddNewArtistAlbumModalProps) {
  const [createAlbumMutation, { isLoading: isCreateAlbumLoading }] = useCreateAlbumMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToAlbumMutation()
  const isLoading = isCreateAlbumLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: ''
    } as AddNewArtistAlbumForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewArtistAlbumValidation)
  })

  async function addAlbum({ title }: AddNewArtistAlbumForm) {
    title = title.trim()

    const res = await createAlbumMutation({ title, artistId }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)

    onCloseWithImage()
    form.reset()
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Album'}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addAlbum)}>
          <Stack>
            <Group align={'center'}>
              <ImageDropzoneWithPreview
                image={image}
                setImage={setImage}
                icon={<IconAlbum size={45} />}
              />
              <TextInput
                flex={1}
                withAsterisk={true}
                maxLength={100}
                label="Title"
                placeholder="The title of the album"
                key={form.key('title')}
                {...form.getInputProps('title')}
              />
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

export default AddNewArtistAlbumModal

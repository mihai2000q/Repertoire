import { Button, Group, Modal, Stack, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import {
  AddNewArtistSongForm,
  addNewArtistSongValidation
} from '../../../validation/artistsForm.ts'
import { toast } from 'react-toastify'
import { useCreateSongMutation, useSaveImageToSongMutation } from '../../../state/songsApi.ts'
import ImageDropzoneWithPreview from '../../image/ImageDropzoneWithPreview.tsx'

interface AddNewArtistSongModalProps {
  opened: boolean
  onClose: () => void
  artistId: string
}

function AddNewArtistSongModal({ opened, onClose, artistId }: AddNewArtistSongModalProps) {
  const [createSongMutation, { isLoading: isCreateSongLoading }] = useCreateSongMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToSongMutation()
  const isLoading = isCreateSongLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: ''
    } as AddNewArtistSongForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewArtistSongValidation)
  })

  async function addSong({ title }: AddNewArtistSongForm) {
    title = title.trim()

    const res = await createSongMutation({ title, description: '', artistId }).unwrap()

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
              <TextInput
                flex={1}
                withAsterisk={true}
                maxLength={100}
                label="Title"
                placeholder="The title of the song"
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

export default AddNewArtistSongModal

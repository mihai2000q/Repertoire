import { Button, Group, Modal, Stack, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import { toast } from 'react-toastify'
import { AddNewArtistForm, addNewArtistValidation } from '../../../validation/artistsForm.ts'
import { useCreateArtistMutation, useSaveImageToArtistMutation } from '../../../state/artistsApi.ts'
import RoundImageDropzoneWithPreview from '../../image/RoundImageDropzoneWithPreview.tsx'

interface AddNewArtistModalProps {
  opened: boolean
  onClose: () => void
}

function AddNewArtistModal({ opened, onClose }: AddNewArtistModalProps) {
  const [createArtistMutation, { isLoading: isCreateArtistLoading }] = useCreateArtistMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToArtistMutation()
  const isLoading = isCreateArtistLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      name: ''
    } as AddNewArtistForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewArtistValidation)
  })

  async function addArtist({ name }: AddNewArtistForm) {
    name = name.trim()

    const res = await createArtistMutation({ name }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${name} added!`)

    onCloseWithImage()
    form.reset()
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Artist'}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addArtist)}>
          <Stack>
            <Group align={'center'}>
              <RoundImageDropzoneWithPreview image={image} setImage={setImage} />
              <TextInput
                flex={1}
                withAsterisk={true}
                maxLength={100}
                label="Name"
                placeholder="The name of the artist"
                key={form.key('name')}
                {...form.getInputProps('name')}
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

export default AddNewArtistModal

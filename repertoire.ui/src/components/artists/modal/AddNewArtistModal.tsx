import { Button, Checkbox, Group, Modal, Stack, Text, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { toast } from 'react-toastify'
import { AddNewArtistForm, addNewArtistSchema } from '../../../validation/artistsForm.ts'
import {
  useCreateArtistMutation,
  useSaveImageToArtistMutation
} from '../../../state/api/artistsApi.ts'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import { IconUserFilled } from '@tabler/icons-react'

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

  const form = useForm<AddNewArtistForm>({
    mode: 'uncontrolled',
    initialValues: {
      name: '',
      isBand: false
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(addNewArtistSchema)
  })

  async function addArtist({ name, isBand }) {
    name = name.trim()

    const res = await createArtistMutation({ name, isBand }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${name} added!`)

    onCloseWithImage()
    form.reset()
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Artist'}>
      <form onSubmit={form.onSubmit(addArtist)}>
        <Stack p={'xs'}>
          <Group>
            <ImageDropzoneWithPreview
              image={image}
              setImage={setImage}
              radius={'50%'}
              icon={<IconUserFilled size={45} />}
            />

            <Stack gap={'xs'} flex={1}>
              <TextInput
                flex={1}
                withAsterisk={true}
                maxLength={100}
                label="Name"
                placeholder="The name of the artist"
                key={form.key('name')}
                {...form.getInputProps('name')}
              />

              <Checkbox
                aria-label={'is-band'}
                label={
                  <Text inline fz={'sm'}>
                    The artist is a <b>band</b>
                  </Text>
                }
                styles={{ label: { paddingLeft: 8 }, labelWrapper: { justifyContent: 'center' } }}
                key={form.key('isBand')}
                {...form.getInputProps('isBand')}
              />
            </Stack>
          </Group>

          <Button style={{ alignSelf: 'center' }} type={'submit'} loading={isLoading}>
            Submit
          </Button>
        </Stack>
      </form>
    </Modal>
  )
}

export default AddNewArtistModal

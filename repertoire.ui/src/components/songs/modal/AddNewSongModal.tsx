import {
  ActionIcon,
  Button,
  FileButton,
  Group,
  Image,
  Modal,
  Stack,
  TextInput,
  Tooltip
} from '@mantine/core'
import { useForm, zodResolver } from '@mantine/form'
import { AddNewSongForm, addNewSongValidation } from '../../../validation/songsForm'
import { useCreateSongMutation, useSaveImageToSongMutation } from '../../../state/songsApi'
import { useState } from 'react'
import { IconPhotoPlus } from '@tabler/icons-react'

interface AddNewSongModalProps {
  opened: boolean
  onClose: () => void
}

function AddNewSongModal({ opened, onClose }: AddNewSongModalProps) {
  const [createSongMutation, { isLoading: isCreateSongLoading }] = useCreateSongMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToSongMutation()
  const isLoading = isCreateSongLoading || isSaveImageLoading

  const [image, setImage] = useState<File>(null)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: ''
    } as AddNewSongForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewSongValidation)
  })

  async function addSong({ title }: AddNewSongForm) {
    const res = await createSongMutation({ title }).unwrap()
    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()
    onClose()
    setImage(null)
    form.reset()
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Add New Song'}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addSong)}>
          <Stack>
            <TextInput
              withAsterisk={true}
              maxLength={100}
              label="Title"
              placeholder="The title of the song"
              key={form.key('title')}
              {...form.getInputProps('title')}
            />

            <Group align={'center'}>
              <FileButton onChange={setImage} accept="image/png,image/jpeg">
                {(props) => (
                  <Tooltip label={'Add a Image'}>
                    <ActionIcon aria-label={'add-image-button'} size={'xl'} {...props}>
                      <IconPhotoPlus />
                    </ActionIcon>
                  </Tooltip>
                )}
              </FileButton>

              {image && (
                <Image
                  src={URL.createObjectURL(image)}
                  h={'100px'}
                  w={'100px'}
                  radius={'md'}
                  alt={'song-image'}
                />
              )}
            </Group>

            <Button type={'submit'} style={{ alignSelf: 'end' }} disabled={isLoading}>
              Add Song
            </Button>
          </Stack>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default AddNewSongModal

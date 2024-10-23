import { Button, Modal, Stack, TextInput } from '@mantine/core'
import { useForm, zodResolver } from '@mantine/form'
import { AddNewSongForm, addNewSongValidation } from '../../../validation/songsForm'
import { useCreateSongMutation } from '../../../state/api'

interface AddNewSongModalProps {
  opened: boolean
  onClose: () => void
}

function AddNewSongModal({ opened, onClose }: AddNewSongModalProps) {
  const [createSongMutation, { isLoading }] = useCreateSongMutation()

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
    await createSongMutation({ title }).unwrap()
    onClose()
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
            <Button type={'submit'} style={{ alignSelf: 'end' }} disabled={isLoading}>
              Add
            </Button>
          </Stack>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default AddNewSongModal

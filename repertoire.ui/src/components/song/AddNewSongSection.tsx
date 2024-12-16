import { useCreateSongSectionMutation } from '../../state/songsApi.ts'
import { Button, Collapse, ComboboxItem, Group, TextInput } from '@mantine/core'
import { useEffect, useState } from 'react'
import { useDidUpdate, useFocusTrap, useInputState } from '@mantine/hooks'
import { toast } from 'react-toastify'
import SongSectionTypeSelect from '../@ui/form/select/SongSectionTypeSelect.tsx'

interface AddNewSongSectionProps {
  songId: string
  opened: boolean
  onClose: () => void
}

function AddNewSongSection({ opened, onClose, songId }: AddNewSongSectionProps) {
  const [createSongSectionMutation, { isLoading: isCreateSongSectionLoading }] =
    useCreateSongSectionMutation()

  const [name, setName] = useInputState('')
  const [nameError, setNameError] = useState(false)
  useDidUpdate(() => setNameError(name.trim().length === 0), [name])

  const [typeError, setTypeError] = useState(false)
  const [type, setType] = useState<ComboboxItem>(null)
  useDidUpdate(() => setTypeError(!type), [type])

  useEffect(() => {
    setNameError(false)
    setTypeError(false)
  }, [opened])

  const nameInputRef = useFocusTrap(opened)

  async function addSection() {
    if (!type || name.trim().length === 0) {
      setTypeError(true)
      setNameError(true)
      return
    }

    const nameTrimmed = name.trim()

    await createSongSectionMutation({ typeId: type.value, name: nameTrimmed, songId }).unwrap()

    toast.success(name + ' added!')

    onClose()
    setName('')
    setType(null)
  }

  return (
    <Collapse in={opened}>
      <Group align={'center'} gap={'xs'} py={'xs'} px={'md'}>
        <SongSectionTypeSelect option={type} onChange={setType} error={typeError} />

        <TextInput
          ref={nameInputRef}
          flex={1}
          maxLength={30}
          placeholder={'Name of Section'}
          value={name}
          onChange={setName}
          error={nameError}
        />

        <Button disabled={isCreateSongSectionLoading} onClick={addSection}>
          Add
        </Button>
      </Group>
    </Collapse>
  )
}

export default AddNewSongSection

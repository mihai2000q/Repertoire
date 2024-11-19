import { useCreateSongSectionMutation, useGetSongSectionTypesQuery } from '../../state/songsApi.ts'
import { Button, Collapse, ComboboxItem, Group, Select, TextInput } from '@mantine/core'
import {useEffect, useState} from 'react'
import { useDidUpdate, useFocusTrap, useInputState } from '@mantine/hooks'
import { toast } from 'react-toastify'

interface AddNewSongSectionProps {
  songId: string
  opened: boolean
  onClose: () => void
}

function AddNewSongSection({ opened, onClose, songId }: AddNewSongSectionProps) {
  const [createSongSectionMutation, { isLoading: isCreateSongSectionLoading }] =
    useCreateSongSectionMutation()

  const { data: songSectionTypesData } = useGetSongSectionTypesQuery()
  const songSectionTypes = songSectionTypesData?.map((type) => ({
    value: type.id,
    label: type.name
  }))

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
        <Select
          clearable={false}
          w={95}
          placeholder={'Type'}
          data={songSectionTypes}
          value={type ? type.value : null}
          onChange={(_, option) => setType(option)}
          maxDropdownHeight={150}
          error={typeError}
        />

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

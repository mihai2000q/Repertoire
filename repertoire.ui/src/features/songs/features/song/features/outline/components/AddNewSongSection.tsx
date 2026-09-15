import { useCreateSongSectionMutation } from '../../sections/state/api/songSectionsApi.ts'
import { Button, Collapse, ComboboxItem, Group, TextInput } from '@mantine/core'
import { useEffect, useState } from 'react'
import { useFocusTrap, useInputState } from '@mantine/hooks'
import { toast } from 'react-toastify'
import SongSectionTypeSelect from '../../../../../../../components/form/select/SongSectionTypeSelect.tsx'

interface AddNewSongSectionProps {
  songId: string
  opened: boolean
  onClose: () => void
  scrollIntoView?: () => void
}

function AddNewSongSection({ songId, opened, onClose, scrollIntoView }: AddNewSongSectionProps) {
  const [createSongSectionMutation, { isLoading }] = useCreateSongSectionMutation()

  const nameInputRef = useFocusTrap(opened)

  const [name, setName] = useInputState('')
  const [nameError, setNameError] = useState(false)
  useEffect(() => setNameError(name.trim().length === 0), [name])

  const [typeError, setTypeError] = useState(false)
  const [type, setType] = useState<ComboboxItem>(null)
  useEffect(() => setTypeError(!type), [type])

  useEffect(() => {
    setNameError(false)
    setTypeError(false)
  }, [opened])

  function handleOnTransitionEnd() {
    if (opened) scrollIntoView()
  }

  async function addSection() {
    if (!type || name.trim().length === 0) {
      setTypeError(!type)
      setNameError(name.trim().length === 0)
      return
    }

    const nameTrimmed = name.trim()

    await createSongSectionMutation({
      typeId: type.value,
      name: nameTrimmed,
      songId: songId,
      partIds: []
    }).unwrap()

    toast.success(nameTrimmed + ' added!')

    onClose()
    setType(null)
    setName('')
  }

  return (
    <Collapse expanded={opened} onTransitionEnd={handleOnTransitionEnd}>
      <Group gap={'xs'} py={'xs'} px={'md'} aria-label={'add-new-song-section'}>
        <Group gap={8}>
          <SongSectionTypeSelect
            w={100}
            option={type}
            onOptionChange={setType}
            error={typeError}
            comboboxProps={{
              position: 'top-start',
              width: 125,
              transitionProps: { duration: 160, transition: 'fade-up' }
            }}
          />
        </Group>

        <TextInput
          ref={nameInputRef}
          flex={1}
          maxLength={30}
          aria-label={'name'}
          placeholder={'Name of Section'}
          value={name}
          onChange={setName}
          error={nameError}
        />

        <Button disabled={isLoading} onClick={addSection}>
          Add
        </Button>
      </Group>
    </Collapse>
  )
}

export default AddNewSongSection

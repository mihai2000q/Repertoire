import { useCreateSongSectionMutation } from '../../state/songsApi.ts'
import { Button, Collapse, ComboboxItem, Group, TextInput } from '@mantine/core'
import { useEffect, useState } from 'react'
import { useDidUpdate, useFocusTrap, useInputState, useScrollIntoView } from '@mantine/hooks'
import { toast } from 'react-toastify'
import SongSectionTypeSelect from '../@ui/form/select/SongSectionTypeSelect.tsx'

interface AddNewSongSectionProps {
  songId: string
  opened: boolean
  onClose: () => void
}

function AddNewSongSection({ opened, onClose, songId }: AddNewSongSectionProps) {
  const [createSongSectionMutation, { isLoading }] = useCreateSongSectionMutation()

  const nameInputRef = useFocusTrap(opened)

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

  const { scrollIntoView, targetRef: scrollIntoViewRef } = useScrollIntoView({
    offset: 20
  })

  function handleOnTransitionEnd() {
    if (opened) scrollIntoView({ alignment: 'end' })
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
      songId: songId
    }).unwrap()

    toast.success(nameTrimmed + ' added!')

    onClose()
    setName('')
    setType(null)
  }

  return (
    <Collapse in={opened} onTransitionEnd={handleOnTransitionEnd}>
      <Group
        ref={scrollIntoViewRef}
        align={'center'}
        gap={'xs'}
        py={'xs'}
        px={'md'}
        aria-label={'add-new-song-section'}
      >
        <SongSectionTypeSelect option={type} onChange={setType} error={typeError} />

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

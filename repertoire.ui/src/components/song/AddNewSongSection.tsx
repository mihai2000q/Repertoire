import { useCreateSongSectionMutation } from '../../state/api/songsApi.ts'
import { Button, Collapse, ComboboxItem, Group, TextInput } from '@mantine/core'
import { useEffect, useState } from 'react'
import { useDidUpdate, useFocusTrap, useInputState, useScrollIntoView } from '@mantine/hooks'
import { toast } from 'react-toastify'
import SongSectionTypeSelect from '../@ui/form/select/SongSectionTypeSelect.tsx'
import { BandMember } from '../../types/models/Artist.ts'
import BandMemberCompactSelect from '../@ui/form/select/BandMemberCompactSelect.tsx'
import InstrumentCompactSelect from '../@ui/form/select/InstrumentCompactSelect.tsx'
import { Instrument, SongSettings } from '../../types/models/Song.ts'

interface AddNewSongSectionProps {
  opened: boolean
  onClose: () => void
  songId: string
  settings: SongSettings
  bandMembers?: BandMember[] | undefined
}

function AddNewSongSection({
  opened,
  onClose,
  songId,
  settings,
  bandMembers
}: AddNewSongSectionProps) {
  const [createSongSectionMutation, { isLoading }] = useCreateSongSectionMutation()

  const { scrollIntoView, targetRef: scrollIntoViewRef } = useScrollIntoView({
    offset: 20
  })
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

  const [bandMember, setBandMember] = useState<BandMember>(null)
  const [instrument, setInstrument] = useState<Instrument>(null)
  useEffect(() => {
    setBandMember(settings.defaultBandMember)
    setInstrument(settings.defaultInstrument)
  }, [settings])

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
      songId: songId,
      bandMemberId: bandMember?.id,
      instrumentId: instrument?.id
    }).unwrap()

    toast.success(nameTrimmed + ' added!')

    onClose()
    setBandMember(settings.defaultBandMember)
    setInstrument(settings.defaultInstrument)
    setType(null)
    setName('')
  }

  return (
    <Collapse in={opened} onTransitionEnd={handleOnTransitionEnd}>
      <Group
        ref={scrollIntoViewRef}
        gap={'xs'}
        py={'xs'}
        px={'md'}
        aria-label={'add-new-song-section'}
      >
        <Group gap={8}>
          <BandMemberCompactSelect
            bandMember={bandMember}
            setBandMember={setBandMember}
            bandMembers={bandMembers}
            position={'top'}
            transitionProps={{ duration: 160, transition: 'fade-up' }}
          />

          <InstrumentCompactSelect
            instrument={instrument}
            setInstrument={setInstrument}
            position={'top'}
            transitionProps={{ duration: 160, transition: 'fade-up' }}
          />

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

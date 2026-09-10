import { useCreateSongPartMutation } from '../state/api/songPartsApi.ts'
import { Button, Collapse, Group, TextInput } from '@mantine/core'
import { useEffect, useState } from 'react'
import { useDidUpdate, useFocusTrap, useInputState } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { BandMember } from '../../../../../../../types/models/Artist.ts'
import BandMemberCompactSelect from '../../../../../../../components/form/select/compact/BandMemberCompactSelect.tsx'
import InstrumentCompactSelect from '../../../../../../../components/form/select/compact/InstrumentCompactSelect.tsx'
import { Instrument, SongSettings } from '../../../../../../../types/models/Song.ts'

interface AddNewSongPartProps {
  opened: boolean
  onClose: () => void
  songId: string
  settings: SongSettings
  bandMembers?: BandMember[] | undefined
  scrollIntoView?: () => void
}

function AddNewSongPart({
  opened,
  onClose,
  songId,
  settings,
  bandMembers,
  scrollIntoView
}: AddNewSongPartProps) {
  const [createSongPartMutation, { isLoading }] = useCreateSongPartMutation()

  const nameInputRef = useFocusTrap(opened)

  const [name, setName] = useInputState('')
  const [nameError, setNameError] = useState(false)
  useEffect(() => setNameError(name.trim().length === 0), [name])

  useEffect(() => {
    setNameError(false)
  }, [opened])

  const [bandMember, setBandMember] = useState<BandMember>(settings.defaultBandMember)
  const [instrument, setInstrument] = useState<Instrument>(settings.defaultInstrument)
  useDidUpdate(() => {
    setBandMember(settings.defaultBandMember)
    setInstrument(settings.defaultInstrument)
  }, [settings])

  function handleOnTransitionEnd() {
    if (opened) scrollIntoView()
  }

  async function addPart() {
    if (name.trim().length === 0) {
      setNameError(name.trim().length === 0)
      return
    }

    const nameTrimmed = name.trim()

    await createSongPartMutation({
      sectionIds: [],
      name: nameTrimmed,
      songId: songId,
      bandMemberId: bandMember?.id,
      instrumentId: instrument?.id
    }).unwrap()

    toast.success(nameTrimmed + ' added!')

    onClose()
    setBandMember(settings.defaultBandMember)
    setInstrument(settings.defaultInstrument)
    setName('')
  }

  return (
    <Collapse expanded={opened} onTransitionEnd={handleOnTransitionEnd}>
      <Group gap={'xs'} py={'xs'} px={'md'} aria-label={'add-new-song-part'}>
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
        </Group>

        <TextInput
          ref={nameInputRef}
          flex={1}
          maxLength={30}
          aria-label={'name'}
          placeholder={'Name of Part'}
          value={name}
          onChange={setName}
          error={nameError}
        />

        <Button disabled={isLoading} onClick={addPart}>
          Add
        </Button>
      </Group>
    </Collapse>
  )
}

export default AddNewSongPart

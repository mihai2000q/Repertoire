import { useCreateSongPartMutation } from '../../parts/state/api/songPartsApi.ts'
import { Button, Collapse, ComboboxItem, Group, TextInput } from '@mantine/core'
import { useEffect, useState } from 'react'
import { useDidUpdate, useFocusTrap, useInputState } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { BandMember } from '../../../../../types/models/Artist.ts'
import BandMemberCompactSelect from './select/compact/BandMemberCompactSelect.tsx'
import InstrumentCompactSelect from './select/compact/InstrumentCompactSelect.tsx'
import { Instrument } from '../../../../../types/models/Song.ts'
import SongSectionSelect from './select/SongSectionSelect.tsx'
import { useSongContext } from '../../../context/SongContext.tsx'

interface AddNewSongPartProps {
  opened: boolean
  onClose: () => void
  scrollIntoView?: () => void
}

function AddNewSongPart({ opened, onClose, scrollIntoView }: AddNewSongPartProps) {
  const { songId, settings, artistBandMembers: bandMembers } = useSongContext()

  const [createSongPartMutation, { isLoading }] = useCreateSongPartMutation()

  const nameInputRef = useFocusTrap(opened)

  const [name, setName] = useInputState('')
  const [nameError, setNameError] = useState(false)
  useEffect(() => setNameError(name.trim().length === 0), [name])

  useEffect(() => {
    setNameError(false)
  }, [opened])

  const [songSection, setSongSection] = useState<ComboboxItem>(null)
  const [bandMember, setBandMember] = useState<BandMember>(settings.defaultBandMember)
  const [instrument, setInstrument] = useState<Instrument>(settings.defaultInstrument)
  useDidUpdate(() => {
    setBandMember(settings.defaultBandMember)
    setInstrument(settings.defaultInstrument)
  }, [settings])

  function handleOnTransitionEnd() {
    if (opened) scrollIntoView?.()
  }

  async function addPart() {
    if (name.trim().length === 0) {
      setNameError(name.trim().length === 0)
      return
    }

    const nameTrimmed = name.trim()

    await createSongPartMutation({
      name: nameTrimmed,
      songId: songId,
      sectionId: songSection?.value,
      bandMemberId: bandMember?.id,
      instrumentId: instrument?.id
    }).unwrap()

    toast.success(nameTrimmed + ' added!')

    onClose()
    setBandMember(settings.defaultBandMember)
    setInstrument(settings.defaultInstrument)
    setSongSection(null)
    setName('')
  }

  return (
    <Collapse expanded={opened} onTransitionEnd={handleOnTransitionEnd}>
      <Group gap={'xs'} py={'xs'} px={'md'} aria-label={'add-new-song-part'}>
        <Group gap={8}>
          <BandMemberCompactSelect
            bandMember={songSection === null ? null : bandMember}
            setBandMember={setBandMember}
            bandMembers={bandMembers}
            position={'top'}
            transitionProps={{ duration: 160, transition: 'fade-up' }}
            disabled={songSection === null}
            tooltipLabel={songSection === null ? 'A song section must be selected' : undefined}
          />

          <InstrumentCompactSelect
            instrument={instrument}
            setInstrument={setInstrument}
            position={'top'}
            transitionProps={{ duration: 160, transition: 'fade-up' }}
          />
        </Group>

        <SongSectionSelect
          w={95}
          option={songSection}
          onOptionChange={setSongSection}
          songId={songId}
          comboboxProps={{
            transitionProps: { duration: 160, transition: 'fade-up' }
          }}
        />

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

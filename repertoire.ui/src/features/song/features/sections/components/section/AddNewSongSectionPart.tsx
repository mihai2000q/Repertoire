import { useCreateSongPartMutation } from '../../../parts/state/api/songPartsApi.ts'
import {
  alpha,
  Button,
  Center,
  CloseButton,
  Collapse,
  Group,
  Stack,
  Text,
  TextInput
} from '@mantine/core'
import { useEffect, useState } from 'react'
import { useDidUpdate, useDisclosure, useFocusTrap, useInputState } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { BandMember } from '../../../../../../types/models/Artist.ts'
import BandMemberCompactSelect from '../../../outline/components/select/compact/BandMemberCompactSelect.tsx'
import InstrumentCompactSelect from '../../../outline/components/select/compact/InstrumentCompactSelect.tsx'
import { Instrument, SongSection } from '../../../../../../types/models/Song.ts'
import { useSongContext } from '../../../../context/SongContext.tsx'
import { IconMusicPlus } from '@tabler/icons-react'

interface AddNewSongPartProps {
  section: SongSection
}

function AddNewSongSectionPart({ section }: AddNewSongPartProps) {
  const { songId, settings, artistBandMembers: bandMembers } = useSongContext()

  const [createSongPartMutation, { isLoading }] = useCreateSongPartMutation()

  const [opened, { open: openAdd, close: closeAdd }] = useDisclosure(false)

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

  async function addPart() {
    if (name.trim().length === 0) {
      setNameError(name.trim().length === 0)
      return
    }

    const nameTrimmed = name.trim()

    await createSongPartMutation({
      name: nameTrimmed,
      songId: songId,
      sectionId: section.id,
      bandMemberId: bandMember?.id,
      instrumentId: instrument?.id
    }).unwrap()

    toast.success(nameTrimmed + ' added!')

    closeAdd()
    setBandMember(settings.defaultBandMember)
    setInstrument(settings.defaultInstrument)
    setName('')
  }

  return (
    <Stack gap={0}>
      <Collapse expanded={!opened}>
        <Group
          px={'sm'}
          py={6}
          gap={'xs'}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.25s',
            color: theme.colors.gray[6],
            '&:hover': {
              boxShadow: theme.shadows.md,
              color: theme.colors.gray[7],
              backgroundColor: alpha(theme.colors.gray[1], 0.35)
            }
          })}
          onClick={openAdd}
          aria-label={`add-new-song-section-part-card-${section.name}`}
        >
          <Center bd={'1px dashed gray'} p={'6px'} style={{ borderRadius: '50%' }}>
            {<IconMusicPlus size={12} />}
          </Center>
          <Text fw={500} fz={'sm'} c={'inherit'} truncate={'end'}>
            Add New Part
          </Text>
        </Group>
      </Collapse>

      <Collapse expanded={opened}>
        <Group gap={'xs'} py={'xs'} px={'md'} aria-label={'add-new-song-part'}>
          <Group gap={8}>
            <CloseButton aria-label={'close'} size={'sm'} onClick={closeAdd} />

            <BandMemberCompactSelect
              bandMember={bandMember}
              setBandMember={setBandMember}
              bandMembers={bandMembers}
              position={'top'}
              transitionProps={{ duration: 160, transition: 'fade-up' }}
              buttonProps={{ size: 25 }}
              iconSize={13}
            />

            <InstrumentCompactSelect
              instrument={instrument}
              setInstrument={setInstrument}
              position={'top'}
              transitionProps={{ duration: 160, transition: 'fade-up' }}
              buttonProps={{ size: 25 }}
              iconSize={13}
            />
          </Group>

          <TextInput
            ref={nameInputRef}
            flex={1}
            size={'xs'}
            maxLength={30}
            aria-label={'name'}
            placeholder={'Name of Part'}
            value={name}
            onChange={setName}
            error={nameError}
          />

          <Button size={'xs'} disabled={isLoading} onClick={addPart}>
            Add
          </Button>
        </Group>
      </Collapse>
    </Stack>
  )
}

export default AddNewSongSectionPart

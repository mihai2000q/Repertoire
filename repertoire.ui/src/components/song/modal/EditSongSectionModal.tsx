import { useUpdateSongSectionMutation } from '../../../state/api/songsApi.ts'
import { useEffect, useState } from 'react'
import {
  Button,
  ComboboxItem,
  Group,
  LoadingOverlay,
  Modal,
  NumberInput,
  Slider,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import { SongSection } from '../../../types/models/Song.ts'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { EditSongSectionForm, editSongSectionSchema } from '../../../validation/songsForm.ts'
import SongSectionTypeSelect from '../../@ui/form/select/compact/SongSectionTypeSelect.tsx'
import { useDidUpdate } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { BandMember } from '../../../types/models/Artist.ts'
import BandMemberSelect from '../../@ui/form/select/BandMemberSelect.tsx'
import InstrumentSelect from '../../@ui/form/select/InstrumentSelect.tsx'

interface EditSongSectionModalProps {
  opened: boolean
  onClose: () => void
  section: SongSection
  bandMembers: BandMember[]
}

function EditSongSectionModal({
  opened,
  onClose,
  section,
  bandMembers
}: EditSongSectionModalProps) {
  const [updateSongSectionMutation, { isLoading }] = useUpdateSongSectionMutation()

  const [hasChanged, setHasChanged] = useState(false)

  const [rehearsalsError, setRehearsalsError] = useState<string | null>()

  const form = useForm<EditSongSectionForm>({
    mode: 'uncontrolled',
    initialValues: {
      name: section.name,
      rehearsals: section.rehearsals,
      confidence: section.confidence,
      typeId: section.songSectionType.id,
      bandMemberId: section.bandMember?.id,
      instrumentId: section.instrument?.id
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(editSongSectionSchema),
    onValuesChange: (values) => {
      setHasChanged(
        values.name !== section.name ||
          (typeof values.rehearsals === 'number' && values.rehearsals !== section.rehearsals) ||
          values.confidence !== section.confidence ||
          values.typeId !== section.songSectionType.id ||
          values.bandMemberId !== section.bandMember?.id ||
          values.instrumentId !== section.instrument?.id
      )

      if (typeof values.rehearsals !== 'number') setRehearsalsError('Cannot be blank')
      else if (values.rehearsals < section.rehearsals)
        setRehearsalsError('Has to be higher than initial value')
      else setRehearsalsError(null)
    }
  })
  useDidUpdate(() => {
    form.setFieldValue('rehearsals', section.rehearsals)
  }, [section])

  const [type, setType] = useState<ComboboxItem>({
    value: section.songSectionType.id,
    label: section.songSectionType.name
  })
  useEffect(() => form.setFieldValue('typeId', type?.value), [type])

  const [bandMember, setBandMember] = useState<BandMember>(section.bandMember)
  useEffect(() => form.setFieldValue('bandMemberId', bandMember?.id), [bandMember])
  useDidUpdate(() => setBandMember(section.bandMember), [section.bandMember])

  const [instrument, setInstrument] = useState<ComboboxItem>(
    section.instrument
      ? {
          value: section.instrument.id,
          label: section.instrument.name
        }
      : undefined
  )
  useEffect(() => form.setFieldValue('instrumentId', instrument?.value), [instrument])
  useDidUpdate(
    () =>
      setInstrument(
        section.instrument
          ? {
              value: section.instrument.id,
              label: section.instrument.name
            }
          : undefined
      ),
    [section.instrument]
  )

  async function updateSongSection({
    name,
    rehearsals,
    confidence,
    bandMemberId,
    instrumentId
  }: EditSongSectionForm) {
    name = name.trim()

    if (rehearsalsError) return

    await updateSongSectionMutation({
      id: section.id,
      typeId: type.value,
      name: name,
      rehearsals: typeof rehearsals !== 'string' ? rehearsals : section.rehearsals,
      confidence: confidence,
      bandMemberId: bandMemberId,
      instrumentId: instrumentId
    }).unwrap()

    onClose()
    setHasChanged(false)
    toast.info(`${name} updated!`)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Section'}>
      <form onSubmit={form.onSubmit(updateSongSection)}>
        <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

        <Stack px={'xs'} py={0}>
          <TextInput
            maxLength={30}
            label="Name"
            placeholder="The name of the song"
            key={form.key('name')}
            {...form.getInputProps('name')}
          />

          <Group>
            <SongSectionTypeSelect
              flex={1}
              label={'Type'}
              placeholder={'Enter Type'}
              option={type}
              onOptionChange={setType}
            />

            <NumberInput
              allowNegative={false}
              allowDecimal={false}
              flex={1}
              label="Rehearsals"
              placeholder="Enter Rehearsals"
              key={form.key('rehearsals')}
              {...form.getInputProps('rehearsals')}
              error={rehearsalsError}
            />
          </Group>

          <Group>
            <BandMemberSelect
              bandMember={bandMember}
              setBandMember={setBandMember}
              bandMembers={bandMembers}
            />
            <InstrumentSelect option={instrument} onOptionChange={setInstrument} flex={1} />
          </Group>

          <Stack gap={0}>
            <Text fw={500} fz={'sm'} c={'black'}>
              Confidence
            </Text>
            <Slider
              thumbLabel={'confidence'}
              label={(value) => `${value}%`}
              key={form.key('confidence')}
              {...form.getInputProps('confidence')}
            />
          </Stack>

          <Tooltip
            disabled={hasChanged}
            label={'You need to make a change before saving'}
            position="bottom"
          >
            <Button
              type={'submit'}
              data-disabled={!hasChanged}
              onClick={(e) => (!hasChanged ? e.preventDefault() : {})}
            >
              Save
            </Button>
          </Tooltip>
        </Stack>
      </form>
    </Modal>
  )
}

export default EditSongSectionModal

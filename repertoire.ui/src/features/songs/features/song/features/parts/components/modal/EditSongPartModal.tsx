import { useUpdateSongPartMutation } from '../../state/api/songPartsApi.ts'
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
import { SongPart } from '../../../../../../../../types/models/Song.ts'
import { schemaResolver, useForm } from '@mantine/form'
import { EditSongPartForm, editSongPartSchema } from '../../../../validation/songForm.ts'
import { useDidUpdate } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { BandMember } from '../../../../../../../../types/models/Artist.ts'
import BandMemberSelect from '../../../outline/components/select/BandMemberSelect.tsx'
import InstrumentSelect from '../../../outline/components/select/InstrumentSelect.tsx'
import { useSongContext } from '../../../../context/SongContext.tsx'

interface EditSongPartModalProps {
  opened: boolean
  onClose: () => void
  part: SongPart
  sectionId?: string
}

function EditSongPartModal({ opened, onClose, part, sectionId }: EditSongPartModalProps) {
  const { artistBandMembers: bandMembers } = useSongContext()

  const [updateSongPartMutation, { isLoading }] = useUpdateSongPartMutation()

  const [hasChanged, setHasChanged] = useState(false)

  const [rehearsalsError, setRehearsalsError] = useState<string | null>()

  const form = useForm<EditSongPartForm>({
    mode: 'uncontrolled',
    initialValues: {
      name: part.name,
      rehearsals: part.rehearsals,
      confidence: part.confidence,
      bandMemberId: sectionId ? part.bandMembers[0]?.id : undefined,
      instrumentId: part.instrument?.id
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: schemaResolver(editSongPartSchema),
    onValuesChange: (values) => {
      setHasChanged(
        values.name !== part.name ||
          (typeof values.rehearsals === 'number' && values.rehearsals !== part.rehearsals) ||
          values.confidence !== part.confidence ||
          values.bandMemberId !== part.bandMembers[0]?.id ||
          values.instrumentId !== part.instrument?.id
      )

      if (typeof values.rehearsals !== 'number') setRehearsalsError('Cannot be blank')
      else if (values.rehearsals < part.rehearsals)
        setRehearsalsError('Has to be higher than initial value')
      else setRehearsalsError(null)
    }
  })
  useEffect(() => {
    form.setFieldValue('rehearsals', part.rehearsals)
    form.setFieldValue('confidence', part.confidence)
  }, [part])

  const [bandMember, setBandMember] = useState<BandMember>(
    sectionId ? part.bandMembers[0] : undefined
  )
  useEffect(() => form.setFieldValue('bandMemberId', bandMember?.id), [bandMember])
  useDidUpdate(() => setBandMember(part.bandMembers[0]), [part.bandMembers])

  const [instrument, setInstrument] = useState<ComboboxItem>(
    part.instrument
      ? {
          value: part.instrument.id,
          label: part.instrument.name
        }
      : undefined
  )
  useEffect(() => form.setFieldValue('instrumentId', instrument?.value), [instrument])
  useDidUpdate(
    () =>
      setInstrument(
        part.instrument
          ? {
              value: part.instrument.id,
              label: part.instrument.name
            }
          : undefined
      ),
    [part.instrument]
  )

  async function updateSongPart({
    name,
    rehearsals,
    confidence,
    bandMemberId,
    instrumentId
  }: EditSongPartForm) {
    name = name.trim()

    if (rehearsalsError) return

    await updateSongPartMutation({
      id: part.id,
      name: name,
      rehearsals: typeof rehearsals !== 'string' ? rehearsals : part.rehearsals,
      confidence: confidence,
      instrumentId: instrumentId,
      sectionId: sectionId,
      bandMemberId: bandMemberId
    }).unwrap()

    onClose()
    setHasChanged(false)
    toast.info(`${name} updated!`)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Part'}>
      <form onSubmit={form.onSubmit(updateSongPart)}>
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
            {sectionId && (
              <BandMemberSelect
                bandMember={bandMember}
                setBandMember={setBandMember}
                bandMembers={bandMembers}
              />
            )}
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
            <Button type={'submit'} disabled={!hasChanged}>
              Save
            </Button>
          </Tooltip>
        </Stack>
      </form>
    </Modal>
  )
}

export default EditSongPartModal

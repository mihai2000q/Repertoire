import { useUpdateSongSectionMutation } from '../../../state/songsApi.ts'
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
import { useForm, zodResolver } from '@mantine/form'
import { EditSongSectionForm, editSongSectionValidation } from '../../../validation/songsForm.ts'
import SongSectionTypeSelect from '../../form/select/SongSectionTypeSelect.tsx'
import { useDidUpdate } from '@mantine/hooks'

interface EditSongSectionModalProps {
  opened: boolean
  onClose: () => void
  section: SongSection
}

function EditSongSectionModal({ opened, onClose, section }: EditSongSectionModalProps) {
  const [updateSongSectionMutation, { isLoading }] = useUpdateSongSectionMutation()

  const [hasChanged, setHasChanged] = useState(false)

  const [rehearsalsError, setRehearsalsError] = useState<string | null>()

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      name: section.name,
      rehearsals: section.rehearsals,
      confidence: section.confidence,
      typeId: section.songSectionType.id
    } as EditSongSectionForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editSongSectionValidation),
    onValuesChange: (values) => {
      setHasChanged(
        values.name !== section.name ||
          (typeof values.rehearsals === 'number' && values.rehearsals !== section.rehearsals) ||
          values.confidence !== section.confidence ||
          values.typeId !== section.songSectionType.id
      )

      if (typeof values.rehearsals !== 'number') setRehearsalsError('Cannot be blank')
      else if (values.rehearsals < section.rehearsals)
        setRehearsalsError('Has to be higher than initial value')
      else setRehearsalsError(null)
    }
  })
  useDidUpdate(() => {
    form.setFieldValue('rehearsals', section.rehearsals) // only rehearsals can be updated from outside
  }, [section])

  const [type, setType] = useState<ComboboxItem>({
    value: section.songSectionType.id,
    label: section.songSectionType.name
  })
  useEffect(() => {
    form.setFieldValue('typeId', type.value)
  }, [type])

  async function updateSongSection({ name, rehearsals, confidence }) {
    if (rehearsalsError) return

    await updateSongSectionMutation({
      id: section.id,
      typeId: type.value,
      name: name,
      rehearsals: rehearsals,
      confidence: confidence
    }).unwrap()

    onClose()
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Section'}>
      <Modal.Body px={'xs'} py={0}>
        <form onSubmit={form.onSubmit(updateSongSection)}>
          <LoadingOverlay visible={isLoading} />

          <Stack>
            <TextInput
              maxLength={100}
              label="Name"
              placeholder="The name of the song"
              key={form.key('name')}
              {...form.getInputProps('name')}
            />

            <Group align={'center'}>
              <SongSectionTypeSelect
                flex={1}
                label={'Type'}
                placeholder={'Enter Type'}
                option={type}
                onChange={setType}
              />

              <NumberInput
                min={0}
                flex={1}
                label="Rehearsals"
                placeholder="Enter Rehearsals"
                key={form.key('rehearsals')}
                {...form.getInputProps('rehearsals')}
                error={rehearsalsError}
              />
            </Group>

            <Stack gap={0}>
              <Text fw={500} fz={'sm'} c={'black'}>
                Confidence
              </Text>
              <Slider
                key={form.key('confidence')}
                {...form.getInputProps('confidence')}
                label={(value) => `${value}%`}
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
      </Modal.Body>
    </Modal>
  )
}

export default EditSongSectionModal
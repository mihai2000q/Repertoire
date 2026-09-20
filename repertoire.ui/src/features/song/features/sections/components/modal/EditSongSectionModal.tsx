import { useUpdateSongSectionMutation } from '../../state/api/songSectionsApi.ts'
import { useEffect, useState } from 'react'
import {
  Button,
  ComboboxItem,
  Group,
  LoadingOverlay,
  Modal,
  Stack,
  TextInput,
  Tooltip
} from '@mantine/core'
import { SongSection } from '../../../../../../types/models/Song.ts'
import { schemaResolver, useForm } from '@mantine/form'
import { EditSongSectionForm, editSongSectionSchema } from '../../../../validation/songForm.ts'
import SongSectionTypeSelect from '../../../../../../components/form/select/SongSectionTypeSelect.tsx'
import { toast } from 'react-toastify'

interface EditSongSectionModalProps {
  opened: boolean
  onClose: () => void
  section: SongSection
}

function EditSongSectionModal({ opened, onClose, section }: EditSongSectionModalProps) {
  const [updateSongSectionMutation, { isLoading }] = useUpdateSongSectionMutation()

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm<EditSongSectionForm>({
    mode: 'uncontrolled',
    initialValues: {
      name: section.name,
      typeId: section.songSectionType.id
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: schemaResolver(editSongSectionSchema),
    onValuesChange: (values) =>
      setHasChanged(values.name !== section.name || values.typeId !== section.songSectionType.id)
  })
  useEffect(() => {
    form.setFieldValue('rehearsals', section.rehearsals)
  }, [section])

  const [type, setType] = useState<ComboboxItem>({
    value: section.songSectionType.id,
    label: section.songSectionType.name
  })
  useEffect(() => form.setFieldValue('typeId', type?.value), [type])

  async function updateSongSection({ name }: EditSongSectionForm) {
    name = name.trim()

    await updateSongSectionMutation({
      id: section.id,
      typeId: type.value,
      name: name,
      partIds: []
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
          </Group>

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

export default EditSongSectionModal

import {
  Combobox,
  Group,
  Input,
  Pill,
  PillsInput,
  PillsInputProps,
  ScrollArea,
  Text,
  useCombobox
} from '@mantine/core'
import { useGetGuitarTuningsQuery } from '../../../../../state/api/songsApi.ts'
import { Dispatch, SetStateAction } from 'react'
import { IconCheck } from '@tabler/icons-react'

interface GuitarTuningMultiSelectProps extends PillsInputProps {
  ids: string[]
  setIds: Dispatch<SetStateAction<string[]>>
  placeholder?: string
  availableIds?: string[]
}

function GuitarTuningMultiSelect({
  ids,
  setIds,
  label,
  placeholder,
  availableIds,
  ...others
}: GuitarTuningMultiSelectProps) {
  const { data, isLoading } = useGetGuitarTuningsQuery()
  const guitarTunings = availableIds
    ? data?.filter((tuning) => availableIds.includes(tuning.id) || ids.includes(tuning.id))
    : data

  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption(),
    onDropdownOpen: () => combobox.updateSelectedOptionIndex('active')
  })

  const handleValueSelect = (id: string) =>
    setIds(ids.includes(id) ? ids.filter((i) => i !== id) : [...ids, id])

  const handleValueRemove = (id: string) => setIds(ids.filter((i) => i !== id))

  const handleValueClear = () => setIds([])

  return (
    <Combobox store={combobox} onOptionSubmit={handleValueSelect} withinPortal={true}>
      <Combobox.DropdownTarget>
        <PillsInput
          disabled={isLoading}
          label={label}
          aria-label={typeof label === 'string' ? label : 'band-member-roles'}
          pointer
          onClick={() => combobox.toggleDropdown()}
          rightSection={
            ids.length > 0 ? (
              <Combobox.ClearButton onClear={handleValueClear} />
            ) : (
              <Combobox.Chevron />
            )
          }
          styles={{ input: { display: 'flex' } }}
          {...others}
        >
          <Pill.Group>
            {ids.length > 0 ? (
              ids.map((id) => (
                <Pill fz={'sm'} key={id} withRemoveButton onRemove={() => handleValueRemove(id)}>
                  {guitarTunings?.find((r) => r.id === id).name}
                </Pill>
              ))
            ) : (
              <Input.Placeholder>{placeholder}</Input.Placeholder>
            )}

            <Combobox.EventsTarget>
              <PillsInput.Field
                type="hidden"
                onBlur={() => combobox.closeDropdown()}
                onKeyDown={(event) => {
                  if (event.key === 'Backspace') {
                    event.preventDefault()
                    handleValueRemove(setIds[ids.length - 1])
                  }
                }}
              />
            </Combobox.EventsTarget>
          </Pill.Group>
        </PillsInput>
      </Combobox.DropdownTarget>

      <Combobox.Dropdown>
        <Combobox.Options>
          <ScrollArea.Autosize mah={150} scrollbarSize={5}>
            {guitarTunings?.map((tuning) => (
              <Combobox.Option value={tuning.id} key={tuning.id} active={ids.includes(tuning.id)}>
                <Group gap={6}>
                  {ids.includes(tuning.id) && <IconCheck size={14} />}
                  <Text fz={'sm'}>{tuning.name}</Text>
                </Group>
              </Combobox.Option>
            ))}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default GuitarTuningMultiSelect

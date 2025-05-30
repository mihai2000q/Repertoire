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
import { useGetBandMemberRolesQuery } from '../../../../../state/api/artistsApi.ts'
import { Dispatch, SetStateAction } from 'react'
import { IconCheck } from '@tabler/icons-react'

interface BandMemberRoleMultiSelectProps extends PillsInputProps {
  ids: string[]
  setIds: Dispatch<SetStateAction<string[]>>
  placeholder?: string
}

function BandMemberRoleMultiSelect({
  ids,
  setIds,
  label,
  placeholder,
  ...others
}: BandMemberRoleMultiSelectProps) {
  const { data: bandMemberRoles, isLoading } = useGetBandMemberRolesQuery()

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
                  {bandMemberRoles?.find((r) => r.id === id).name}
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
            {bandMemberRoles?.map((role) => (
              <Combobox.Option value={role.id} key={role.id} active={ids.includes(role.id)}>
                <Group gap={6}>
                  {ids.includes(role.id) && <IconCheck size={14} />}
                  <Text fz={'sm'}>{role.name}</Text>
                </Group>
              </Combobox.Option>
            ))}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default BandMemberRoleMultiSelect

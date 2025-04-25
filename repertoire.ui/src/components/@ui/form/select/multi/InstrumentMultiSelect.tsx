import {
  Box, Center,
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
import { useGetInstrumentsQuery } from '../../../../../state/api/songsApi.ts'
import { Dispatch, SetStateAction } from 'react'
import useInstrumentIcon from '../../../../../hooks/useInstrumentIcon.tsx'
import { IconCheck } from '@tabler/icons-react'

interface InstrumentMultiSelectProps extends PillsInputProps {
  ids: string[]
  setIds: Dispatch<SetStateAction<string[]>>
  placeholder?: string
  availableIds?: string[]
}

function InstrumentMultiSelect({
  ids,
  setIds,
  label,
  placeholder,
  availableIds,
  ...others
}: InstrumentMultiSelectProps) {
  const { data, isLoading } = useGetInstrumentsQuery()
  const instruments = availableIds
    ? data?.filter((instrument) => availableIds.includes(instrument.id))
    : data

  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption(),
    onDropdownOpen: () => combobox.updateSelectedOptionIndex('active')
  })

  const getInstrumentIcon = useInstrumentIcon()

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
                  <Center c={'primary.7'} w={15} h={'100%'}>
                    {getInstrumentIcon(instruments?.find((r) => r.id === id).name)}
                  </Center>
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
            {instruments?.map((ins) => (
              <Combobox.Option value={ins.id} key={ins.id} active={ids.includes(ins.id)}>
                <Group flex={1} gap={'xs'}>
                  <Box c={'primary.7'} w={19} h={19}>
                    {getInstrumentIcon(ins.name)}
                  </Box>
                  <Text fz={'sm'} fw={500}>
                    {ins.name}
                  </Text>
                  {ids.includes(ins.id) && (
                    <IconCheck
                      size={18}
                      opacity={0.6}
                      stroke={1.5}
                      color={'currentColor'}
                      style={{ marginInlineStart: 'auto' }}
                    />
                  )}
                </Group>
              </Combobox.Option>
            ))}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default InstrumentMultiSelect

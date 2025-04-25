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
import { Dispatch, SetStateAction } from 'react'
import { IconCheck } from '@tabler/icons-react'
import Difficulty from '../../../../../types/enums/Difficulty.ts'

interface DifficultyMultiSelectProps extends PillsInputProps {
  difficulties: Difficulty[]
  setDifficulties: Dispatch<SetStateAction<Difficulty[]>>
  placeholder?: string
  availableDifficulties?: Difficulty[]
}

function DifficultyMultiSelect({
  difficulties,
  setDifficulties,
  label,
  placeholder,
  availableDifficulties,
  ...others
}: DifficultyMultiSelectProps) {
  availableDifficulties ??= Object.entries(Difficulty).map(([_, value]) => value)

  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption(),
    onDropdownOpen: () => combobox.updateSelectedOptionIndex('active')
  })

  const handleValueSelect = (difficulty: Difficulty) =>
    setDifficulties(
      difficulties.includes(difficulty)
        ? difficulties.filter((d) => d !== difficulty)
        : [...difficulties, difficulty]
    )

  const handleValueRemove = (difficulty: string) =>
    setDifficulties(difficulties.filter((d) => d !== difficulty))

  const handleValueClear = () => setDifficulties([])

  return (
    <Combobox store={combobox} onOptionSubmit={handleValueSelect} withinPortal={true}>
      <Combobox.DropdownTarget>
        <PillsInput
          label={label}
          aria-label={typeof label === 'string' ? label : 'difficulties'}
          pointer
          onClick={() => combobox.toggleDropdown()}
          rightSection={
            difficulties.length > 0 ? (
              <Combobox.ClearButton onClear={handleValueClear} />
            ) : (
              <Combobox.Chevron />
            )
          }
          styles={{ input: { display: 'flex' } }}
          {...others}
        >
          <Pill.Group>
            {difficulties.length > 0 ? (
              difficulties.map((d) => (
                <Pill fz={'sm'} key={d} withRemoveButton onRemove={() => handleValueRemove(d)}>
                  {difficulties?.find((d2) => d === d2)}
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
                    handleValueRemove(setDifficulties[difficulties.length - 1])
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
            {availableDifficulties?.map((diff) => (
              <Combobox.Option value={diff} key={diff} active={difficulties.includes(diff)}>
                <Group gap={6}>
                  {difficulties.includes(diff) && <IconCheck size={14} />}
                  <Text fz={'sm'}>{diff}</Text>
                </Group>
              </Combobox.Option>
            ))}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default DifficultyMultiSelect

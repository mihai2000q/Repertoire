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

const allDifficulties = Object.entries(Difficulty).map(([_, value]) => value)

function DifficultyMultiSelect({
  difficulties,
  setDifficulties,
  label,
  placeholder,
  availableDifficulties,
  ...others
}: DifficultyMultiSelectProps) {
  availableDifficulties = availableDifficulties
    ? allDifficulties.filter(
        (diff) => availableDifficulties.includes(diff) || difficulties.includes(diff)
      )
    : allDifficulties

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
              difficulties.map((diff) => (
                <Pill
                  fz={'sm'}
                  key={diff}
                  withRemoveButton
                  onRemove={() => handleValueRemove(diff)}
                >
                  {capitalizeFirstLetter(diff)}
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
            {availableDifficulties.map((diff) => (
              <Combobox.Option value={diff} key={diff} active={difficulties.includes(diff)}>
                <Group gap={6}>
                  {difficulties.includes(diff) && <IconCheck size={14} />}
                  <Text fz={'sm'}>{capitalizeFirstLetter(diff)}</Text>
                </Group>
              </Combobox.Option>
            ))}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

function capitalizeFirstLetter(val: string) {
  return String(val).charAt(0).toUpperCase() + String(val).slice(1);
}

export default DifficultyMultiSelect

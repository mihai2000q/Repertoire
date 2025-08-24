import {
  ActionIcon,
  ActionIconProps,
  Combobox,
  Group,
  ScrollArea,
  Text,
  Tooltip,
  useCombobox
} from '@mantine/core'
import { IconCheck, IconSearch, IconStarFilled } from '@tabler/icons-react'
import { forwardRef, ReactNode, useEffect, useState } from 'react'
import Difficulty from '../../../../../types/enums/Difficulty.ts'
import { useInputState } from '@mantine/hooks'

const allDifficultiesMap = new Map<Difficulty, string>(
  Object.entries(Difficulty).map(([key, value]) => [value, key])
)
const allDifficulties = Array.from(allDifficultiesMap.keys())

interface DifficultySelectButtonProps extends ActionIconProps {
  difficulty: Difficulty | null
  setDifficulty: (difficulty: Difficulty | null) => void
  icon?: ReactNode
}

const DifficultySelectButton = forwardRef<HTMLButtonElement, DifficultySelectButtonProps>(
  ({ difficulty, setDifficulty, icon, ...others }, ref) => {
    const [value, setValue] = useState<string>(allDifficultiesMap[difficulty] ?? '')
    const [search, setSearch] = useInputState(allDifficultiesMap[difficulty] ?? '')

    useEffect(() => {
      setValue(allDifficultiesMap.get(difficulty) ?? '')
      setSearch(allDifficultiesMap.get(difficulty) ?? '')
    }, [difficulty])

    const combobox = useCombobox({
      onDropdownClose: () => {
        combobox.resetSelectedOption()
        combobox.focusTarget()
        setSearch(value)
      },
      onDropdownOpen: () => {
        combobox.focusSearchInput()
      }
    })

    const filteredDifficulties =
      search.trim() !== ''
        ? allDifficulties?.filter((d) => d.includes(search.toLowerCase().trim()))
        : allDifficulties

    const DifficultyOption = ({ difficultyOption }: { difficultyOption: Difficulty }) => (
      <Combobox.Option
        key={difficultyOption}
        value={difficultyOption}
        aria-label={difficultyOption}
        onClick={() => setDifficulty(difficultyOption === difficulty ? null : difficultyOption)}
      >
        <Group gap={'xs'} wrap={'nowrap'}>
          <Text fz={'sm'} lineClamp={1}>
            {allDifficultiesMap.get(difficultyOption)}
          </Text>
          {difficultyOption === difficulty && (
            <IconCheck
              size={16}
              opacity={0.6}
              stroke={1.5}
              color={'currentColor'}
              style={{ marginInlineStart: 'auto' }}
            />
          )}
        </Group>
      </Combobox.Option>
    )

    function handleSubmit(valueString: string) {
      setValue(valueString)
      setSearch(valueString)
      combobox.closeDropdown()
    }

    function handleClear() {
      setDifficulty(null)
    }

    return (
      <Combobox onOptionSubmit={handleSubmit} store={combobox} withArrow>
        <Combobox.Target withAriaAttributes={false}>
          <Tooltip
            disabled={combobox.dropdownOpened}
            label={
              difficulty !== null
                ? `${allDifficultiesMap.get(difficulty)} is selected`
                : 'Select a difficulty'
            }
            openDelay={500}
          >
            <ActionIcon
              ref={ref}
              variant={'form'}
              aria-label={'difficulty'}
              aria-selected={difficulty !== null}
              onClick={() => combobox.toggleDropdown()}
              {...others}
            >
              {icon ?? <IconStarFilled size={16} />}
            </ActionIcon>
          </Tooltip>
        </Combobox.Target>

        <Combobox.Dropdown miw={180} pt={2} px={'xxs'} pb={0}>
          <Combobox.Search
            size={'xs'}
            px={'xxs'}
            pos={'relative'}
            maxLength={100}
            aria-label={'search'}
            placeholder={'Search'}
            leftSection={<IconSearch size={12} />}
            rightSection={difficulty && <Combobox.ClearButton onClear={handleClear} />}
            value={search}
            onChange={setSearch}
            sx={{
              '.mantine-Input-section': {
                position: 'absolute',
                pointerEvents: 'all'
              },
              '.mantine-Combobox-input': { paddingTop: 1 }
            }}
          />
          <Combobox.Options>
            <ScrollArea.Autosize mah={200} scrollbarSize={5}>
              {filteredDifficulties?.length === 0 ? (
                <Combobox.Empty>No Difficulties found</Combobox.Empty>
              ) : (
                filteredDifficulties?.map((difficulty) => (
                  <DifficultyOption key={difficulty} difficultyOption={difficulty} />
                ))
              )}
            </ScrollArea.Autosize>
          </Combobox.Options>
        </Combobox.Dropdown>
      </Combobox>
    )
  }
)

DifficultySelectButton.displayName = 'DifficultySelectButton'

export default DifficultySelectButton

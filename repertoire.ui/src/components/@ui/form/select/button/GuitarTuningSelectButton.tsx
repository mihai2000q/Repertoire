import {
  ActionIcon,
  ActionIconProps,
  Center,
  Combobox,
  Group,
  ScrollArea,
  Text,
  Tooltip,
  useCombobox
} from '@mantine/core'
import { useGetGuitarTuningsQuery } from '../../../../../state/api/songsApi.ts'
import { IconCheck, IconSearch } from '@tabler/icons-react'
import { forwardRef, ReactNode, useEffect, useState } from 'react'
import { GuitarTuning } from '../../../../../types/models/Song.ts'
import CustomIconGuitarHead from '../../../icons/CustomIconGuitarHead.tsx'

interface GuitarTuningSelectButtonProps extends ActionIconProps {
  guitarTuning: GuitarTuning | null
  setGuitarTuning: (guitarTuning: GuitarTuning | null) => void
  icon?: ReactNode
}

const GuitarTuningSelectButton = forwardRef<HTMLButtonElement, GuitarTuningSelectButtonProps>(
  ({ guitarTuning, setGuitarTuning, icon, ...others }, ref) => {
    const { data: guitarTunings, isLoading } = useGetGuitarTuningsQuery()

    const [value, setValue] = useState<string>(guitarTuning?.name ?? '')
    const [search, setSearch] = useState(guitarTuning?.name ?? '')

    useEffect(() => {
      setValue(guitarTuning?.name ?? '')
      setSearch(guitarTuning?.name ?? '')
    }, [guitarTuning])

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

    const filteredGuitarTunings =
      search.trim() !== ''
        ? guitarTunings?.filter((gt) => gt.name.toLowerCase().includes(search.toLowerCase().trim()))
        : guitarTunings

    const GuitarTuningOption = ({ guitarTuningOption }: { guitarTuningOption: GuitarTuning }) => (
      <Combobox.Option
        key={guitarTuningOption.id}
        value={guitarTuningOption.name}
        aria-label={guitarTuningOption.name}
        onClick={() =>
          setGuitarTuning(guitarTuningOption.id === guitarTuning?.id ? null : guitarTuningOption)
        }
      >
        <Group gap={'xs'} wrap={'nowrap'}>
          <Text fz={'sm'} lineClamp={1}>
            {guitarTuningOption.name}
          </Text>
          {guitarTuningOption.id === guitarTuning?.id && (
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
      setGuitarTuning(null)
    }

    return (
      <Combobox onOptionSubmit={handleSubmit} store={combobox} withArrow>
        <Combobox.Target withAriaAttributes={false}>
          <Tooltip
            disabled={combobox.dropdownOpened}
            label={
              guitarTuning !== null ? `${guitarTuning.name} is selected` : 'Select a guitar tuning'
            }
            openDelay={500}
          >
            <ActionIcon
              ref={ref}
              variant={'form'}
              aria-label={'guitar-tuning'}
              aria-selected={guitarTuning !== null}
              disabled={isLoading}
              onClick={() => combobox.toggleDropdown()}
              {...others}
            >
              {icon ?? (
                <Center mr={2}>
                  <CustomIconGuitarHead size={17} />
                </Center>
              )}
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
            rightSection={guitarTuning && <Combobox.ClearButton onClear={handleClear} />}
            value={search}
            onChange={(e) => setSearch(e.currentTarget.value)}
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
              {filteredGuitarTunings?.length === 0 ? (
                <Combobox.Empty>No Guitar Tunings found</Combobox.Empty>
              ) : (
                filteredGuitarTunings?.map((guitarTuning) => (
                  <GuitarTuningOption key={guitarTuning.id} guitarTuningOption={guitarTuning} />
                ))
              )}
            </ScrollArea.Autosize>
          </Combobox.Options>
        </Combobox.Dropdown>
      </Combobox>
    )
  }
)

GuitarTuningSelectButton.displayName = 'GuitarTuningSelectButton'

export default GuitarTuningSelectButton

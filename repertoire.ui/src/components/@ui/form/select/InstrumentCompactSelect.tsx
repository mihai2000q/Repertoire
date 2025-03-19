import {
  ActionIcon,
  alpha,
  Box,
  Combobox,
  ComboboxProps,
  Group,
  ScrollArea,
  Text,
  Tooltip,
  useCombobox
} from '@mantine/core'
import { useGetInstrumentsQuery } from '../../../../state/api/songsApi.ts'
import { IconCheck, IconSearch } from '@tabler/icons-react'
import useInstrumentIcon from '../../../../hooks/useInstrumentIcon.tsx'
import { forwardRef, useEffect, useState } from 'react'
import CustomIconTriangleMusic from '../../icons/CustomIconTriangleMusic.tsx'
import { Instrument } from '../../../../types/models/Song.ts'

interface InstrumentCompactSelectProps extends ComboboxProps {
  instrument: Instrument | null
  setInstrument: (instrument: Instrument | null) => void
  tooltipLabel?: string
}

const InstrumentCompactSelect = forwardRef<HTMLButtonElement, InstrumentCompactSelectProps>(
  ({ instrument, setInstrument, tooltipLabel, ...others }, ref) => {
    const { data: instruments, isLoading } = useGetInstrumentsQuery()

    const getInstrumentIcon = useInstrumentIcon()

    const [value, setValue] = useState<string>(instrument?.name ?? '')
    const [search, setSearch] = useState(instrument?.name ?? '')

    useEffect(() => {
      setValue(instrument?.name ?? '')
      setSearch(instrument?.name ?? '')
    }, [instrument])

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

    const filteredInstruments =
      search.trim() !== ''
        ? instruments?.filter((i) => i.name.toLowerCase().includes(search.toLowerCase().trim()))
        : instruments

    const InstrumentOption = ({ instrumentOption }: { instrumentOption: Instrument }) => (
      <Combobox.Option
        key={instrumentOption.id}
        value={instrumentOption.name}
        aria-label={instrumentOption.name}
        onClick={() =>
          setInstrument(instrumentOption.id === instrument?.id ? null : instrumentOption)
        }
      >
        <Group gap={'xs'} wrap={'nowrap'}>
          <Box c={'primary.7'} w={19} h={19}>
            {getInstrumentIcon(instrumentOption.name)}
          </Box>
          <Text fz={'sm'} fw={500} lineClamp={1}>
            {instrumentOption.name}
          </Text>
          {instrumentOption.id === instrument?.id && (
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
      setInstrument(null)
    }

    return (
      <Combobox onOptionSubmit={handleSubmit} store={combobox} withArrow {...others}>
        <Combobox.Target withAriaAttributes={false}>
          {!instrument ? (
            <Tooltip
              label={tooltipLabel ?? 'Choose an instrument'}
              openDelay={200}
              position={'top'}
              disabled={isLoading || combobox.dropdownOpened}
            >
              <ActionIcon
                ref={ref}
                aria-label={'select-instrument'}
                radius={'50%'}
                sx={(theme) => ({
                  transition: '0.25s',
                  backgroundColor: 'transparent',
                  color: theme.colors.gray[5],
                  border: `1px dashed ${theme.colors.gray[5]}`,

                  '&:hover': {
                    backgroundColor: theme.colors.gray[1],
                    color: alpha(theme.colors.gray[6], 0.85),
                    borderColor: alpha(theme.colors.gray[6], 0.85)
                  },
                  '&[data-disabled="true"]': {
                    backgroundColor: 'transparent',
                    color: alpha(theme.colors.gray[4], 0.8),
                    borderColor: alpha(theme.colors.gray[4], 0.8)
                  }
                })}
                onClick={() => combobox.toggleDropdown()}
                disabled={isLoading}
              >
                <CustomIconTriangleMusic size={14} />
              </ActionIcon>
            </Tooltip>
          ) : (
            <Tooltip
              label={`${instrument.name} is selected`}
              openDelay={200}
              position={'top'}
              disabled={combobox.dropdownOpened}
            >
              <ActionIcon
                ref={ref}
                aria-label={instrument.name}
                variant={'subtle'}
                radius={'50%'}
                c={'primary.7'}
                onClick={() => combobox.toggleDropdown()}
              >
                <Box w={15} h={15}>
                  {getInstrumentIcon(instrument.name)}
                </Box>
              </ActionIcon>
            </Tooltip>
          )}
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
            rightSection={instrument && <Combobox.ClearButton onClear={handleClear} />}
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
              {filteredInstruments?.length === 0 ? (
                <Combobox.Empty>No Instruments found</Combobox.Empty>
              ) : (
                filteredInstruments?.map((instrument) => (
                  <InstrumentOption key={instrument.id} instrumentOption={instrument} />
                ))
              )}
            </ScrollArea.Autosize>
          </Combobox.Options>
        </Combobox.Dropdown>
      </Combobox>
    )
  }
)

InstrumentCompactSelect.displayName = 'InstrumentCompactSelect'

export default InstrumentCompactSelect

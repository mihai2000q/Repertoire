import {
  ActionIcon,
  alpha,
  Box,
  Combobox,
  ComboboxItem,
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
import { useEffect, useState } from 'react'
import CustomIconTriangleMusic from '../../icons/CustomIconTriangleMusic.tsx'

interface InstrumentCompactSelectProps extends ComboboxProps {
  option: ComboboxItem | null
  onOptionChange: (comboboxItem: ComboboxItem | null) => void
}

function InstrumentCompactSelect({
  option,
  onOptionChange,
  ...others
}: InstrumentCompactSelectProps) {
  const { data: instruments, isLoading } = useGetInstrumentsQuery()
  const data = instruments?.map((type) => ({
    value: type.id,
    label: type.name
  }))

  const getInstrumentIcon = useInstrumentIcon()

  const [value, setValue] = useState<string>(option?.label ?? '')
  const [search, setSearch] = useState(option?.label ?? '')

  useEffect(() => {
    setValue(option?.label ?? '')
    setSearch(option?.label ?? '')
  }, [option])

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

  const filteredData =
    search.trim() !== ''
      ? data?.filter((instrument) =>
          instrument.label.toLowerCase().includes(search.toLowerCase().trim())
        )
      : data

  const InstrumentOption = ({ instrumentOption }: { instrumentOption: ComboboxItem }) => (
    <Combobox.Option
      key={instrumentOption.value}
      value={instrumentOption.label}
      aria-label={instrumentOption.label}
      onClick={() => {
        if (instrumentOption.value === option?.value) onOptionChange(null)
        else onOptionChange(instrumentOption)
      }}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <Box c={'primary.7'} w={19} h={19}>
          {getInstrumentIcon(instrumentOption.label)}
        </Box>
        <Text fz={'sm'} fw={500} lineClamp={1}>
          {instrumentOption.label}
        </Text>
        {instrumentOption.value === option?.value && (
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
    setValue(value)
    setSearch(valueString)
    combobox.closeDropdown()
  }

  function handleClear() {
    onOptionChange(null)
  }

  return (
    <Combobox onOptionSubmit={handleSubmit} store={combobox} withArrow {...others}>
      <Combobox.Target withAriaAttributes={false}>
        {!option ? (
          <Tooltip
            label={'Choose an instrument'}
            openDelay={200}
            position={'top'}
            disabled={isLoading || combobox.dropdownOpened}
          >
            <ActionIcon
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
          <Tooltip label={`${option.label} is selected`} openDelay={200} position={'top'}>
            <ActionIcon
              aria-label={option.label}
              variant={'subtle'}
              radius={'50%'}
              c={'primary.7'}
              onClick={() => combobox.toggleDropdown()}
            >
              <Box w={15} h={15}>
                {getInstrumentIcon(option.label)}
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
          rightSection={option && <Combobox.ClearButton onClear={handleClear} />}
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
            {filteredData?.length === 0 ? (
              <Combobox.Empty>No Instruments found</Combobox.Empty>
            ) : (
              filteredData?.map((instrument) => (
                <InstrumentOption key={instrument.value} instrumentOption={instrument} />
              ))
            )}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default InstrumentCompactSelect

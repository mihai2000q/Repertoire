import {Box, ComboboxItem, Group, Select, SelectProps, Text} from '@mantine/core'
import { useGetInstrumentsQuery } from '../../../../state/api/songsApi.ts'
import { IconCheck } from '@tabler/icons-react'
import useInstrumentIcon from '../../../../hooks/useInstrumentIcon.tsx'

interface InstrumentSelectProps extends SelectProps {
  option: ComboboxItem | null
  onOptionChange: (comboboxItem: ComboboxItem | null) => void
}

function InstrumentSelect({ option, onOptionChange, ...others }: InstrumentSelectProps) {
  const { data: instruments, isLoading } = useGetInstrumentsQuery()
  const data = instruments?.map((type) => ({
    value: type.id,
    label: type.name
  }))

  const getInstrumentIcon = useInstrumentIcon()

  const renderSelectOption: SelectProps['renderOption'] = ({ option, checked }) => (
    <Group flex={1} gap={'xs'}>
      <Box c={'primary.7'} w={19} h={19}>
        {getInstrumentIcon(option?.label)}
      </Box>
      <Text fz={'sm'} fw={500}>{option.label}</Text>
      {checked && (
        <IconCheck
          size={18}
          opacity={0.6}
          stroke={1.5}
          color={'currentColor'}
          style={{ marginInlineStart: 'auto' }}
        />
      )}
    </Group>
  )

  return (
    <Select
      disabled={isLoading}
      label={'Instrument'}
      placeholder={'Choose an instrument'}
      data={data}
      value={option?.value ?? null}
      leftSection={
        <Box c={option ? 'primary.7' : 'gray.6'} w={22} h={22}>
          {getInstrumentIcon(option?.label)}
        </Box>
      }
      onChange={(_, option) => onOptionChange(option)}
      renderOption={renderSelectOption}
      maxDropdownHeight={150}
      clearable
      searchable
      {...others}
    />
  )
}

export default InstrumentSelect

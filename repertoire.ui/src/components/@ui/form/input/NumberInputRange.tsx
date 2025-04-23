import { Group, MantineStyleProps, NumberInput, Stack, Text } from '@mantine/core'

interface NumberInputRangeProps extends MantineStyleProps {
  label: string
  value1: string | number
  onChange1: (value: string | number) => void
  value2: string | number
  onChange2: (value: string | number) => void
  max?: number
  isLoading?: boolean
}

function NumberInputRange({
  label,
  value1,
  onChange1,
  value2,
  onChange2,
  max,
  isLoading,
  ...others
}: NumberInputRangeProps) {
  return (
    <Stack gap={2} {...others}>
      {label && (
        <Text fw={500} fz={'sm'}>
          {label}
        </Text>
      )}
      <Group gap={'xs'}>
        <NumberInput
          flex={1}
          aria-label={`min ${label}`}
          allowNegative={false}
          allowDecimal={false}
          max={max}
          disabled={isLoading}
          value={value1}
          onChange={onChange1}
        />
        <Text>-</Text>
        <NumberInput
          flex={1}
          aria-label={`max ${label}`}
          allowNegative={false}
          allowDecimal={false}
          disabled={isLoading}
          value={value2}
          onChange={onChange2}
        />
      </Group>
    </Stack>
  )
}

export default NumberInputRange

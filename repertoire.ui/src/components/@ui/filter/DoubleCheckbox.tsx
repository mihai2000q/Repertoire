import { Checkbox, Group, MantineStyleProps, Stack, Text } from '@mantine/core'
import { useDidUpdate } from '@mantine/hooks'

interface DoubleCheckboxProps extends MantineStyleProps {
  checked1: boolean
  onChange1: (value: boolean) => void
  checked2: boolean
  onChange2: (value: boolean) => void
  title?: string
  label1?: string
  label2?: string
  disabled?: boolean
}

function DoubleCheckbox({
  checked1,
  onChange1,
  checked2,
  onChange2,
  title,
  label1,
  label2,
  disabled
}: DoubleCheckboxProps) {
  useDidUpdate(() => {
    if (checked1 && checked2) {
      onChange2(false)
    }
  }, [checked1])
  useDidUpdate(() => {
    if (checked2 && checked1) {
      onChange1(false)
    }
  }, [checked2])

  return (
    <Stack gap={'xxs'}>
      {title && (
        <Text fw={500} fz={'sm'}>
          {title}
        </Text>
      )}
      <Group aria-label={title}>
        <Checkbox
          label={label1}
          styles={{ label: { paddingLeft: 8 } }}
          disabled={disabled}
          checked={checked1}
          onChange={(value) => onChange1(value.currentTarget.checked)}
        />
        <Checkbox
          label={label2}
          styles={{ label: { paddingLeft: 8 } }}
          disabled={disabled}
          checked={checked2}
          onChange={(value) => onChange2(value.currentTarget.checked)}
        />
      </Group>
    </Stack>
  )
}

export default DoubleCheckbox

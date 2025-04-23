import { Checkbox, Group, MantineStyleProps, Stack, Text } from '@mantine/core'

interface DoubleCheckboxProps extends MantineStyleProps {
  checked1: boolean
  onChange1: (value: boolean) => void
  checked2: boolean
  onChange2: (value: boolean) => void
  title?: string
  label1?: string
  label2?: string
  isLoading?: boolean
}

function DoubleCheckbox({
  checked1,
  onChange1,
  checked2,
  onChange2,
  title,
  label1,
  label2,
  isLoading
}: DoubleCheckboxProps) {
  return (
    <Stack gap={'xxs'}>
      {title && (
        <Text fw={500} fz={'sm'}>
          {title}
        </Text>
      )}
      <Group>
        <Checkbox
          label={label1}
          styles={{ label: { paddingLeft: 8 } }}
          disabled={isLoading}
          checked={checked1}
          onChange={(value) => {
            onChange1(value.currentTarget.checked)
            if (value.currentTarget.checked) {
              onChange2(false)
            }
          }}
        />
        <Checkbox
          label={label2}
          styles={{ label: { paddingLeft: 8 } }}
          disabled={isLoading}
          checked={checked2}
          onChange={(value) => {
            onChange2(value.currentTarget.checked)
            if (value.currentTarget.checked) {
              onChange1(false)
            }
          }}
        />
      </Group>
    </Stack>
  )
}

export default DoubleCheckbox

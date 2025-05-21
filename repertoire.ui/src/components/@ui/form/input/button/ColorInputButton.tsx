import { Dispatch, SetStateAction } from 'react'
import { Button, ColorPicker, Combobox, Group, Popover, Stack, Text, Tooltip } from '@mantine/core'

interface ColorInputButtonProps {
  color: string
  setColor: Dispatch<SetStateAction<string>>
  swatches?: string[]
}

function ColorInputButton({ color, setColor, swatches }: ColorInputButtonProps) {
  return (
    <Stack gap={6}>
      <Text fz={'sm'} fw={500} mt={1}>
        Color
      </Text>
      <Group gap={'xxs'}>
        <Popover shadow={'sm'} withArrow>
          <Popover.Target>
            <Tooltip label={'Choose a color'} openDelay={200}>
              <Button
                aria-label={'color-input'}
                w={25}
                h={25}
                p={0}
                bg={color ?? 'transparent'}
                style={(theme) => ({
                  border: `1px solid ${theme.colors.gray[4]}`,
                  borderWidth: !color ? 1 : 0
                })}
              />
            </Tooltip>
          </Popover.Target>

          <Popover.Dropdown>
            <ColorPicker value={color} onChange={setColor} swatches={swatches} />
          </Popover.Dropdown>
        </Popover>
        <Tooltip label={'Clear color'} openDelay={200} disabled={!color}>
          <Combobox.ClearButton
            opacity={color ? 1 : 0}
            onClear={() => setColor(undefined)}
            style={{ transition: '0.15s', cursor: !color && 'default' }}
          />
        </Tooltip>
      </Group>
    </Stack>
  )
}

export default ColorInputButton

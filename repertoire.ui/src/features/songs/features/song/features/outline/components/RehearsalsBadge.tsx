import { IconRefresh } from '@tabler/icons-react'
import { Group, NumberFormatter, Text } from '@mantine/core'

interface RehearsalsBadgeProps {
  rehearsals: number
}

function RehearsalsBadge({ rehearsals }: RehearsalsBadgeProps) {
  return (
    <Group
      bg={'gray.1'}
      c={'gray.6'}
      gap={'xxs'}
      px={'5px'}
      py={'2px'}
      style={(theme) => ({
        borderRadius: theme.radius.md
      })}
    >
      <IconRefresh size={13} />
      <Text fw={500} fz={'xs'} c={'gray.6'} inline>
        <NumberFormatter value={rehearsals} />
      </Text>
    </Group>
  )
}

export default RehearsalsBadge

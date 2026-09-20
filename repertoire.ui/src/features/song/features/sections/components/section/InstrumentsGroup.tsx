import {
  alpha,
  Avatar,
  AvatarGroup,
  AvatarGroupProps,
  Group,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { Instrument } from '../../../../../../types/models/Song.ts'
import useInstrumentIcon from '../../../../../../hooks/useInstrumentIcon.tsx'

interface InstrumentsGroupProps extends AvatarGroupProps {
  instruments: Instrument[]
}

function InstrumentsGroup({ instruments, ...props }: InstrumentsGroupProps) {
  const getInstrumentIcon = useInstrumentIcon()

  return (
    <Tooltip
      label={
        <Stack gap={'xxs'}>
          {instruments.map((instrument) => (
            <Group key={instrument.id} gap={'xxs'}>
              <Avatar
                bg={'transparent'}
                key={instrument.id}
                aria-label={instrument.name}
                size={'15px'}
                bd={0}
                styles={(theme) => ({
                  placeholder: {
                    color: theme.colors.gray[2],
                    backgroundColor: 'transparent'
                  }
                })}
              >
                {getInstrumentIcon(instrument)}
              </Avatar>
              <Text fz={'xs'} c={'gray.2'}>
                {instrument.name}
              </Text>
            </Group>
          ))}
        </Stack>
      }
      openDelay={300}
    >
      <AvatarGroup {...props}>
        {instruments.map((instrument) => (
          <Avatar
            bg={'transparent'}
            key={instrument.id}
            aria-label={instrument.name}
            size={'28px'}
            bd={0}
            styles={(theme) => ({
              placeholder: {
                padding: '3px',
                color: alpha(theme.colors.gray[6], 0.9),
                backgroundColor: 'transparent'
              }
            })}
          >
            {getInstrumentIcon(instrument)}
          </Avatar>
        ))}
      </AvatarGroup>
    </Tooltip>
  )
}

export default InstrumentsGroup

import { alpha, Avatar, AvatarGroup, AvatarGroupProps, Tooltip } from '@mantine/core'
import { Instrument } from '../../../../../../../../types/models/Song.ts'
import useInstrumentIcon from '../../../../../../../../hooks/useInstrumentIcon.tsx'

interface InstrumentsGroupProps extends AvatarGroupProps {
  instruments: Instrument[]
}

function InstrumentsGroup({ instruments, ...props }: InstrumentsGroupProps) {
  const getInstrumentIcon = useInstrumentIcon()

  const tooltipLabel = instruments.reduce(
    (res, i, idx) => res + i.name + (idx === instruments.length - 1 ? '' : ', '),
    ''
  )

  return (
    <Tooltip label={tooltipLabel} openDelay={300}>
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

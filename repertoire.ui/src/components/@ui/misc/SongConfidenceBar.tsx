import {MantineStyleProps, NumberFormatter, Progress, Tooltip} from '@mantine/core'

interface SongConfidenceBarProps extends MantineStyleProps {
  confidence: number
  size?: number | string
}

function SongConfidenceBar({ confidence, size = 'sm', ...props }: SongConfidenceBarProps) {
  return (
    <Tooltip.Floating
      role={'tooltip'}
      label={
        <>
          <NumberFormatter value={confidence} />%
        </>
      }
    >
      <Progress
        aria-label={'confidence'}
        size={size}
        value={confidence}
        {...props}
      />
    </Tooltip.Floating>
  )
}

export default SongConfidenceBar

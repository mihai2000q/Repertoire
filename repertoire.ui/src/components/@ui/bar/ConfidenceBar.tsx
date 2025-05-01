import { MantineStyleProps, NumberFormatter, Progress, Tooltip } from '@mantine/core'

interface ConfidenceBarProps extends MantineStyleProps {
  confidence: number
  size?: number | string
}

function ConfidenceBar({ confidence, size = 'sm', ...props }: ConfidenceBarProps) {
  return (
    <Tooltip.Floating
      role={'tooltip'}
      label={
        <>
          <NumberFormatter value={confidence} />%
        </>
      }
    >
      <Progress aria-label={'confidence'} size={size} value={confidence} {...props} />
    </Tooltip.Floating>
  )
}

export default ConfidenceBar

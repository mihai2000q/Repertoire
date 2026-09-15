import { MantineStyleProps, Progress, Tooltip } from '@mantine/core'

interface SongOutlineConfidenceBarProps extends MantineStyleProps {
  confidence: number
  size?: string | number
}

function SongOutlineConfidenceBar({
  confidence,
  size = 'xs',
  ...props
}: SongOutlineConfidenceBarProps) {
  return (
    <Tooltip.Floating role={'tooltip'} label={`Confidence: ${confidence}%`}>
      <Progress aria-label={'confidence'} {...props} size={size} value={confidence} />
    </Tooltip.Floating>
  )
}

export default SongOutlineConfidenceBar

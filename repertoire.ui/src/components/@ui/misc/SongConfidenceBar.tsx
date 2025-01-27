import { NumberFormatter, Progress, Tooltip } from '@mantine/core'

interface SongConfidenceBarProps {
  confidence: number
  flex?: string | number
  w?: number
  maw?: number
  miw?: number
  size?: number | string
}

function SongConfidenceBar({ confidence, flex, w, maw, miw, size = 'sm' }: SongConfidenceBarProps) {
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
        flex={flex}
        w={w}
        maw={maw}
        miw={miw}
        size={size}
        value={confidence}
      />
    </Tooltip.Floating>
  )
}

export default SongConfidenceBar

import { NumberFormatter, Progress, Tooltip } from '@mantine/core'

interface ProgressBarProps {
  progress: number
  flex?: string | number
  w?: number
  maw?: number
  miw?: number
  size?: number | string
}

function SongProgressBar({ progress, flex, w, maw, miw, size = 'sm' }: ProgressBarProps) {
  return (
    <Tooltip.Floating role={'tooltip'} label={<NumberFormatter value={progress} />}>
      <Progress
        aria-label={'progress'}
        w={w}
        maw={maw}
        miw={miw}
        flex={flex}
        size={size}
        value={progress / 10}
        color={'green'}
      />
    </Tooltip.Floating>
  )
}

export default SongProgressBar

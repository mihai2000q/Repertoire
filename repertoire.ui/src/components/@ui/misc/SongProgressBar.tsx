import {MantineStyleProps, NumberFormatter, Progress, Tooltip} from '@mantine/core'

interface ProgressBarProps extends MantineStyleProps {
  progress: number
  size?: string | number
}

function SongProgressBar({ progress, size = 'sm', ...props }: ProgressBarProps) {
  return (
    <Tooltip.Floating role={'tooltip'} label={<NumberFormatter value={progress} />}>
      <Progress
        aria-label={'progress'}
        {...props}
        size={size}
        value={progress / 10}
        color={'green'}
      />
    </Tooltip.Floating>
  )
}

export default SongProgressBar

import { MantineStyleProps, NumberFormatter, Progress, Tooltip } from '@mantine/core'

interface SongSectionTypeBadgeProps extends MantineStyleProps {
  progress: number
  maxProgress: number
  size?: string | number
}

function SongOutlineProgressBar({
  progress,
  maxProgress,
  size = 'xs',
  ...props
}: SongSectionTypeBadgeProps) {
  return (
    <Tooltip.Floating
      role={'tooltip'}
      label={
        <>
          Progress: <NumberFormatter value={progress} />
        </>
      }
    >
      <Progress
        aria-label={'progress'}
        {...props}
        size={size}
        value={progress === 0 ? 0 : (progress / maxProgress) * 100}
        color={'green'}
      />
    </Tooltip.Floating>
  )
}

export default SongOutlineProgressBar

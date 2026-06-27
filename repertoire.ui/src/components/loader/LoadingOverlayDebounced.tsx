import { LoadingOverlay, LoadingOverlayProps } from '@mantine/core'
import { useDebouncedValue } from '@mantine/hooks'

interface LoadingOverlayDebouncedProps extends LoadingOverlayProps {
  timeout?: number
}

function LoadingOverlayDebounced({
  timeout = 500,
  visible,
  ...props
}: LoadingOverlayDebouncedProps) {
  const [debouncedVisible] = useDebouncedValue(visible, timeout)

  return <LoadingOverlay visible={debouncedVisible} {...props} />
}

export default LoadingOverlayDebounced
